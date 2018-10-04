package chainlib

import (
	"time"
)

type Transaction struct {
	Category string `json:"category"`

	TransactionID string `json:"transactionid"`

	From string `json:"from"`

	To string `json:"to"`

	Amount float64 `json:"amount"`

	Time time.Time `json:"time"`

	BlockNum int64 `json:"blocknum"`

	IsIrrevisible bool `json:"isirreversible"`

	Memo string `json:"memo"`
}

type ExtractTransaction struct {
	To string `json:"to"`

	Value string `json:"value"`

	Extras string `json:"extras"`

	TrxID string `json:"trxid"`
}
