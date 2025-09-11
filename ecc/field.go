package ecc

import (
	"fmt"
	"math/big"
)

/*
	3
	{0, 1, 2}
	3 + 0 = 0 (mod 3)
	2 + 1 = 0 (mod 3)
	negate = 3 - 3
	negate = 3 - 2 = 1
	{0, 1, 2}
	2 - 1 = 1
	subtract = 1 - 2
	1 + 1 = 2
	3 + 1
	a.negate() = -a
*/

type FieldElement struct {
	order *big.Int
	num   *big.Int
}

func S256Field(num *big.Int) *FieldElement {
	twoExp256 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	twoExp32 := new(big.Int).Exp(big.NewInt(2), big.NewInt(32), nil)
	p := new(big.Int).Sub(new(big.Int).Sub(twoExp256, twoExp32), big.NewInt(977))
	return &FieldElement{p, num}
}

func NewFieldElement(order *big.Int, num *big.Int) *FieldElement {
	if order.Cmp(num) == -1 {
		err := fmt.Sprintf("Num not in the range of 0 to %v", order)
		panic(err)
	}
	return &FieldElement{order, num}
}

// PUBLIC METHODS

func (fe *FieldElement) Sqrt() *FieldElement {
	// make sure (p + 1) % 4 == 0
	orderAddOne := new(big.Int).Add(fe.order, big.NewInt(1))
	modRes := new(big.Int).Mod(orderAddOne, big.NewInt(4))

	if modRes.Cmp(big.NewInt(0)) != 0 {
		panic("Order plus one mod 4 is not 0")
	}

	return fe.Power(new(big.Int).Div(orderAddOne, big.NewInt(4)))
}

func (fe *FieldElement) Divide(other *FieldElement) *FieldElement {
	fe.checkOrder(other)

	op := new(big.Int)
	otherReverse := other.Power(op.Sub(fe.order, big.NewInt(int64(2))))

	return fe.Multiply(otherReverse)
}

func (fe *FieldElement) Inverse() *FieldElement {
	return fe.Power(new(big.Int).Sub(fe.order, big.NewInt(int64(2))))
}

func (fe *FieldElement) ScalarMul(val *big.Int) *FieldElement {
	var op big.Int
	return NewFieldElement(fe.order, op.Mod(op.Mul(fe.num, val), fe.order))
}

func (fe *FieldElement) Power(power *big.Int) *FieldElement {
	var op big.Int
	t := op.Mod(power, op.Sub(fe.order, big.NewInt(int64(1))))
	res := op.Exp(fe.num, t, fe.order)
	return NewFieldElement(fe.order, res)
}

func (fe *FieldElement) Multiply(other *FieldElement) *FieldElement {
	fe.checkOrder(other)
	var op big.Int

	return NewFieldElement(
		fe.order,
		op.Mod(op.Mul(fe.num, other.num), fe.order),
	)
}

func (fe *FieldElement) Substract(other *FieldElement) *FieldElement {
	return fe.Add(other.Negate())
}

func (fe *FieldElement) Add(other *FieldElement) *FieldElement {
	fe.checkOrder(other)
	var op big.Int
	return NewFieldElement(fe.order, op.Mod(op.Add(fe.num, other.num), fe.order))
}

func (fe *FieldElement) Negate() *FieldElement {
	var op big.Int
	return NewFieldElement(fe.order, op.Sub(fe.order, fe.num))
}

func (fe *FieldElement) String() string {
	return fmt.Sprintf("FieldElement{order: %s, num: %s}", fe.order.String(), fe.num.String())
}

func (fe *FieldElement) EqualTo(other *FieldElement) bool {
	return fe.num.Cmp(other.num) == 0 && other.order.Cmp(other.order) == 0
}

// PRIVATE METHODS

func (fe *FieldElement) checkOrder(other *FieldElement) {
	if fe.order.Cmp(other.order) != 0 {
		panic("Order are not equal!")
	}
}
