package qiwi

import (
	"fmt"
	"strconv"
)

type currency string

const (
	RUB currency = "RUB"
	USD currency = "USD"
	EUR currency = "EUR"
	//GBP Currency = "GBP"
)

// kopeeksInRuble used to make float amount value from int
const kopeeksInRuble float64 = 100

type money float64

func toMoney(a int) money {
	return money(float64(a) / kopeeksInRuble)
}

// Int returns amount in kopeks/cents
func (m money) Int() int {
	return int(float64(m) * kopeeksInRuble)
}

// String returns float with amount representation
func (m money) String() string {
	return fmt.Sprintf("%.2f", m)
}

func (m money) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.2f", float64(m))
	return []byte(s), nil
}

func (m *money) UnmarshalJSON(data []byte) error {
	const doublequotesign byte = 32 // "
	if data[0] == doublequotesign {
		data = data[1 : len(data)-1]
	}
	am, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return fmt.Errorf("[QIWI] Amount JSON error: %w", ErrBadJSON)
	}

	*m = money(am)

	return err
}

type Amount struct {
	Value    money    `json:"value"`
	Currency currency `json:"currency"`
}

func newAmount(a int, cur currency) Amount {
	return Amount{Value: toMoney(a), Currency: cur}
}

// NewAmountInRubles make sets rubles amount
func NewAmountInRubles(a int) Amount {
	return newAmount(a, RUB)
}

// NewAmountInDollars make sets dollars amount
func NewAmountInDollars(a int) Amount {
	return newAmount(a, USD)
}

// NewAmountInEuros make sets euros amount
func NewAmountInEuros(a int) Amount {
	return newAmount(a, EUR)
}
