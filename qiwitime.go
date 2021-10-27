// QIWI uses different time format
// for that we implement custom time parsers for JSON
// example: 2021-07-29T16:30:00+03:00
package qiwi

import (
	"fmt"
	"time"
)

// QIWITime holds ISO8601 datetime with timezone YYYY-MM-DDThh:mm:ss±hh:mm
type QIWITime struct {
	time.Time
}

const qiwidate = "2006-01-02T15:04:05"
const qiwitime = "2006-01-02T15:04:05+07:00"

// UnmarshalJSON unpacks QIWI datetime format in go time.Time
func (qt *QIWITime) UnmarshalJSON(b []byte) (err error) {
	s := string(b[1 : len(b)-1])
	qt.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		qt.Time, err = time.Parse(qiwidate, s)
	}
	return err
}
// MarshalJSON packs time.Time to QIWI datetime format
func (qt QIWITime) MarshalJSON() ([]byte, error) {
	// if qt.IsZero() {
	// 	return nil, nil
	// }
	return []byte(fmt.Sprintf(`"%s"`, qt.Format(qiwitime))), nil
}
