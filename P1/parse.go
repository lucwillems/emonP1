package P1

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	errCOSEMNoMatch     = errors.New("COSEM was no match")
	telegramHeaderRegex = regexp.MustCompile(`^\/(.+)$`)
	cosemOBISRegex      = regexp.MustCompile(`^(\d+-\d+:\d+\.\d+\.\d+)(?:\(([^\)]+)\))+$`)
	cosemUnitRegex      = regexp.MustCompile(`^([\d\.]+)\*(?i)([a-z0-9]+)$`)
)

// parsedTelegram parses lines from P1 data, or telegrams
func ParseTelegram(lines []string) *Telegram {
	tgram := NewTelegram()
	tgram.Timestamp = time.Now()
	tgram.Version = ""

	for _, l := range lines {
		// try to detect identification header
		match := telegramHeaderRegex.FindStringSubmatch(l)
		if len(match) > 0 {
			tgram.Device = match[1]
			continue
		}

		if obj, err := ParseTelegramLine(strings.TrimSpace(l)); err == nil {
			if _, exists := tgram.Objects[obj.Id]; exists == false {
				//telegram timestamp
				if obj.Id == OBISTypeDateTimestamp {
					if t, err := obj.AsDateTime(); err == nil {
						tgram.Timestamp = t
					}
				}
				if obj.Id == OBISTypeVersionInformation {
					tgram.Version = obj.Value().Value
				}
				//store obj
				tgram.Objects[obj.Id] = obj
			} else {
				os.Stderr.WriteString("already exist: " + string(obj.Id))
			}
		}
	}
	return tgram
}

func ParseTelegramLine(line string) (*TelegramObject, error) {
	matches := cosemOBISRegex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, errCOSEMNoMatch
	}

	var obj *TelegramObject
	// is this a known COSEM object
	if i, ok := TypeInfo[matches[1]]; ok {
		obj = &TelegramObject{
			Id:   OBISId(matches[1]),
			Info: i,
		}
	}
	if obj == nil {
		return nil, errCOSEMNoMatch
	}

	for _, v := range matches[2:] {
		ov := TelegramValue{}
		// check if the unit of the value is specified as well
		match := cosemUnitRegex.FindStringSubmatch(v)
		if len(match) > 1 {
			ov.Value = match[1]
			ov.Unit = match[2]
		} else {
			ov.Value = v
		}
		ov.Valid = true
		obj.Values = append(obj.Values, &ov)
	}

	return obj, nil
}
