package qiwi

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestAmount(t *testing.T) {
	tests := []struct {
		input int
		want  money
	}{
		{100, 1.00},
		{500, 5.00},
		{101, 1.01},
		{1, 0.01},
		{99, 0.99},
		{150, 1.50},
		{199, 1.99},
	}

	for _, test := range tests {
		a := newAmount(test.input, rub)
		if test.want != a.Value {
			t.Errorf("Wrong amount %0.2f, must be %0.2f", a.Value, test.want)
		}
	}
}

func TestJSONAmount(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{1, `{"value":0.01,"currency":"RUB"}`},
		{10, `{"value":0.10,"currency":"RUB"}`},
		{100, `{"value":1.00,"currency":"RUB"}`},
		{10101, `{"value":101.01,"currency":"RUB"}`},
		{10000, `{"value":100.00,"currency":"RUB"}`},
	}

	for _, test := range tests {
		a := NewAmountInRubles(test.input)

		j, err := json.Marshal(a)
		if err != nil {
			t.Errorf("JSON Marshal error: %s", err)
		}

		if string(j) != test.want {
			t.Errorf("Wrong JSON result: %s", j)
		}

	}
}

func TestJSONUnmarshalAmount(t *testing.T) {
	tests := []struct {
		want  int
		err   error
		input string
	}{
		{1, nil, `{"value":0.01,"currency":"RUB"}`},
		{10, nil, `{"value":0.10,"currency":"RUB"}`},
		{100, nil, `{"value":1.00,"currency":"RUB"}`},
		{10101, nil, `{"value":101.01,"currency":"RUB"}`},
		{10000, nil, `{"value":100.00,"currency":"RUB"}`},
		{12000, nil, `{"value":"120.00","currency":"RUB"}`},
		{0, nil, `{"value":0.00,"currency":"RUB"}`},
		{0, nil, `{"value":-0.00,"currency":"RUB"}`},
		{-100, nil, `{"value":-1.00,"currency":"RUB"}`},
		{0, ErrBadJSON, `{"value": "BADVAL","currency":"RUB"}`},
	}

	for _, test := range tests {
		var am Amount
		err := json.Unmarshal([]byte(test.input), &am)
		if err != nil {
			if !errors.Is(err, test.err) {
				t.Errorf("JSON Marshal error: %s", err)
			} // 123
		}

		if am.Value.Int() != test.want {
			t.Errorf("Wrong JSON result: %d, expected %d", am.Value.Int(), test.want)
		}

	}
}

func TestAmountInRuble(t *testing.T) {
	a := NewAmountInRubles(100)
	if a.Currency != rub {
		t.Error("Bad currency value")
	}
}

func TestAmountInDollars(t *testing.T) {
	a := NewAmountInDollars(100)
	if a.Currency != usd {
		t.Error("Bad currency value")
	}
}

func TestAmountInEuros(t *testing.T) {
	a := NewAmountInEuros(100)
	if a.Currency != eur {
		t.Error("Bad currency value")
	}
}
