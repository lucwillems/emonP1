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
	Failures  int
	Objects   map[OBISId]*TelegramData
}

func NewTelegram() *Telegram {
	var t Telegram
	t.Objects = make(map[OBISId]*TelegramData)
	return &t
}

// TelegramData is the structured representation of a singl line in a P1 data
type TelegramData struct {
	Id        OBISId
	Value     interface{} //holds the converted values from string -> int64/float64/time..
	Unit      string
	Timestamp time.Time
	//internal state data
	info     OType
	rawValue string
	err      error
}

type TST string

func (t *Telegram) Get(id OBISId) *TelegramData {
	if i, ok := t.Objects[id]; ok == true {
		return i
	}
	//nil object
	x := TelegramData{}
	x.Id = OBISTypeNil
	x.Value = nil
	x.info = TypeInfo[OBISTypeNil]
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

func (to *TelegramData) AsString() (string, error) {
	if to.info.Type == Timestamp {
		if t, err := to.AsDateTime(); err == nil {
			return t.Format(time.RFC3339), nil
		} else {
			return "", err
		}
	}
	if to.info.Type == Hex {
		s, err := hex.DecodeString(to.rawValue)
		return string(s), err
	}
	return to.rawValue, nil
}

func (to *TelegramData) AsFloat() (float64, error) {
	if f, err := strconv.ParseFloat(to.rawValue, 64); err == nil {
		return f, nil
	} else {
		return 0, err
	}
}

func (to *TelegramData) AsInt() (int64, error) {
	if i, err := strconv.ParseInt(to.rawValue, 10, 64); err == nil {
		return i, nil
	} else {
		return 0, err
	}
}

func (to *TelegramData) AsDateTime() (time.Time, error) {
	return toTimestamp(to.rawValue)
}

func (to *TelegramData) AsBool() (bool, error) {
	if b, err := strconv.ParseBool(to.rawValue); err == nil {
		return b, nil
	} else {
		return false, err
	}
}

func (to *TelegramData) IsValid() (bool, error) {
	return to.err == nil, to.err
}

func (to *TelegramData) String() string {
	s, _ := to.AsString()
	return fmt.Sprintf("%-15s: %s %s (%s) %s", to.Id, s, to.info.Unit, to.info.Type, to.info.Description)
}
