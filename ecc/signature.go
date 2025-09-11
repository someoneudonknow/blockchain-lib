package ecc

import "fmt"

type Signature struct {
	r *FieldElement
	s *FieldElement
}

func NewSignature(r *FieldElement, s *FieldElement) *Signature {
	return &Signature{r, s}
}

/*
1. Set the first byte to 0x30
2. Second byte is the total length of s and r
3. The first byte is 0x02 it is indicator for the beginning of the byte array for r
4. if the first byte of r is >= 0x80, then we need to append 0x00 as beginning byte of the bytes array of r, compute the length of the bytes array of rand append the length behind the 0x02 of step 2
5. Insert 0x02 behide the last byte of the r byte array, as indicator for the beginning
6. Do the same for s as step 4
total length of 0x44 or 0x45
*/
func (s *Signature) DER() []byte {
	rBytes := s.r.num.Bytes()
	sBytes := s.s.num.Bytes()

	if len(rBytes) == 0 {
		panic("r is empty")
	}

	if len(sBytes) == 0 {
		panic("s is empty")
	}

	if rBytes[0] >= byte(0x80) {
		rBytes = append([]byte{0x00}, rBytes...)
	}

	if sBytes[0] >= byte(0x80) {
		sBytes = append([]byte{0x00}, sBytes...)
	}

	rBytes = append([]byte{0x02, byte(len(rBytes))}, rBytes...)
	sBytes = append([]byte{0x02, byte(len(sBytes))}, sBytes...)

	encoded := []byte{0x30}
	encoded = append(encoded, byte(len(sBytes)+len(rBytes)))

	encoded = append(encoded, rBytes...)
	encoded = append(encoded, sBytes...)

	return encoded
}

func (s *Signature) String() string {
	return fmt.Sprintf("Signature(r: {%s}, s: {%s})", s.r.String(), s.s.String())
}
