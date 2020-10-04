package P1

import (
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
	telegramHeaderRegex = regexp.MustCompile(`^\/(.+)$`)
	cosemOBISRegex      = regexp.MustCompile(`^(\d+-\d+:\d+\.\d+\.\d+)(.+)$`)
	cosemValueUnitRegex = regexp.MustCompile(`^\(([\w\.]+)(([\*\s])([\w]+))?\)$`)
	cosemMBusValueUnit  = regexp.MustCompile(`^\((\d{12}\w)\)\(([\d\.]+)(([\*\s])([\w]+))?\)$`)
	cosemLogDataRegex   = regexp.MustCompile(`^\((\d{0,2})\)\((\d+-\d+:\d+\.\d+\.\d+)\)(.+)$`)
	cosemLogValueUnit   = regexp.MustCompile(`\(([\w\.]+)[\*\s]?([\w]+)?\)`)
)

// parsedTelegram parses lines from P1 data, or telegrams
func ParseTelegram(lines []string) *Telegram {
	tgram := NewTelegram()
	tgram.Timestamp = time.Now()
	tgram.Version = ""

	for n, l := range lines {
		// try to detect identification header
		match := telegramHeaderRegex.FindStringSubmatch(l)
		if len(match) > 0 {
			tgram.Device = match[1]
			continue
		}

		if obj, err := ParseTelegramLine(strings.TrimSpace(l)); err == nil && obj != nil {
			if _, exists := tgram.Objects[obj.Id]; exists == false {
				//telegram timestamp
				if obj.Id == OBISTypeDateTimestamp {
					if t, err := obj.AsDateTime(); err == nil {
						tgram.Timestamp = t
					}
				}
				if obj.Id == OBISTypeVersionInformation || obj.Id == OBISTypeBEVersionInfo {
					tgram.Version = obj.rawValue
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

func (data *TelegramData) handleCOSUMValues(rawValue string, unit string) (*TelegramData, error) {
	data.rawValue = rawValue
	data.Unit = unit
	convert(data)
	return data, data.err
}

func (data *TelegramData) handleCOSUMMBusValues(timestamp string, rawValue string, unit string) (*TelegramData, error) {
	data.rawValue = rawValue
	data.Unit = unit
	if data.Timestamp, data.err = toTimestamp(timestamp); data.err == nil {
		convert(data)
		return data, nil
	}
	return data, fmt.Errorf("%s: %w", data.Id, data.err)
}

func (data *TelegramData) handleCOSUMLog(numbers string, obisId string, logs string) (*TelegramData, error) {
	data.rawValue = logs
	data.Unit = ""
	var n int64
	if n, data.err = strconv.ParseInt(numbers, 10, 64); data.err != nil {
		return data, fmt.Errorf("%s,%w", data.Id, data.err)
	}
	logData := LogData{}
	logData.Id = OBISId(obisId)
	logData.Logs = make([]*Log, n)
	//we need to parse N logs here
	if match := cosemLogValueUnit.FindAllStringSubmatch(logs, -1); len(match) == (int(n) * 2) {
		for i := 0; i < int(n); i++ {
			log := Log{}
			log.Timestamp, logData.err = toTimestamp(match[i*2][1])
			log.Value = match[(i*2)+1][1]
			log.Unit = match[(i*2)+1][2]
			logData.Logs[i] = &log
		}
	}
	data.Value = logData
	return data, nil
}

func ParseTelegramLine(line string) (*TelegramData, error) {
	if line == "" {
		return nil, nil
	}
	matches := cosemOBISRegex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, errCOSEMNoMatch
	}

	var obj *TelegramData
	// is this a known COSEM object
	if i, ok := TypeInfo[matches[1]]; ok {
		obj = &TelegramData{
			Id:   OBISId(matches[1]),
			info: i,
		}
	} else {
		return nil, errors.New(matches[1] + ": unknown OBIS id")
	}

	var x = matches[2]
	//preset common values
	obj.Timestamp = time.Unix(0, 0) //epoch 0
	obj.Unit = ""

	if match := cosemValueUnitRegex.FindStringSubmatch(x); len(match) > 1 {
		//single (<value>*<unit>) or (<value>) match ?
		return obj.handleCOSUMValues(match[1], match[4])
	} else if match := cosemMBusValueUnit.FindStringSubmatch(x); len(match) > 1 {
		//MBus (<TST>)(<value>*<unit>) or (<TST>)(<value>) match ?
		return obj.handleCOSUMMBusValues(match[1], match[2], match[5])
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

func convert(o *TelegramData) error {
	if o.info.Type == Timestamp {
		o.Value, o.err = o.AsDateTime()
		return o.err
	}
	if o.info.Type == Integer {
		o.Value, o.err = o.AsInt()
		return o.err
	}
	if o.info.Type == Float || o.info.Type == MBusFloat {
		o.Value, o.err = o.AsFloat()
		return o.err
	}
	if o.info.Type == Hex || o.info.Type == String {
		o.Value, o.err = o.AsString()
		return o.err
	}
	return nil
}
