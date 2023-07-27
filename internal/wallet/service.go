package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/ybalcin/wallet-service/pkg/errr"
	"github.com/ybalcin/wallet-service/pkg/utility"
)

type (
	// Service is an interface for wallet use cases
	Service interface {
		// CreateWallet creates wallet
		CreateWallet(ctx context.Context, dto *CreateWalletRequest) (*ID, *errr.Error)
		// DepositMoney provides to deposit(add) money to wallet
		DepositMoney(ctx context.Context, id string, req *MoneyTransactionRequest) (*Wallet, *errr.Error)
		// WithdrawMoney provides to withdraw(remove) monet from wallet
		WithdrawMoney(ctx context.Context, id string, req *MoneyTransactionRequest) (*Wallet, *errr.Error)
	}

	// ServiceImplementation is an implementation of Service interface
	ServiceImplementation struct {
		repository Repository
	}
)

// NewService creates new instance of ServiceImplementation
func NewService(repository Repository) *ServiceImplementation {
	return &ServiceImplementation{repository: repository}
}

// CreateWallet creates wallet
func (s *ServiceImplementation) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*ID, *errr.Error) {
	wallet, err := New(req.Username)
	if err != nil {
		return nil, errr.ThrowBadRequestError(err)
	}

	if err = s.repository.InsertWallet(ctx, wallet); err != nil {
		return nil, errr.ThrowInternalServerError(err)
	}

	return NewID(wallet.ID), nil
}

// DepositMoney provides to deposit(add) money to wallet
func (s *ServiceImplementation) DepositMoney(ctx context.Context, walletID string, req *MoneyTransactionRequest) (*Wallet, *errr.Error) {
	money, err := NewMoney(req.Amount)
	if err != nil {
		return nil, errr.ThrowBadRequestError(err)
	}

	return s.applyTransactionToWallet(ctx, NewDepositTransaction(walletID, *money))
}

// WithdrawMoney provides to withdraw(remove) monet from wallet
func (s *ServiceImplementation) WithdrawMoney(ctx context.Context, walletID string, req *MoneyTransactionRequest) (*Wallet, *errr.Error) {
	money, err := NewMoney(req.Amount)
	if err != nil {
		return nil, errr.ThrowBadRequestError(err)
	}

	return s.applyTransactionToWallet(ctx, NewWithdrawTransaction(walletID, *money))
}

func (s *ServiceImplementation) findWalletWithCurrentState(ctx context.Context, walletID string) (*Wallet, *errr.Error) {
	if utility.IsStrEmpty(walletID) {
		return nil, errr.ThrowBadRequestError(errors.New("provide a valid wallet id"))
	}

	wallet, err := s.repository.FindWalletByID(ctx, walletID)
	if err != nil {
		return nil, errr.ThrowInternalServerError(err)
	}
	if wallet == nil {
		return nil, errr.ThrowNotFoundError(fmt.Errorf("wallet with id: %s not found", walletID))
	}

	transactions, err := s.repository.FindTransactionsByWalletID(ctx, wallet.ID)
	if err != nil {
		return nil, errr.ThrowInternalServerError(err)
	}
	if err = wallet.Mutate(transactions...); err != nil {
		return nil, errr.ThrowBadRequestError(err)
	}

	return wallet, nil
}

func (s *ServiceImplementation) applyTransactionToWallet(ctx context.Context, transaction Transaction) (*Wallet, *errr.Error) {
	wallet, err := s.findWalletWithCurrentState(ctx, transaction.WalletID)
	if err != nil {
		return nil, err
	}

	if ex := wallet.Apply(transaction); ex != nil {
		return nil, errr.ThrowBadRequestError(ex)
	}
	if len(wallet.Changes) > 0 {
		if ex := s.repository.InsertTransactions(ctx, wallet.Changes...); ex != nil {
			return nil, errr.ThrowInternalServerError(ex)
		}
	}

	return wallet, nil
}
