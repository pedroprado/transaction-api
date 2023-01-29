package values

import "github.com/pkg/errors"

type TransactionType string

const (
	TransactionTypePixOut TransactionType = "PIX_OUT"
)

var ValidTransactionTypes = map[TransactionType]TransactionType{
	TransactionTypePixOut: TransactionTypePixOut,
}

func (transactionType TransactionType) IsValid() (bool, error) {
	if _, exists := ValidTransactionTypes[transactionType]; !exists {
		return false, errors.New("not valid transaction type")
	}
	return true, nil
}
