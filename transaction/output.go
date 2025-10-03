package transaction

import (
	"bufio"
	"math/big"
)

type TransactionOutput struct {
	amount       *big.Int
	scriptPubKey *ScriptSig
	scriptLength *big.Int
}

func NewTransactionOutput(reader *bufio.Reader) *TransactionOutput {
	return &TransactionOutput{
		amount:       big.NewInt(0),
		scriptPubKey: nil,
		scriptLength: nil,
	}
}
