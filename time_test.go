package qiwi

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	type TimeCompare struct {
		Time time.Time `json:"datetime"`
	}

	tests := []struct {
		input, want, format string
	}{
		{`{"datetime": "2021-10-25T16:30:00+03:00"}`, "25 Oct 21 13:30 +0000", time.RFC822Z},
		{`{"datetime": "2021-07-29T16:30:00"}`, "29 Jul 21 16:30", time.RFC822},
		{`{"datetime": "2021-07-29T16:30:00+03:00"}`, "29 Jul 21 16:30 +0300", time.RFC822Z},
	}

	for _, test := range tests {
		var p TimeCompare
		correcttime, _ := time.Parse(test.format, test.want)

		_ = json.Unmarshal([]byte(test.input), &p)
		if !correcttime.Equal(p.Time) {
			t.Errorf("Time parse fail %s %s", p.Time, correcttime)
		}
	}
}

func TestNowInMoscow(t *testing.T) {
	d := 3 * time.Hour
	now := time.Now().UTC().Add(d)

	moscowTime := NowInMoscow()

	if moscowTime.Hour() != now.Hour() {
		t.Error("Wrong NowInMoscow")
	}
}

func TestTimeString(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"25 Oct 21 13:30 +0000", "2021-10-25T13:30:00+00:00"},
		{"29 Jul 21 12:45 +0300", "2021-07-29T12:45:00+03:00"},
		{"29 Jul 21 16:30 +0300", "2021-07-29T16:30:00+03:00"},
	}

	for _, test := range tests {
		parsedTime, _ := time.Parse(time.RFC822Z, test.input)
		qt := &Time{Time: parsedTime}

		if test.want != qt.String() {
			t.Errorf("Time strings not equal %s != %s", qt, test.want)
		}
	}
}
