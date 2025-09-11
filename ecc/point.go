package ecc

import (
	"fmt"
	"math/big"
)

type OP_TYPE int

const (
	ADD OP_TYPE = iota
	SUB
	MUL
	DIV
	EXP
)

type Point struct {
	a *FieldElement
	b *FieldElement
	x *FieldElement
	y *FieldElement
}

func OpOnField(x *FieldElement, y *FieldElement, scalar *big.Int, opType OP_TYPE) *FieldElement {
	switch opType {
	case ADD:
		return x.Add(y)
	case SUB:
		return x.Substract(y)
	case MUL:
		if y != nil {
			return x.Multiply(y)
		}
		if scalar != nil {
			return x.ScalarMul(scalar)
		}
		panic("Error on multiply")
	case DIV:
		return x.Divide(y)
	case EXP:
		if scalar == nil {
			panic("Scalar can't be nil on EXP operation")
		}
		return x.Power(scalar)
	default:
		panic("Operation type not supported!")
	}
}

// y^2 = x^3 + ax + b
func NewEllipticCurvePoint(
	x *FieldElement,
	y *FieldElement,
	a *FieldElement,
	b *FieldElement,
) *Point {
	// x = nil and y = nil => Identity point
	if x == nil && y == nil {
		return &Point{
			a: a,
			b: b,
			x: x,
			y: y,
		}
	}

	left := OpOnField(y, nil, big.NewInt(2), EXP)
	x3 := OpOnField(x, nil, big.NewInt(3), EXP)
	ax := OpOnField(a, x, nil, MUL)
	right := OpOnField(OpOnField(x3, ax, nil, ADD), b, nil, ADD)

	if left.EqualTo(right) != true {
		err := fmt.Sprintf("Point(%v, %v) is not on the curve with a:%v, b:%v\n", x, y, a, b)
		panic(err)
	}

	return &Point{
		a: a,
		b: b,
		x: x,
		y: y,
	}
}

func S256Point(x *big.Int, y *big.Int) *Point {
	a := S256Field(big.NewInt(0))
	b := S256Field(big.NewInt(7))

	if x == nil && y == nil {
		return &Point{
			a: a,
			b: b,
			x: nil,
			y: nil,
		}
	}

	return &Point{
		a: a,
		b: b,
		x: S256Field(x),
		y: S256Field(y),
	}
}

func (p *Point) Address(compressed bool, testnet bool) string {
	hash160 := p.hash160(compressed)
	prefix := []byte{}

	if testnet {
		prefix = append(prefix, 0x6f)
	} else {
		prefix = append(prefix, 0x00)
	}

	return Base58Checksum(append(prefix, hash160...))
}

func (p *Point) hash160(compressed bool) []byte {
	_, secBytes := p.SEC(compressed)
	return Hash160(secBytes)
}

// uncompressed = 0x04 + x (32 bytes) + y (32 bytes)
func (p *Point) SEC(compressed bool) (string, []byte) {
	secBytes := []byte{}
	if !compressed {
		secBytes = append(secBytes, 0x04)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		secBytes = append(secBytes, p.y.num.Bytes()...)

		return fmt.Sprintf("04%064x%064x", p.x.num, p.y.num), secBytes
	}
	if new(big.Int).Mod(p.y.num, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		secBytes = append(secBytes, 0x02)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		return fmt.Sprintf("02%064x", p.x.num), secBytes
	} else {
		secBytes = append(secBytes, 0x03)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		return fmt.Sprintf("03%064x", p.x.num), secBytes
	}
}

func (p *Point) Verify(e *FieldElement, sig *Signature) bool {
	/*
		e = hash(message)
		u1 = e * s^-1
		u2 = r * s^-1
		P = u1*G + u2*Q(current Point) => (xP, yP)
		r = xP -> Verify success
	*/
	sInverse := sig.s.Inverse()

	u1 := OpOnField(e, sInverse, nil, MUL)
	u2 := OpOnField(sig.r, sInverse, nil, MUL)
	G := GeneratorPoint()
	total := (G.ScalarMul(u1.num)).Add(p.ScalarMul(u2.num))

	return total.x.num.Cmp(sig.r.num) == 0
}

/*
k*G, k = 13, => 13G
k = 1101 (2^3 + 2^2 + 2^0) * G => 2^3G + 2^2G + 2^0G
=> (G << 3) + (G << 2) + (G << 0)
1 trillition, 40 bits in binary form
we at most do 40 times of addition => 1 trilliion times
*/

func (p *Point) ScalarMul(scalar *big.Int) *Point {
	if scalar == nil {
		panic("Scalar can't be nil")
	}

	binaryForm := fmt.Sprintf("%b", scalar)
	result := NewEllipticCurvePoint(nil, nil, p.a, p.b)
	current := p

	for i := len(binaryForm) - 1; i >= 0; i-- {
		if binaryForm[i] == '1' {
			result = result.Add(current)
		}
		// doubling step 2G -> 4G -> 8G
		current = current.Add(current)
	}

	return result
}

func (p *Point) Add(other *Point) *Point {
	// Check if two points are on the same curve, a and b are constants so if the a and b of 2 point is different, its the two different curve. => Can't perform add operation
	if !p.a.EqualTo(other.a) || !p.b.EqualTo(other.b) {
		panic("Given two point not on the same curve")
	}

	if p.x == nil {
		return other
	}

	if other.x == nil {
		return p
	}

	// points are on the vertical A(x, y), B(x, -y)
	zero := NewFieldElement(p.y.order, big.NewInt(0))
	if p.x.EqualTo(other.x) && OpOnField(p.y, other.y, nil, ADD).EqualTo(zero) {
		return &Point{
			x: nil,
			y: nil,
			a: p.a,
			b: p.b,
		}
	}

	slope := p.SlopeTo(other)
	slopeSquared := OpOnField(slope, nil, big.NewInt(2), EXP)

	// x3 = m^2 - x1 - x2
	x3 := OpOnField(OpOnField(slopeSquared, p.x, nil, SUB), other.x, nil, SUB)
	// y3 = m(x3 - x1) + y1
	y3 := OpOnField(OpOnField(slope, OpOnField(x3, p.x, nil, SUB), nil, MUL), p.y, nil, ADD)
	minusY3 := OpOnField(y3, nil, big.NewInt(-1), MUL)

	return &Point{
		x: x3,
		y: minusY3,
		a: p.a,
		b: p.b,
	}
}

func (p *Point) SlopeTo(other *Point) *FieldElement {
	var numerator *FieldElement
	var denominator *FieldElement

	if p.Equal(other) {
		numerator = OpOnField(
			OpOnField(OpOnField(p.x, nil, big.NewInt(2), EXP), nil, big.NewInt(3), MUL),
			p.a,
			nil,
			ADD,
		)
		denominator = OpOnField(p.y, nil, big.NewInt(2), MUL)
	} else {
		numerator = OpOnField(other.y, p.y, nil, SUB)
		denominator = OpOnField(other.x, p.x, nil, SUB)
	}

	return OpOnField(numerator, denominator, nil, DIV)
}

func (p *Point) String() string {
	xString := "nil"
	yString := "nil"

	if p.x != nil {
		xString = p.x.String()
	}

	if p.y != nil {
		yString = p.y.String()
	}

	return fmt.Sprintf(
		"(x:%s, y:%s, a:%s, b:%s)",
		xString,
		yString,
		p.a.String(),
		p.b.String(),
	)
}

func (p *Point) Equal(other *Point) bool {
	return p.a.EqualTo(other.a) && p.b.EqualTo(other.b) && p.x.EqualTo(other.x) &&
		p.y.EqualTo(other.y)
}

func (p *Point) NotEqual(other *Point) bool {
	return !p.Equal(other)
}
