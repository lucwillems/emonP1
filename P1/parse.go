package P1

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	errCOSEMNoMatch     = errors.New("COSEM was no match")
	telegramHeaderRegex = regexp.MustCompile(`^\/(.+)$`)
	//cosemOBISRegex      = regexp.MustCompile(`^(\d+-\d+:\d+\.\d+\.\d+)(?:\(([^\)]+)\))+$`)
	cosemOBISRegex      = regexp.MustCompile(`^(\d+-\d+:\d+\.\d+\.\d+)(.+)$`)
	cosemValueUnitRegex = regexp.MustCompile(`^\(([\d\.]+)\*(?i)([a-z0-9]+)\)$`)
	cosemValueRegex     = regexp.MustCompile(`^\(([a-zA-Z0-9]+)\)$`)
	cosemMBusValueUnit  = regexp.MustCompile(`^\((\d{12}[WS]+)\)\(([\d\.]+)\*(?i)([a-z0-9]+)`)
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
				if obj.Id == OBISTypeVersionInformation || obj.Id == OBISTypeBEVersionInfo {
					tgram.Version = obj.Value
				}
				//store obj
				tgram.Objects[obj.Id] = obj
			} else {
				fmt.Fprintf(os.Stderr, "Already exists: %s\n", obj.Id)
			}
		} else {
			//fmt.Fprintf(os.Stderr,err.Error());
			//fmt.Fprintf(os.Stderr,"\n%d: %s",n,l);
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
		fmt.Fprintf(os.Stderr, "unknown id: %s", matches[1])
		return nil, errCOSEMNoMatch
	}
	var x = matches[2]
	//preset common values
	obj.Timestamp = time.Unix(0, 0) //epoch 0
	obj.Unit = ""

	//single (<value>) match ?
	if match := cosemValueRegex.FindStringSubmatch(x); len(match) > 1 {
		obj.Value = match[1]
		return obj, nil
	}

	//single (<value>*<unit>) match ?
	if match := cosemValueUnitRegex.FindStringSubmatch(x); len(match) > 1 {
		obj.Value = match[1]
		obj.Unit = match[2]
		return obj, nil
	}

	if match := cosemMBusValueUnit.FindStringSubmatch(x); len(match) > 1 {
		obj.Value = match[2]
		obj.Unit = match[3]
		obj.Timestamp = toTimestamp(match[1])
		return obj, nil
	}

	//others ?
	obj.Value = x
	return obj, nil
}
