package number

import "service/internal/pkg/core"

type RentalRedisNumber struct{}

func (n *RentalRedisNumber) GenerateRentalNumber() string {
	return n.generate("RNT")
}

func (n *RentalRedisNumber) GenerateRentalPaymentNumber() string {
	return n.generate("PAY")
}

func (n *RentalRedisNumber) GenerateRentalRefundNumber() string {
	return n.generate("REF")
}

func (n *RentalRedisNumber) generate(prefix string) string {
	return core.RedisNumber(prefix)
}
