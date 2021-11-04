package qiwi

import (
	"fmt"
	"time"
)

// QIWITime holds ISO8601 datetime with timezone.
// Pattern: YYYY-MM-DDThh:mm:ss±hh:mm.
type QIWITime struct {
	time.Time
}

const (
	moscowtz = "Europe/Moscow"
	qiwidate = "2006-01-02T15:04:05"
	qiwitime = "2006-01-02T15:04:05-07:00"
)

// NowInMoscow returns current time in Moscow (MSK+3).
func NowInMoscow() QIWITime {
	tz, _ := time.LoadLocation(moscowtz)
	return QIWITime{Time: time.Now().In(tz)}
}

// Add delta to time.
func (qt QIWITime) Add(d time.Duration) *QIWITime {
	qt.Time = qt.Time.Add(d)
	return &qt
}

// UnmarshalJSON unpacks QIWI datetime format in go time.Time.
func (qt *QIWITime) UnmarshalJSON(b []byte) (err error) {
	s := string(b[1 : len(b)-1])
	qt.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		qt.Time, err = time.Parse(qiwidate, s)
	}
	return err
}

// MarshalJSON packs time.Time to QIWI datetime format.
func (qt *QIWITime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, qt.Format(qiwitime))), nil
}
