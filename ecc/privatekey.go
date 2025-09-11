package ecc

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	d *big.Int
	Q *Point
}

func NewPrivateKey(secret *big.Int) *PrivateKey {
	G := GeneratorPoint()

	return &PrivateKey{
		d: secret,
		Q: G.ScalarMul(secret),
	}
}

func (pk *PrivateKey) Sign(e *big.Int) *Signature {
	/*
				All calculation on finite field element
				choose k randomly from 1 -> n-1
				R = kG
				r = xR; if r = 0; choose different k
				s = k^-1 x (e (hashed message) + d (private key) x r); if s = 0; choose different k
		    Signature{r, s}
	*/
	n := BitcoinN()
	G := GeneratorPoint()

	for {
		k, err := rand.Int(rand.Reader, n)
		if err != nil {
			panic("Error occur while generating private key")
		}

		r := G.ScalarMul(k).x.num
		if r.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		kField := NewFieldElement(n, k)
		eField := NewFieldElement(n, e)
		rField := NewFieldElement(n, r)
		dField := NewFieldElement(n, pk.d)
		kInverse := kField.Inverse()
		dxr := dField.Multiply(rField)
		ePlusDxr := eField.Add(dxr)
		s := kInverse.Multiply(ePlusDxr)

		if s.num.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		/*
		   if s > n / 2 we need to change it to n - s, when doing signature verify, s and n - s are equivalance doing this change is for malleability reason
		*/
		if s.num.Cmp(new(big.Int).Div(n, big.NewInt(2))) > 0 {
			s = NewFieldElement(n, new(big.Int).Sub(n, s.num))
		}

		return &Signature{
			r: rField,
			s: s,
		}
	}
}

func (pk *PrivateKey) Public() *Point {
	return pk.Q
}

func (pk *PrivateKey) String() string {
	return fmt.Sprintf("Private key hex:{%s}", pk.d)
}
