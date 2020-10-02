package P1

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"time"
)

const P1TimestampFormat = "060102150405"

type Telegram struct {
	Device    string
	Version   string
	Timestamp time.Time
	Objects   map[OBISId]*TelegramObject
}

func NewTelegram() *Telegram {
	var t Telegram
	t.Objects = make(map[OBISId]*TelegramObject)
	return &t
}

// TelegramObject is the structured representation of a sinle line in a P1 data
// dump. It can have one or more values
type TelegramObject struct {
	Id        OBISId
	Info      OType
	Value     string
	Unit      string
	Timestamp time.Time
}

func (t *Telegram) Get(id OBISId) *TelegramObject {
	if i, ok := t.Objects[id]; ok == true {
		return i
	}
	//nil object
	x := TelegramObject{}
	x.Id = OBISTypeNil
	x.Info = TypeInfo[OBISTypeNil]
	return &x
}

func (t *Telegram) SortedIds() []string {
	// To store the keys in slice in sorted order
	keys := make([]string, t.Size())
	i := 0
	for k := range t.Objects {
		keys[i] = string(k)
		i++
	}
	sort.Strings(keys)
	return keys
}

func (t *Telegram) Size() int {
	return len(t.Objects)
}
func (t *Telegram) Has(id OBISId) bool {
	_, ok := t.Objects[id]
	return ok
}

func (to *TelegramObject) AsString() (string, error) {
	if to.Info.Type == Timestamp {
		if t, err := to.AsDateTime(); err == nil {
			return t.Format(time.RFC3339), nil
		} else {
			return "", err
		}
	}
	if to.Info.Type == Hex {
		s, err := hex.DecodeString(to.Value)
		return string(s), err
	}
	return to.Value, nil
}

func (to *TelegramObject) AsFloat() (float64, error) {
	if f, err := strconv.ParseFloat(to.Value, 64); err == nil {
		return f, nil
	} else {
		return 0, err
	}
}

func (to *TelegramObject) AsInt() (int64, error) {
	if i, err := strconv.ParseInt(to.Value, 10, 64); err == nil {
		return i, nil
	} else {
		return 0, err
	}
}

func (to *TelegramObject) AsDateTime() (time.Time, error) {
	return toTimestamp(to.Value), nil
}

func toTimestamp(s string) time.Time {
	// Remove the DST indicator from the timestamp
	rawDateTime := s[:len(s)-1]
	if location, err := time.LoadLocation(GetTimeZone()); err == nil {
		if dateTime, err := time.ParseInLocation(P1TimestampFormat, rawDateTime, location); err == nil {
			return dateTime
		} else {
			return time.Unix(0, 0)
		}
	}
	return time.Unix(0, 0)
}

func (to *TelegramObject) AsBool() (bool, error) {
	if b, err := strconv.ParseBool(to.Value); err == nil {
		return b, nil
	} else {
		return false, err
	}
}

func (to *TelegramObject) ToString() string {
	s, _ := to.AsString()
	return fmt.Sprintf("%-15s: %s %s (%s) %s", to.Id, s, to.Info.Unit, to.Info.Type, to.Info.Description)
}

func GetTimeZone() string {
	//TODO : make config for this
	return "Europe/Brussels"
}
