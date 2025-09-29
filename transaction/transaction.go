package transaction

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
)

type Transaction struct {
	version  *big.Int
	txInputs []*TransactionInput
	txOutput []*TransactionOutput
	lockTime *big.Int
	testnet  bool
}

func getInputCount(bufReader *bufio.Reader) *big.Int {
	// if the first byte of input is 0, then witness transaction, we need to skip the first two bytes (0x00, 0x01)
	fitstByte, err := bufReader.Peek(1)
	if err != nil {
		panic(err)
	}

	if fitstByte[0] == 0x00 {
		skipBuf := make([]byte, 2)
		_, err = bufReader.Read(skipBuf)
		if err != nil {
			panic(err)
		}
	}

	count := ReadVarint(bufReader)
	return count
}

func ParseTransaction(binary []byte) *Transaction {
	transaction := Transaction{}
	reader := bytes.NewReader(binary)
	bufReader := bufio.NewReader(reader)

	versionBuf := make([]byte, 4)
	bufReader.Read(versionBuf)

	version := LittleEndianToBigInt(versionBuf, LITTLE_ENDIAN_4_BYTES)
	transaction.version = version
	fmt.Printf("Transaction version is: %d\n", version)

	inputCount := getInputCount(bufReader)
	fmt.Printf("Transaction input count is: %d\n", inputCount)

	inputs := []*TransactionInput{}
	for i := 0; i < int(inputCount.Int64()); i++ {
		input := NewTransactionInput(bufReader)
		inputs = append(inputs, input)
	}
	transaction.txInputs = inputs

	return &transaction
}
