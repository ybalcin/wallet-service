package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ybalcin/wallet-service/internal/wallet"
	"github.com/ybalcin/wallet-service/pkg/errr"
	"go.uber.org/mock/gomock"
	"testing"
)

func setupMockRepo(t *testing.T) *MockRepository {
	return NewMockRepository(gomock.NewController(t))
}

func TestServiceImplementation(t *testing.T) {
	ctx := context.Background()
	mockRepo := setupMockRepo(t)
	service := wallet.NewService(mockRepo)

	t.Run("CreateWallet", func(t *testing.T) {
		req := &wallet.CreateWalletRequest{Username: "user"}

		mockRepo.EXPECT().InsertWallet(ctx, gomock.Any()).Return(nil)

		id, err := service.CreateWallet(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, id)
	})

	t.Run("DepositMoney", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			req := &wallet.MoneyTransactionRequest{Amount: 10}
			w := &wallet.Wallet{
				ID:      uuid.NewString(),
				Balance: wallet.Money{Amount: 10},
			}

			mockRepo.EXPECT().FindWalletByID(ctx, w.ID).Return(w, nil)
			mockRepo.EXPECT().FindTransactionsByWalletID(ctx, w.ID).Return(nil, nil)
			mockRepo.EXPECT().InsertTransactions(ctx, gomock.Any()).Return(nil)

			expected, err := service.DepositMoney(ctx, w.ID, req)
			assert.Nil(t, err)
			assert.Equal(t, expected, w)
		})

		t.Run("should return ErrInvalidWalletID if walletID is empty", func(t *testing.T) {
			req := &wallet.MoneyTransactionRequest{Amount: 10}
			w := &wallet.Wallet{
				ID:      "",
				Balance: wallet.Money{Amount: 10},
			}

			w, err := service.DepositMoney(ctx, w.ID, req)
			assert.Nil(t, w)
			assert.Equal(t, errr.ThrowBadRequestError(errors.New(wallet.ErrInvalidWalletID)), err)
		})

		t.Run("should return error if repository.FindWalletByID returns error", func(t *testing.T) {
			req := &wallet.MoneyTransactionRequest{Amount: 10}
			w := &wallet.Wallet{
				ID:      uuid.NewString(),
				Balance: wallet.Money{Amount: 10},
			}

			e := errors.New("")

			mockRepo.EXPECT().FindWalletByID(ctx, w.ID).Return(w, e)

			w, err := service.DepositMoney(ctx, w.ID, req)
			assert.Nil(t, w)
			assert.Equal(t, errr.ThrowInternalServerError(e), err)
		})

		t.Run("should return error if wallet not found", func(t *testing.T) {
			req := &wallet.MoneyTransactionRequest{Amount: 10}

			id := uuid.NewString()

			mockRepo.EXPECT().FindWalletByID(ctx, id).Return(nil, nil)

			w, err := service.DepositMoney(ctx, id, req)
			assert.Nil(t, w)
			assert.Equal(t, errr.ThrowNotFoundError(fmt.Errorf(wallet.ErrWalletNotFound, id)), err)
		})

		t.Run("should return error if repository.FindTransactionsByWalletID returns error", func(t *testing.T) {
			req := &wallet.MoneyTransactionRequest{Amount: 10}
			w := &wallet.Wallet{
				ID:      uuid.NewString(),
				Balance: wallet.Money{Amount: 10},
			}

			e := errors.New("")

			mockRepo.EXPECT().FindWalletByID(ctx, w.ID).Return(w, nil)
			mockRepo.EXPECT().FindTransactionsByWalletID(ctx, w.ID).Return(nil, e)

			w, err := service.DepositMoney(ctx, w.ID, req)
			assert.Nil(t, w)
			assert.Equal(t, errr.ThrowInternalServerError(e), err)
		})

		t.Run("should return error if repository.InsertTransactions returns error", func(t *testing.T) {
			req := &wallet.MoneyTransactionRequest{Amount: 10}
			w := &wallet.Wallet{
				ID:      uuid.NewString(),
				Balance: wallet.Money{Amount: 10},
			}

			e := errors.New("")

			mockRepo.EXPECT().FindWalletByID(ctx, w.ID).Return(w, nil)
			mockRepo.EXPECT().FindTransactionsByWalletID(ctx, w.ID).Return(nil, nil)
			mockRepo.EXPECT().InsertTransactions(ctx, gomock.Any()).Return(e)

			w, err := service.DepositMoney(ctx, w.ID, req)
			assert.Nil(t, w)
			assert.Equal(t, errr.ThrowInternalServerError(e), err)
		})
	})

	t.Run("WithdrawMoney", func(t *testing.T) {
		req := &wallet.MoneyTransactionRequest{Amount: 1}
		w := &wallet.Wallet{
			ID:      uuid.NewString(),
			Balance: wallet.Money{Amount: 10},
		}

		mockRepo.EXPECT().FindWalletByID(ctx, w.ID).Return(w, nil)
		mockRepo.EXPECT().FindTransactionsByWalletID(ctx, w.ID).Return(nil, nil)
		mockRepo.EXPECT().InsertTransactions(ctx, gomock.Any()).Return(nil)

		expected, err := service.WithdrawMoney(ctx, w.ID, req)
		assert.Nil(t, err)
		assert.Equal(t, expected, w)
	})

	t.Run("GetWallet", func(t *testing.T) {
		w := &wallet.Wallet{
			ID:      uuid.NewString(),
			Balance: wallet.Money{Amount: 10},
		}

		mockRepo.EXPECT().FindWalletByID(ctx, w.ID).Return(w, nil)
		mockRepo.EXPECT().FindTransactionsByWalletID(ctx, w.ID).Return(nil, nil)

		expected, err := service.GetWallet(ctx, w.ID)
		assert.Nil(t, err)
		assert.Equal(t, expected, w)
	})
}
