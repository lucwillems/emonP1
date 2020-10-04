package P1

import (
	"sort"
	"time"
)

const P1TimestampFormat = "060102150405" //YYmmddHHMMSS

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
