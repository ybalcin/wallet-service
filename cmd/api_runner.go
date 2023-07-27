package cmd

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/ybalcin/wallet-service/internal/config"
	"github.com/ybalcin/wallet-service/internal/wallet"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RunApi(ctx context.Context) error {
	cfg, err := config.Read()
	if err != nil {
		return err
	}

	opts := options.Client().ApplyURI(cfg.MongoSettings.URI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	db := client.Database(cfg.MongoSettings.Database)

	walletRepo := wallet.NewMongoRepository(db)
	walletService := wallet.NewService(walletRepo)
	walletApi := wallet.NewApi(walletService)
	root := NewApiRoot(walletApi)

	go func() {
		if err = root.Listen(cfg.Port); err != nil {
			log.Fatalf("server listen error: %v\n", err)
		}
	}()

	// graceful
	<-ctx.Done()
	if err = root.app.Shutdown(); err != nil {
		log.Fatalf("server shutdown failed: %+v", err)
	}
	fmt.Println("server exited properly")

	return nil
}
