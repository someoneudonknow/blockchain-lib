package main

import (
	"ecc"
	"fmt"
	"math/big"
)

/*
	Bitcoin parameters
	a = 0
	b = 7
	Gx = 0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798
	Gy = 0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8
	p = 2^256 - 2^32 - 977
	n = 0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141
	y^2 = x^3 + 7
*/

func main() {
	// e := big.NewInt(12345)
	// message := new(big.Int)
	// message.SetBytes(ecc.Hash256("TuTran transfer 1 BTC to TuTran"))
	//
	// pk := ecc.NewPrivateKey(e)
	//
	// sig := pk.Sign(message)
	//
	// fmt.Printf("Sign message is %s\n", sig)
	//
	// pubKey := pk.Public()
	// n := ecc.BitcoinN()
	//
	// fmt.Printf("SEC of pubKey is: %s\n", pubKey.SEC(true))
	//
	// messageFieldElement := ecc.NewFieldElement(n, message)
	//
	// if pubKey.Verify(messageFieldElement, sig) {
	// 	fmt.Println("The signature is valid")
	// }
	//
	// secBinUnCompressed := new(big.Int)
	// secBinUnCompressed.SetString(pubKey.SEC(false), 16)
	// unUnCompressedDecode := ecc.ParseSEC(secBinUnCompressed.Bytes())
	//
	// fmt.Printf("Parse SEC of pubKey is: %s\n", unUnCompressedDecode)

	// r := new(big.Int)
	// r.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	// rField := ecc.S256Field(r)
	// s := new(big.Int)
	// s.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	// sField := ecc.S256Field(s)
	// sig := ecc.NewSignature(rField, sField)
	// derEncode := sig.DER()
	// fmt.Printf("der encoding for signature is %x\n", derEncode)

	// val := new(big.Int)
	// val.SetString("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", 16)
	// fmt.Printf("base58 encoding is %s\n", ecc.EncodeBase58(val.Bytes()))

	// privateKey := ecc.NewPrivateKey(big.NewInt(5002))
	// pubKey := privateKey.Public()
	// fmt.Printf("Wallet address for 5002*G is %s\n", pubKey.Address(false, true))

	p := new(big.Int)
	p.SetString("12345678", 16)
	bytes := p.Bytes()
	fmt.Printf("bytes for 0x12345678 is %x\n", bytes)

	littleEndianByte := ecc.BigIntToLittleEndian(p, ecc.LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("little endian for 0x12345678 is %x\n", littleEndianByte)

	littleEndianByteToInt64 := ecc.LittleEndianToBigInt(littleEndianByte, ecc.LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("little endian bytes into int is %x\n", littleEndianByteToInt64)
}
