package main

import (
	"encoding/hex"
	"fmt"
	tx "transaction"
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
	hexStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}

	tx.ParseTransaction(decoded)

	opCode := tx.NewBitcoinOpCode()

	encoded := opCode.EncodeNum(-1)
	fmt.Printf("encode -1: %x\n", encoded)
	fmt.Printf("decode -1: %x\n", opCode.DecodeNum(encoded))
}
