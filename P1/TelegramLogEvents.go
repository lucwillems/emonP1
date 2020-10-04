package P1

import (
	"fmt"
	"time"
)

type Log struct {
	Timestamp time.Time
	Value     interface{}
	Unit      string
	err       error
}

type LogData struct {
	Id   OBISId
	Logs []*Log
	err  error
}

func (log *Log) String() string {
	s := fmt.Sprint(log.Value)
	return fmt.Sprintf("%s: %s %s", log.Timestamp, s, log.Unit)
}
func (log *Log) IsValid() bool {
	return log.err == nil
}
func (logData *LogData) String() string {
	return fmt.Sprintf("%s [%d] : %s", logData.Id, len(logData.Logs), logData.Logs)
}
func (logData *LogData) IsValid() bool {
	return logData.err == nil
}
