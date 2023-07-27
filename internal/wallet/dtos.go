package wallet

type (
	CreateWalletRequest struct {
		Username string `json:"username"`
	}

	MoneyTransactionRequest struct {
		Amount float32 `json:"amount"`
	}
)
