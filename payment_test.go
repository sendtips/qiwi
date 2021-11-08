package qiwi

import (
	"encoding/json"
	"testing"
)

func TestStatusCodeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input string
		want  StatusCode
	}{
		{`{"value": "CREATED"}`, StatusCreated},
		{`{"value": "WAITING"}`, StatusWait},
		{`{"value": "COMPLETED"}`, StatusCompleted},
		{`{"value": "SUCCESS"}`, StatusOK},
		{`{"value": "DECLINE"}`, StatusFail},
		{`{"value": "DECLINED"}`, StatusFail},
		{`{"value": ""}`, StatusFail},
	}

	for _, test := range tests {
		var res Status
		_ = json.Unmarshal([]byte(test.input), &res)

		if res.Value != test.want {
			t.Errorf("Incorrect status value have: %v, want: %s", res.Value, test.want)
		}
	}
}
