package mathutils

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

// DecimalAdd 加法
func DecimalAdd(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Add(d2)
}

// DecimalSub 减法
func DecimalSub(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Sub(d2)
}

// DecimalMul 乘法
func DecimalMul(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Mul(d2)
}

// DecimalDiv 除法
func DecimalDiv(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Div(d2)
}

// DecimalToInt64 int
func DecimalToInt64(d decimal.Decimal) int64 {
	return d.IntPart()
}

func DecimalString(value string) decimal.Decimal {
	price, _ := decimal.NewFromString(value)
	return price
}

func DecimalInt64(value int64) decimal.Decimal {
	price, _ := decimal.NewFromString(cast.ToString(value))
	return price
}
func DecimalUInt64(value uint64) decimal.Decimal {
	price, _ := decimal.NewFromString(cast.ToString(value))
	return price
}

func DecimalInt(value int) decimal.Decimal {
	price, _ := decimal.NewFromString(cast.ToString(value))
	return price
}

// DecimalFloat64 float
func DecimalFloat64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

// StringFloat64 float
func StringFloat64(d string) float64 {
	f, _ := DecimalString(d).Float64()
	return f
}

func Fen2Yuan(price int64) string {
	d := decimal.New(1, 2)
	return DecimalInt64(price).DivRound(d, 2).String()
}
