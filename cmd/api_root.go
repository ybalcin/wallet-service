package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ybalcin/wallet-service/internal/wallet"
)

type ApiRoot struct {
	app *fiber.App

	walletApi *wallet.Api
}

func NewApiRoot(walletApi *wallet.Api) *ApiRoot {
	app := fiber.New()

	root := &ApiRoot{
		app:       app,
		walletApi: walletApi,
	}
	root.RegisterRoutes(walletApi)

	return root
}

func (r *ApiRoot) RegisterRoutes(walletApi *wallet.Api) {
	group := r.app.Group("api")
	walletApi.AddRoutesTo(group)
}

func (r *ApiRoot) Listen(port string) error {
	return r.app.Listen(":" + port)
}
