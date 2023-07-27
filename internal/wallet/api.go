package wallet

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ybalcin/wallet-service/pkg/errr"
	"github.com/ybalcin/wallet-service/pkg/response"
)

type Api struct {
	service Service
}

func NewApi(service Service) *Api {
	return &Api{service: service}
}

func (a *Api) AddRoutesTo(r fiber.Router) {
	wallets := r.Group("wallets")

	wallets.Post("/", a.CreateWallet)
	wallets.Put("/:id/deposit", a.DepositMoney)
	wallets.Put("/:id/withdraw", a.WithdrawMoney)
	wallets.Get("/:id", a.GetWallet)
}

func (a *Api) CreateWallet(c *fiber.Ctx) error {
	createWalletDto := new(CreateWalletRequest)
	if err := c.BodyParser(createWalletDto); err != nil {
		return response.New(c).Error(errr.ThrowBadRequestError(err)).JSON()
	}

	id, err := a.service.CreateWallet(c.UserContext(), createWalletDto)
	if err != nil {
		return response.New(c).Error(err).JSON()
	}

	return response.New(c).Data(id).JSON()
}

func (a *Api) DepositMoney(c *fiber.Ctx) error {
	id := c.Params("id")
	transaction := new(MoneyTransactionRequest)
	if err := c.BodyParser(transaction); err != nil {
		return response.New(c).Error(errr.ThrowBadRequestError(err)).JSON()
	}

	wallet, err := a.service.DepositMoney(c.UserContext(), id, transaction)
	if err != nil {
		return response.New(c).Error(err).JSON()
	}

	return response.New(c).Data(wallet).JSON()
}

func (a *Api) WithdrawMoney(c *fiber.Ctx) error {
	id := c.Params("id")
	transaction := new(MoneyTransactionRequest)
	if err := c.BodyParser(transaction); err != nil {
		return response.New(c).Error(errr.ThrowBadRequestError(err)).JSON()
	}

	wallet, err := a.service.WithdrawMoney(c.UserContext(), id, transaction)
	if err != nil {
		return response.New(c).Error(err).JSON()
	}

	return response.New(c).Data(wallet).JSON()
}

func (a *Api) GetWallet(c *fiber.Ctx) error {
	id := c.Params("id")
	wallet, err := a.service.GetWallet(c.UserContext(), id)
	if err != nil {
		return response.New(c).Error(err).JSON()
	}

	return response.New(c).Data(wallet).JSON()
}
