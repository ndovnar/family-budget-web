package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ndovnar/family-budget-api/internal/store"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ store.Store = (*Mongo)(nil)

type Mongo struct {
	database *mongo.Database
}

func New(ctx context.Context, config Config, applicationName string) (*Mongo, error) {
	connectionDebugInfo := databaseConnectionInfoString(config)

	clientOptions := options.Client().
		SetHosts(config.Hosts).
		SetAuth(options.Credential{
			AuthSource: config.Database,
			Username:   config.Username,
			Password:   config.Password,
		}).
		SetConnectTimeout(5 * time.Second).
		SetAppName(applicationName)

	if config.UseTLS {
		clientOptions.SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		err := fmt.Errorf("could not create MongoDB client [%v]: %w", connectionDebugInfo, err)
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	err = client.Ping(ctxTimeout, nil)
	cancel()
	if err != nil {
		return nil, fmt.Errorf("connection to MongoDB [%v] worked, but ping failed: %w", connectionDebugInfo, err)
	}

	database := client.Database(config.Database)
	if err := initializeDatabaseIfNeeded(ctx, database); err != nil {
		return nil, err
	}

	return &Mongo{
		database: database,
	}, nil
}

func databaseConnectionInfoString(config Config) string {
	return fmt.Sprintf(
		"mongodb://%s@%s/%s",
		config.Username,
		strings.Join(config.Hosts, ","),
		config.Database,
	)
}

func initializeDatabaseIfNeeded(ctx context.Context, database *mongo.Database) error {
	usersCollection := database.Collection(CollectionUsers)
	_, err := usersCollection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			// an index for key "_id" is created by default
			{Keys: bson.M{"id": 1}},
			{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
		},
	)
	return err
}

func (m *Mongo) Status(ctx context.Context) error {
	if client := m.database.Client(); client != nil {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := client.Ping(ctxTimeout, nil)
		cancel()
		return err
	}
	return errors.New("mongo client not initialized")
}
