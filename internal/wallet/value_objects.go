package wallet

import "errors"

type Money struct {
	Amount float32 `json:"amount"`
}

func NewMoney(amount float32) (*Money, error) {
	if amount < 0 {
		return nil, errors.New(ErrInvalidMoneyAmount)
	}

	return &Money{Amount: amount}, nil
}

type ID struct {
	Id string `json:"id"`
}

func NewID(id string) *ID {
	return &ID{id}
}
