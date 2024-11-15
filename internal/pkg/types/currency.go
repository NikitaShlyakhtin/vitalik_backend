package types

import "fmt"

type Currency string

const (
	BTC  Currency = "BTC"
	USDT Currency = "USDT"
	ETH  Currency = "ETH"
)

var availableCurrencies = []Currency{BTC, USDT, ETH}

func (c *Currency) String() string {
	return string(*c)
}

func (c *Currency) Validate() bool {
	switch *c {
	case BTC, USDT, ETH:
		return true
	default:
		return false
	}
}

type CurrencyPair struct {
	Currency1 Currency `json:"currency1"`
	Currency2 Currency `json:"currency2"`
}

func (p *CurrencyPair) String() string {
	if p == nil {
		return ""
	}

	return fmt.Sprintf("%s/%s", p.Currency1.String(), p.Currency2.String())
}

func (p *CurrencyPair) StringReverse() string {
	if p == nil {
		return ""
	}

	return fmt.Sprintf("%s/%s", p.Currency2.String(), p.Currency1.String())
}

func (p *CurrencyPair) Equals(p2 *CurrencyPair) bool {
	if p == nil || p2 == nil {
		return false
	}

	return (p.Currency1 == p2.Currency1 && p.Currency2 == p2.Currency2) ||
		(p.Currency1 == p2.Currency2 && p.Currency2 == p2.Currency1)
}

func (p *CurrencyPair) Validate() bool {
	if p == nil {
		return false
	}

	return p.Currency1.Validate() && p.Currency2.Validate()
}

func ListAvailableCurrencies() []Currency {
	return availableCurrencies
}
