package tinvest

import (
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/shopspring/decimal"
)

type unitsNano interface {
	GetUnits() int64
	GetNano() int32
}

func ToDecimal(q unitsNano) decimal.Decimal {
	if q == nil {
		return decimal.Zero
	}
	return decimal.NewFromInt(q.GetUnits()).Add(decimal.New(int64(q.GetNano()), -9))
}

func FromDecimal(d decimal.Decimal) *investapi.Quotation {
	units := d.Truncate(0).IntPart()

	fractionalPart := d.Sub(decimal.NewFromInt(units))
	nanoDecimal := fractionalPart.Mul(decimal.New(1, 9))
	nano := int32(nanoDecimal.IntPart()) // nolint: gosec

	return &investapi.Quotation{
		Units: units,
		Nano:  nano,
	}
}
