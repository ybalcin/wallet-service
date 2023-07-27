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
		// WithdrawMoney provides to withdraw(sub) monet from wallet
		WithdrawMoney(ctx context.Context, id string, req *MoneyTransactionRequest) (*Wallet, *errr.Error)
		// GetWallet gets wallet with current state
		GetWallet(ctx context.Context, walletID string) (*Wallet, *errr.Error)
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

func (s *ServiceImplementation) findWalletWithCurrentState(ctx context.Context, walletID string) (*Wallet, *errr.Error) {
	if utility.IsStrEmpty(walletID) {
		return nil, errr.ThrowBadRequestError(errors.New(ErrInvalidWalletID))
	}

	wallet, err := s.repository.FindWalletByID(ctx, walletID)
	if err != nil {
		return nil, errr.ThrowInternalServerError(err)
	}
	if wallet == nil {
		return nil, errr.ThrowNotFoundError(fmt.Errorf(ErrWalletNotFound, walletID))
	}

	transactions, err := s.repository.FindTransactionsByWalletID(ctx, wallet.ID)
	if err != nil {
		return nil, errr.ThrowInternalServerError(err)
	}
	wallet.Mutate(transactions...)

	return wallet, nil
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

	wallet, ex := s.findWalletWithCurrentState(ctx, walletID)
	if ex != nil {
		return nil, ex
	}

	if err = wallet.DepositMoney(*money); err != nil {
		return nil, errr.ThrowBadRequestError(err)
	}

	return s.saveWalletChanges(ctx, wallet)
}

// WithdrawMoney provides to withdraw(sub) monet from wallet
func (s *ServiceImplementation) WithdrawMoney(ctx context.Context, walletID string, req *MoneyTransactionRequest) (*Wallet, *errr.Error) {
	money, err := NewMoney(req.Amount)
	if err != nil {
		return nil, errr.ThrowBadRequestError(err)
	}

	wallet, ex := s.findWalletWithCurrentState(ctx, walletID)
	if ex != nil {
		return nil, ex
	}

	if err = wallet.WithdrawMoney(*money); err != nil {
		return nil, errr.ThrowBadRequestError(err)
	}

	return s.saveWalletChanges(ctx, wallet)
}

// GetWallet gets wallet with current state
func (s *ServiceImplementation) GetWallet(ctx context.Context, walletID string) (*Wallet, *errr.Error) {
	return s.findWalletWithCurrentState(ctx, walletID)
}

func (s *ServiceImplementation) saveWalletChanges(ctx context.Context, wallet *Wallet) (*Wallet, *errr.Error) {
	if len(wallet.Changes) > 0 {
		if err := s.repository.InsertTransactions(ctx, wallet.Changes...); err != nil {
			return nil, errr.ThrowInternalServerError(err)
		}
	}

	return wallet, nil
}
