package P1

import (
	"fmt"
	"time"
)

// TelegramData is the structured representation of a single line in a P1 data
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

func (td *TelegramData) IsValid() (bool, error) {
	return td.err == nil, td.err
}

func (td *TelegramData) String() string {
	s := fmt.Sprint(td.Value)
	return fmt.Sprintf("%-15s: %s %s (%s) %s", td.Id, s, td.info.Unit, td.info.Type, td.info.Description)
}
