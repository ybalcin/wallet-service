package wallet

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/ybalcin/wallet-service/internal/wallet"
	"testing"
)

func TestNewWallet(t *testing.T) {
	t.Run("should return a new Wallet instance", func(t *testing.T) {
		w, err := wallet.New("user")
		assert.Nil(t, err)
		assert.NotNil(t, w)
	})

	t.Run("should return ErrInvalidUsername if username is empty", func(t *testing.T) {
		usernames := []string{" ", ""}
		for _, u := range usernames {
			w, err := wallet.New(u)
			assert.Nil(t, w)
			assert.Equal(t, errors.New(wallet.ErrInvalidUsername), err)
		}
	})
}

func TestWallet_WithdrawMoney(t *testing.T) {
	t.Run("should apply withdraw and append transaction to changes", func(t *testing.T) {
		money, _ := wallet.NewMoney(10)
		w, _ := wallet.New("user")
		w.Balance = wallet.Money{Amount: 20}
		err := w.WithdrawMoney(*money)
		assert.Nil(t, err)
		assert.Equal(t, float32(10), w.Balance.Amount)
		assert.Equal(t, 1, len(w.Changes))
	})

	t.Run("should return error", func(t *testing.T) {
		cases := []struct {
			name     string
			amount   float32
			expected error
		}{
			{
				"if money is 0",
				0,
				errors.New(wallet.ErrInvalidMoneyAmount),
			},
			{
				"if money is greater than wallet balance",
				10,
				errors.New(wallet.ErrInsufficientMoneyAmount),
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				money, _ := wallet.NewMoney(tt.amount)
				w, _ := wallet.New("user")
				w.Balance = wallet.Money{Amount: 1}
				actual := w.WithdrawMoney(*money)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})
}

func TestWallet_DepositMoney(t *testing.T) {
	t.Run("should apply deposit and append transaction to changes", func(t *testing.T) {
		money, _ := wallet.NewMoney(10)
		w, _ := wallet.New("username")
		err := w.DepositMoney(*money)
		assert.Nil(t, err)
		assert.Equal(t, float32(10), w.Balance.Amount)
		assert.Equal(t, 1, len(w.Changes))
	})

	t.Run("should return ErrInvalidMoneyAmount if money is 0", func(t *testing.T) {
		money, _ := wallet.NewMoney(0)
		w, _ := wallet.New("username")
		actual := w.DepositMoney(*money)
		assert.Equal(t, errors.New(wallet.ErrInvalidMoneyAmount), actual)
	})
}

func TestWallet_Mutate(t *testing.T) {
	cases := []struct {
		transaction wallet.Transaction
		balance     wallet.Money
		expected    float32
	}{
		{
			transaction: wallet.Transaction{
				Type:  wallet.DepositTransactionType,
				Money: wallet.Money{Amount: 10},
			},
			balance:  wallet.Money{Amount: 10},
			expected: 20,
		},
		{
			transaction: wallet.Transaction{
				Type:  wallet.DepositTransactionType,
				Money: wallet.Money{Amount: 10},
			},
			balance:  wallet.Money{Amount: 1.11},
			expected: 11.11,
		},
		{
			transaction: wallet.Transaction{
				Type:  wallet.WithdrawTransactionType,
				Money: wallet.Money{Amount: 9},
			},
			balance:  wallet.Money{Amount: 10},
			expected: 1,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			w, _ := wallet.New("user")
			w.Balance = tt.balance
			w.Mutate(tt.transaction)
			assert.Equal(t, tt.expected, w.Balance.Amount)
		})
	}
}

func TestNewMoney(t *testing.T) {
	t.Run("should return ErrInvalidMoneyAmount error", func(t *testing.T) {
		money, err := wallet.NewMoney(-1)
		assert.Nil(t, money)
		assert.Equal(t, errors.New(wallet.ErrInvalidMoneyAmount), err)
	})

	t.Run("should create instance of money", func(t *testing.T) {
		money, err := wallet.NewMoney(10)
		assert.Nil(t, err)
		assert.Equal(t, float32(10), money.Amount)
	})
}

func TestMoney_Add(t *testing.T) {
	money, _ := wallet.NewMoney(10)
	to := wallet.Money{Amount: 10}
	money.Add(&to)
	assert.Equal(t, float32(20), to.Amount)
}

func TestMoney_Sub(t *testing.T) {
	money, _ := wallet.NewMoney(1)
	from := wallet.Money{Amount: 10}
	money.Sub(&from)
	assert.Equal(t, float32(9), from.Amount)
}
