package wallet

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	walletsCollection      = "wallets"
	transactionsCollection = "transactions"
)

type (
	// Repository is an interface that do db operations
	Repository interface {
		// InsertWallet inserts wallet to collection
		InsertWallet(ctx context.Context, w *Wallet) error
		// FindWalletByID finds wallet by id
		FindWalletByID(ctx context.Context, id string) (*Wallet, error)
		// FindTransactionsByWalletID finds transaction by wallet id
		FindTransactionsByWalletID(ctx context.Context, walletID string) ([]Transaction, error)

		// InsertTransactions inserts transactions to collection
		InsertTransactions(ctx context.Context, transactions ...Transaction) error
	}

	// MongoRepository is a concrete implementation of Repository interface
	MongoRepository struct {
		wallets      *mongo.Collection
		transactions *mongo.Collection
	}
)

// NewMongoRepository creates instance of MongoRepository
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		wallets:      db.Collection(walletsCollection),
		transactions: db.Collection(transactionsCollection),
	}
}

// InsertWallet inserts wallet to collection
func (r *MongoRepository) InsertWallet(ctx context.Context, w *Wallet) error {
	_, err := r.wallets.InsertOne(ctx, w)
	if err != nil {
		return err
	}

	return nil
}

// InsertTransactions inserts transactions to collection
func (r *MongoRepository) InsertTransactions(ctx context.Context, transactions ...Transaction) error {
	documents := make([]interface{}, len(transactions))
	for i, t := range transactions {
		documents[i] = t
	}

	_, err := r.transactions.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	return nil
}

// FindWalletByID finds wallet by id
func (r *MongoRepository) FindWalletByID(ctx context.Context, id string) (*Wallet, error) {
	res := r.wallets.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, res.Err()
	}

	wallet := new(Wallet)
	if err := res.Decode(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// FindTransactionsByWalletID finds transaction by wallet id
func (r *MongoRepository) FindTransactionsByWalletID(ctx context.Context, walletID string) ([]Transaction, error) {
	opts := options.Find()
	opts.Sort = bson.M{"created_at": 1}

	cursor, err := r.transactions.Find(ctx, bson.M{
		"wallet_id": walletID,
	}, opts)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
