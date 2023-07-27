package wallet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ybalcin/wallet-service/pkg/utility"
	"sync"
	"time"
)

type (
	Wallet struct {
		sync.Mutex `bson:"-" json:"-"`
		ID         string    `bson:"_id" json:"id"`
		Username   string    `bson:"username" json:"username"`
		Balance    Money     `bson:"-" json:"balance"`
		CreatedAt  time.Time `bson:"created_at" json:"-"`

		Changes []Transaction `bson:"-" json:"-"`
	}

	Transaction struct {
		ID        string          `bson:"_id"`
		WalletID  string          `bson:"wallet_id"`
		Type      TransactionType `bson:"type"`
		Money     Money           `bson:"money"`
		CreatedAt time.Time       `bson:"created_at"`
	}
)

// New creates new instance of Wallet
func New(username string) (*Wallet, error) {
	if utility.IsStrEmpty(username) {
		return nil, errors.New(ErrInvalidUsername)
	}

	return &Wallet{
		ID:        uuid.NewString(),
		Username:  username,
		CreatedAt: time.Now(),
	}, nil
}

// WithdrawMoney withdraw(sub) money from wallet and adds transaction to the changes
func (w *Wallet) WithdrawMoney(money Money) error {
	if money.Amount == 0 {
		return errors.New(ErrInvalidMoneyAmount)
	}
	if w.Balance.Amount < money.Amount {
		return errors.New(ErrInsufficientMoneyAmount)
	}

	t := Transaction{
		ID:        uuid.NewString(),
		WalletID:  w.ID,
		Type:      WithdrawTransactionType,
		Money:     money,
		CreatedAt: time.Now(),
	}
	w.apply(t)

	return nil
}

// DepositMoney deposit(add) money to the wallet and adds transaction to the changes
func (w *Wallet) DepositMoney(money Money) error {
	if money.Amount == 0 {
		return errors.New(ErrInvalidMoneyAmount)
	}

	t := Transaction{
		ID:        uuid.NewString(),
		WalletID:  w.ID,
		Type:      DepositTransactionType,
		Money:     money,
		CreatedAt: time.Now(),
	}
	w.apply(t)

	return nil
}

// Mutate mutates wallet current state by transactions
func (w *Wallet) Mutate(transactions ...Transaction) {
	for _, transaction := range transactions {
		switch transaction.Type {
		case DepositTransactionType:
			w.Balance = add(w.Balance, transaction.Money)
		case WithdrawTransactionType:
			w.Balance = sub(w.Balance, transaction.Money)
		}
	}
}

func add(to Money, add Money) Money {
	to.Amount += add.Amount
	return to
}

func sub(from Money, sub Money) Money {
	from.Amount -= sub.Amount
	return from
}

func (w *Wallet) apply(transaction Transaction) {
	w.Lock()
	defer w.Unlock()

	w.Changes = append(w.Changes, transaction)
	w.Mutate(transaction)
}
