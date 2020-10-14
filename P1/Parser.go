package P1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	errCOSEMNoMatch     = errors.New("line was no COSEM match")
	telegramHeaderRegex = regexp.MustCompile(`^/(.+)$`)
	checksumRegex       = regexp.MustCompile(`^!(.+)$`)
	cosemOBISRegex      = regexp.MustCompile(`^(\d+-\d+:\d+\.\d+\.\d+)(.+)$`)
	cosemValueUnitRegex = regexp.MustCompile(`^\(([\w.]+)(([*\s])([\w]+))?\)$`)
	cosemMBusValueUnit  = regexp.MustCompile(`^\((\d{12}\w)\)\(([\d.]+)(([*\s])([\w]+))?\)$`)
	cosemLogDataRegex   = regexp.MustCompile(`^\((\d{0,2})\)\((\d+-\d+:\d+\.\d+\.\d+)\)(.+)$`)
	cosemLogValueUnit   = regexp.MustCompile(`\(([\w.]+)[*\s]?([\w]+)?\)`)
)

// parsedTelegram parses lines from P1 data, or telegrams
func Parse(message string) *Telegram {
	lines := strings.Split(message, "\n")

	tgram := NewTelegram()
	tgram.Timestamp = time.Now()
	tgram.Version = ""

	for n, l := range lines {
		//skip empty lines
		if len(l) == 0 {
			continue
		}

		// try to detect identification header
		match := telegramHeaderRegex.FindStringSubmatch(l)
		if len(match) > 0 {
			tgram.Device = match[1]
			continue
		}

		// try to detect checksum
		match = checksumRegex.FindStringSubmatch(l)
		if len(match) > 0 {
			tgram.Checksum = match[1]
			continue
		}

		if obj, err := ParseTelegramLine(strings.TrimSpace(l)); err == nil && obj != nil {
			if _, exists := tgram.Objects[obj.Id]; exists == false {
				//telegram timestamp
				if obj.Id == OBISTypeDateTimestamp {
					obj.Timestamp = obj.Value.(time.Time)
				}
				if obj.Id == OBISTypeVersionInformation || obj.Id == OBISTypeBEVersionInfo {
					tgram.Version = obj.Value.(string)
				}
				//store obj
				tgram.Objects[obj.Id] = obj
			} else {
				tgram.Failures++
				fmt.Fprintf(os.Stderr, "%d | Already exists: %s\n", n, obj.Id)
			}
		} else {
			if err != nil {
				tgram.Failures++
				fmt.Fprintf(os.Stderr, "%d | %s\n", n, err.Error())
			}
		}
	}
	return tgram
}

func (td *COSEMInstance) handleCOSUMValues(rawValue string) (*COSEMInstance, error) {
	td.rawValue = rawValue
	td.Value, td.err = convert(td.rawValue, td.info.Type)
	return td, td.err
}

func (td *COSEMInstance) handleCOSUMMBusValues(timestamp string, rawValue string) (*COSEMInstance, error) {
	td.rawValue = rawValue
	if td.Timestamp, td.err = toTimestamp(timestamp); td.err == nil {
		td.Value, td.err = convert(rawValue, td.info.Type)
		return td, nil
	}
	return td, fmt.Errorf("%s: %w", td.Id, td.err)
}

func (td *COSEMInstance) handleCOSUMLog(numbers string, obisId string, logs string) (*COSEMInstance, error) {
	td.rawValue = logs
	var n int64
	if n, td.err = strconv.ParseInt(numbers, 10, 64); td.err != nil {
		return td, fmt.Errorf("%s,%w", td.Id, td.err)
	}

	if i, ok := TypeInfo[obisId]; ok == true {
		logData := LogData{}
		logData.Id = obisId
		logData.Logs = make([]*Log, n)
		logData.info = i
		logData.err = nil
		logData.rawValue = logs

		//we need td parse N logs here
		if match := cosemLogValueUnit.FindAllStringSubmatch(logs, -1); len(match) == (int(n) * 2) {
			for i := 0; i < int(n); i++ {
				log := Log{}
				log.Timestamp, logData.err = toTimestamp(match[i*2][1])
				log.Value, log.err = convert(match[(i*2)+1][1], logData.info.Type)
				logData.Logs[i] = &log
			}
		}
		td.Value = logData
		return td, nil
	} else {
		return nil, errors.New(obisId + ": unknown log event OBIS id")
	}
}

func ParseTelegramLine(line string) (*COSEMInstance, error) {
	matches := cosemOBISRegex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, errCOSEMNoMatch
	}

	var obj *COSEMInstance
	// is this a known COSEM object
	if i, ok := TypeInfo[matches[1]]; ok {
		obj = &COSEMInstance{
			Id:   matches[1],
			info: i,
		}
	} else {
		return nil, errors.New(matches[1] + ": unknown OBIS id")
	}

	var x = matches[2]
	//preset common values
	obj.Timestamp = time.Unix(0, 0) //epoch 0
	if match := cosemValueUnitRegex.FindStringSubmatch(x); len(match) > 1 {
		//single (<value>*<unit>) or (<value>) match ?
		return obj.handleCOSUMValues(match[1])
	} else if match := cosemMBusValueUnit.FindStringSubmatch(x); len(match) > 1 {
		//MBus (<TST>)(<value>*<unit>) or (<TST>)(<value>) match ?
		return obj.handleCOSUMMBusValues(match[1], match[2])
	} else if match := cosemLogDataRegex.FindStringSubmatch(x); len(match) > 1 {
		return obj.handleCOSUMLog(match[1], match[2], match[3])
	} else {
		obj.rawValue = x
		return obj, obj.err
	}
}

func GetTimeZone() string {
	//TODO : make config for this
	return "CET"
}

func toTimestamp(s string) (time.Time, error) {
	// Remove the DST indicator from the timestamp
	rawDateTime := s[:len(s)-1]
	if location, err := time.LoadLocation(GetTimeZone()); err == nil {
		if dateTime, err := time.ParseInLocation(P1TimestampFormat, rawDateTime, location); err == nil {
			return dateTime, nil
		} else {
			return time.Unix(0, 0), err
		}
	} else {
		return time.Unix(0, 0), err
	}
}

func convert(rawData string, dataType string) (interface{}, error) {
	if dataType == Timestamp {
		return toTimestamp(rawData)
	} else if dataType == Integer {
		return strconv.ParseInt(rawData, 10, 64)
	} else if dataType == Float || dataType == MBusFloat {
		return strconv.ParseFloat(rawData, 64)
	} else if dataType == Hex {
		s, err := hex.DecodeString(rawData)
		return string(s), err
	} else if dataType == Bool {
		return strconv.ParseBool(rawData)
	} else if dataType == String {
		return rawData, nil
	} else {
		return rawData, nil
	}
}
