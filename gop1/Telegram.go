package gop1

import (
	"errors"
	"fmt"
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
	Id     OBISId
	Info   OType
	Values []*TelegramValue
}

// TelegramValue is one value of a P1 data line, optionally with a specific unit
// of measurement
type TelegramValue struct {
	Valid bool
	Value string
	Unit  string
}

func (t *Telegram) Get(id OBISId) *TelegramObject {
	if i, ok := t.Objects[id]; ok == true {
		return i
	}
	//nil object
	x := TelegramObject{}
	x.Id = OBISTypeNil
	x.Info = TypeInfo[OBISTypeNil]
	x.Values = make([]*TelegramValue, 0)
	return &x
}

func (t *Telegram) Size() int {
	return len(t.Objects)
}
func (t *Telegram) Has(id OBISId) bool {
	_, ok := t.Objects[id]
	return ok
}

func (to *TelegramObject) len() int {
	return len(to.Values)
}

func (to *TelegramObject) Value() *TelegramValue {
	if len(to.Values) > 0 {
		return to.Values[0]
	}
	//return nil value
	x := TelegramValue{}
	x.Valid = false
	x.Value = ""
	x.Unit = ""
	return &x
}

func (to *TelegramObject) AsString() (string, error) {
	v := to.Value()
	if v != nil && v.Valid {
		return v.Value, nil
	}
	return "", errors.New("nil or invalid value")
}

func (to *TelegramObject) AsFloat() (float64, error) {
	v := to.Value()
	if v != nil && v.Valid {
		if f, err := strconv.ParseFloat(v.Value, 64); err == nil {
			return f, nil
		} else {
			return 0, err
		}
	}
	return 0, errors.New("nil TelegramValue")
}

func (to *TelegramObject) AsInt() (int64, error) {
	v := to.Value()
	if v != nil && v.Valid {
		if i, err := strconv.ParseInt(v.Value, 10, 64); err == nil {
			return i, nil
		} else {
			return 0, err
		}
	}
	return 0, errors.New("nil TelegramValue")
}

func (to *TelegramObject) AsDateTime() (time.Time, error) {
	v := to.Value()
	if v != nil && v.Valid {
		return toTimestamp(v.Value)
	}
	return time.Now(), errors.New("nil TelegramValue")
}

func toTimestamp(s string) (time.Time, error) {
	// Remove the DST indicator from the timestamp
	rawDateTime := s[:len(s)-1]
	if location, err := time.LoadLocation(GetTimeZone()); err == nil {
		if dateTime, err := time.ParseInLocation(P1TimestampFormat, rawDateTime, location); err == nil {
			return dateTime, nil
		} else {
			return time.Now(), err
		}
	} else {
		return time.Now(), err
	}
}

func (to *TelegramObject) AsBool() (bool, error) {
	v := to.Value()
	if v != nil && v.Valid {
		if b, err := strconv.ParseBool(v.Value); err == nil {
			return b, nil
		} else {
			return false, err
		}
	}
	return false, errors.New("nil TelegramValue")
}

func (to *TelegramObject) ToString() string {
	v := to.Value()
	return fmt.Sprintf("%-15s: %s %s (%s) %s", to.Id, v.Value, to.Info.Unit, to.Info.Type, to.Info.Description)
}

func GetTimeZone() string {
	return "Europe/Brussels"
}
