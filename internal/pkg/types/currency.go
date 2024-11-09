package types

type Currency string

const (
	BTC  Currency = "BTC"
	USDT Currency = "USDT"
	ETH  Currency = "ETH"
)

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
