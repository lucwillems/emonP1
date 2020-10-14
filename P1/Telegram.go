package P1

import (
	"time"
)

const P1TimestampFormat = "060102150405" //YYmmddHHMMSS

type Telegram struct {
	Device    string
	Version   string
	Timestamp time.Time
	Failures  int
	Objects   map[string]*COSEMInstance
	Checksum  string
}

func NewTelegram() *Telegram {
	var t Telegram
	t.Objects = make(map[string]*COSEMInstance)
	return &t
}

func (t *Telegram) Get(id string) (*COSEMInstance, bool) {
	if i, ok := t.Objects[id]; ok == true {
		return i, true
	}
	//nil object
	return nil, false
}

func (t *Telegram) OBISIds() []string {
	// To store the keys in slice in sorted order
	keys := make([]string, t.Size())
	i := 0
	for k := range t.Objects {
		keys[i] = string(k)
		i++
	}
	return keys
}

func (t *Telegram) Size() int {
	return len(t.Objects)
}
func (t *Telegram) Has(id string) bool {
	_, ok := t.Objects[id]
	return ok
}
