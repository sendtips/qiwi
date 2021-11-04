package qiwi

import (
	"fmt"
	"strconv"
)

type currency string

const (
	rub currency = "RUB"
	usd currency = "USD"
	eur currency = "EUR"
)

// kopeeksInRuble used to make float amount value from int.
const kopeeksInRuble float64 = 100

type money float64

func toMoney(a int) money {
	return money(float64(a) / kopeeksInRuble)
}

// Int returns amount in kopeks/cents.
func (m money) Int() int {
	return int(float64(m) * kopeeksInRuble)
}

// String returns float with amount representation.
func (m money) String() string {
	return fmt.Sprintf("%.2f", m)
}

func (m money) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.2f", float64(m))
	return []byte(s), nil
}

func (m *money) UnmarshalJSON(data []byte) error {
	const doublequotesign byte = 34 // "
	const moneybitsize = 64         // int64

	if data[0] == doublequotesign {
		data = data[1 : len(data)-1]
	}
	am, err := strconv.ParseFloat(string(data), moneybitsize)
	if err != nil {
		return fmt.Errorf("[QIWI] Amount JSON error: %w", ErrBadJSON)
	}

	*m = money(am)

	return err
}

// Amount carry money amount and currency ISO 3-ALPHA code.
type Amount struct {
	Value    money    `json:"value"`
	Currency currency `json:"currency"`
}

func newAmount(a int, cur currency) Amount {
	return Amount{Value: toMoney(a), Currency: cur}
}

// NewAmountInRubles set amount in Russian Rubles.
func NewAmountInRubles(a int) Amount {
	return newAmount(a, rub)
}

// NewAmountInDollars set amount in US Dollars.
func NewAmountInDollars(a int) Amount {
	return newAmount(a, usd)
}

// NewAmountInEuros set amount to amount in Euros.
func NewAmountInEuros(a int) Amount {
	return newAmount(a, eur)
}
