package ecc

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

// 0x30 || total_length || 0x02 || r_length || r_bytes || 0x02 || s_length || s_bytes
func ParseSigBin(sigBin []byte) *Signature {
	reader := bytes.NewReader(sigBin)
	bufReader := bufio.NewReader(reader)

	firstByte := make([]byte, 1)
	bufReader.Read(firstByte)

	if firstByte[0] != 0x30 {
		panic("Bad signature")
	}

	totalLength := make([]byte, 1)
	bufReader.Read(totalLength)

	if totalLength[0]+2 != byte(len(sigBin)) {
		panic("Bad signature")
	}

	rMarker := make([]byte, 1)
	bufReader.Read(rMarker)
	if rMarker[0] != 0x02 {
		panic("Bad singature")
	}

	rLenBuf := make([]byte, 1)
	bufReader.Read(rLenBuf)
	rLen := rLenBuf[0]
	rBin := make([]byte, rLen)
	bufReader.Read(rBin)
	rBig := new(big.Int)
	rBig.SetBytes(rBin)

	sMarker := make([]byte, 1)
	bufReader.Read(sMarker)
	if sMarker[0] != 0x02 {
		panic("Bad signature")
	}

	sLenBuf := make([]byte, 1)
	bufReader.Read(sLenBuf)
	sLen := sLenBuf[0]
	sBin := make([]byte, sLen)
	bufReader.Read(sBin)
	sBig := new(big.Int)
	sBig.SetBytes(sBin)

	if int(sLen+rLen+6) != len(sigBin) {
		panic("Signature wrong")
	}

	n := BitcoinN()

	return &Signature{
		r: NewFieldElement(n, rBig),
		s: NewFieldElement(n, sBig),
	}
}

func Hash160(s []byte) []byte {
	sha256 := sha256.Sum256(s)
	hasher := ripemd160.New()
	hasher.Write(sha256[:])
	hashed := hasher.Sum(nil)
	return hashed
}

func Base58Checksum(s []byte) string {
	hash256 := Hash256(string(s))
	return EncodeBase58(append(s, hash256[:4]...))
}

func EncodeBase58(s []byte) string {
	const BASE58_ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	count := 0
	divFactor := big.NewInt(58)

	for idx := range s {
		if s[idx] == 0 {
			count += 1
		} else {
			break
		}
	}

	prefix := ""
	for i := 0; i < count; i++ {
		prefix += "1"
	}
	result := ""
	num := new(big.Int)
	num.SetBytes(s)

	for num.Cmp(big.NewInt(0)) > 0 {
		modRes := new(big.Int).Mod(num, divFactor)
		num = new(big.Int).Div(num, divFactor)

		result = string(BASE58_ALPHABET[modRes.Int64()]) + result
	}

	return prefix + result
}

func ParseSEC(secBin []byte) *Point {
	if secBin[0] == 4 {
		// uncompressed
		x := new(big.Int)
		x.SetBytes(secBin[1:33])
		y := new(big.Int)
		y.SetBytes(secBin[33:65])

		return S256Point(x, y)
	}

	isEven := (secBin[0] == 2)
	x := new(big.Int)
	x.SetBytes(secBin[1:])
	y2 := S256Field(x).Power(big.NewInt(3)).Add(S256Field(big.NewInt(7)))
	y := y2.Sqrt()
	var yEven *FieldElement
	var yOdd *FieldElement

	if new(big.Int).Mod(y.num, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		yEven = y
		yOdd = y.Negate()
	} else {
		yOdd = y
		yEven = y.Negate()
	}

	if isEven {
		return S256Point(x, yEven.num)
	} else {
		return S256Point(x, yOdd.num)
	}
}

// Hash it 2 times to reduce the risk
// This is Bitcoin's design
// Hash 1 time is reversible so we need to hash it twice
func Hash256(text string) []byte {
	hashOnce := sha256.Sum256([]byte(text))
	hashTwice := sha256.Sum256(hashOnce[:])
	return hashTwice[:]
}

func GeneratorPoint() *Point {
	Gx := new(big.Int)
	Gx.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	Gy := new(big.Int)
	Gy.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	G := S256Point(Gx, Gy)
	return G
}

func BitcoinN() *big.Int {
	n := new(big.Int)
	n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	return n
}
