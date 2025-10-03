package transaction

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type TransactionFetcher struct {
}

func NewTransactionFetcher() *TransactionFetcher {
	return &TransactionFetcher{}
}

func (tf *TransactionFetcher) getURL(testnet bool) string {
	if testnet {
		return "https://blockstream.info/testnet/api/tx"
	}
	return "https://blockstream.info/api/tx"
}

func (tf *TransactionFetcher) Fetch(txid string, testnet bool) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/hex", tf.getURL(testnet), txid)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	buf, err := hex.DecodeString(string(body))
	if err != nil {
		return nil, err
	}

	return buf, nil
}
