package transaction

import (
	"bufio"
	"math/big"
)

type TransactionInput struct {
	preTxID   []byte
	preTxIdx  *big.Int
	scriptSig *ScriptSig
	sequence  *big.Int
}

func NewTransactionInput(reader *bufio.Reader) *TransactionInput {
	// First 32 bytes are previous hash256 of previous transaction
	transactionInput := &TransactionInput{}
	previousTx := make([]byte, 32)

	reader.Read(previousTx)
	// reverse the byte to convert from little endian to big endian
	transactionInput.preTxID = reverseByteSlice(previousTx)

	// next 4 bytes is the previous transaction index in little endian
	preTxIdx := make([]byte, 4)
	reader.Read(preTxIdx)
	transactionInput.preTxIdx = LittleEndianToBigInt(preTxIdx, LITTLE_ENDIAN_4_BYTES)

	// next is scriptSig
	transactionInput.scriptSig = NewScriptSig(reader)

	sequence := make([]byte, 4)
	reader.Read(sequence)
	transactionInput.sequence = LittleEndianToBigInt(sequence, LITTLE_ENDIAN_4_BYTES)

	return transactionInput
}

func reverseByteSlice(bytes []byte) []byte {
	reversed := []byte{}
	for i := len(bytes) - 1; i >= 0; i-- {
		reversed = append(reversed, reversed[i])
	}
	return reversed
}
