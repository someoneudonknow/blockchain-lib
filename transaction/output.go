package transaction

import (
	"bufio"
)

type TransactionOutput struct {
	reader *bufio.Reader
}

func NewTransactionOutput(reader *bufio.Reader) *TransactionInput {
	return &TransactionInput{
		reader,
	}
}
