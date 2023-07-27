package wallet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ybalcin/wallet-service/pkg/utility"
	"time"
)

type (
	Wallet struct {
		ID        string    `bson:"_id" json:"id"`
		Username  string    `bson:"username" json:"username"`
		Balance   Money     `bson:"-" json:"balance"`
		CreatedAt time.Time `bson:"created_at" json:"-"`

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
		return nil, errors.New("provide a valid username")
	}

	return &Wallet{
		ID:        uuid.NewString(),
		Username:  username,
		CreatedAt: time.Now(),
	}, nil
}

// NewWithdrawTransaction creates new instance of transaction with WithdrawTransactionType
func NewWithdrawTransaction(walletID string, money Money) Transaction {
	return Transaction{
		ID:        uuid.NewString(),
		WalletID:  walletID,
		Type:      WithdrawTransactionType,
		Money:     money,
		CreatedAt: time.Now(),
	}
}

// NewDepositTransaction creates new instance of transaction with DepositTransactionType
func NewDepositTransaction(walletID string, money Money) Transaction {
	return Transaction{
		ID:        uuid.NewString(),
		WalletID:  walletID,
		Type:      DepositTransactionType,
		Money:     money,
		CreatedAt: time.Now(),
	}
}

// DepositMoney provides deposit(add) money to the wallet
func (w *Wallet) DepositMoney(money Money) error {
	if money.Amount == 0 {
		return errors.New("provide valid money amount")
	}

	newBalance, err := NewMoney(w.Balance.Amount + money.Amount)
	if err != nil {
		return err
	}
	w.Balance = *newBalance

	return nil
}

// WithdrawMoney provides withdraw(remove) money from the wallet
func (w *Wallet) WithdrawMoney(money Money) error {
	if money.Amount == 0 {
		return errors.New("provide valid money amount")
	}
	if w.Balance.Amount < money.Amount {
		return errors.New("you don't have this amount of money in your wallet")
	}

	newBalance, err := NewMoney(w.Balance.Amount - money.Amount)
	if err != nil {
		return err
	}
	w.Balance = *newBalance

	return nil
}

// Mutate mutates wallet current state by transactions
func (w *Wallet) Mutate(transactions ...Transaction) error {
	for _, transaction := range transactions {
		switch transaction.Type {
		case DepositTransactionType:
			if err := w.DepositMoney(transaction.Money); err != nil {
				return err
			}
		case WithdrawTransactionType:
			if err := w.WithdrawMoney(transaction.Money); err != nil {
				return err
			}
		}
	}

	return nil
}

// Apply applies transactions to the wallet, adds transaction to the wallet.Changes
func (w *Wallet) Apply(transaction Transaction) error {
	w.Changes = append(w.Changes, transaction)
	return w.Mutate(transaction)
}
