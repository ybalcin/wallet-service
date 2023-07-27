package cmd

import (
	"context"
	"github.com/ybalcin/wallet-service/internal/wallet"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RunApi(ctx context.Context) error {
	opts := options.Client().ApplyURI("mongodb+srv://wallet-service-user:EmGcfRCHL9tAoyYZ@cluster0.l1pmb.mongodb.net/?retryWrites=true&w=majority")

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	db := client.Database("wallet-service-db")

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	walletRepo := wallet.NewMongoRepository(db)
	walletService := wallet.NewService(walletRepo)
	walletApi := wallet.NewApi(walletService)

	root := NewApiRoot(walletApi)

	return root.Listen()
}
