package wallet

import "errors"

type Money struct {
	Amount float32 `json:"amount" bson:"amount"`
}

// NewMoney creates new instance of money
func NewMoney(amount float32) (*Money, error) {
	if amount < 0 {
		return nil, errors.New(ErrInvalidMoneyAmount)
	}

	return &Money{Amount: amount}, nil
}

// Add adds amount to `to` arg
func (m Money) Add(to *Money) {
	to.Amount += m.Amount
}

// Sub subtract amount from `from` arg
func (m Money) Sub(from *Money) {
	from.Amount -= m.Amount
}

type ID struct {
	Id string `json:"id"`
}

func NewID(id string) *ID {
	return &ID{id}
}
