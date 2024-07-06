package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"pkg/common/errors"
	"pkg/common/retry"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConnector implements the Connector interface for MongoDB.
type MongoConnector struct {
	client *mongo.Client
	config *Config
}

// NewMongoConnector creates a new MongoConnector with the given configuration.
func NewMongoConnector(config *Config) *MongoConnector {
	return &MongoConnector{config: config}
}

// Connect establishes a connection to the MongoDB database.
func (c *MongoConnector) Connect(ctx context.Context) error {
	uri := c.buildConnectionString()

	clientOptions := options.Client().ApplyURI(uri)

	log.Printf("Connecting to MongoDB: %s", uri)

	var client *mongo.Client
	err := retry.Retry(ctx, func() error {
		var err error
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("Failed to connect to MongoDB: %v", err)
			return errors.NewError(errors.ErrorTypeConnection, "failed to connect to MongoDB", err)
		}
		return client.Ping(ctx, nil)
	}, retry.DefaultConfig())

	if err != nil {
		return err
	}

	c.client = client
	return nil
}

func (c *MongoConnector) buildConnectionString() string {
	query := url.Values{}
	for k, v := range c.config.Options {
		query.Add(k, fmt.Sprintf("%v", v))
	}

	var baseURL string
	if c.config.Port > 0 {
		// If port is provided, use standard MongoDB protocol
		baseURL = fmt.Sprintf("mongodb://%s:%s@%s:%d",
			url.QueryEscape(c.config.Username),
			url.QueryEscape(c.config.Password),
			c.config.Host,
			c.config.Port)
	} else {
		// If no port, assume it's MongoDB Atlas and use srv protocol
		baseURL = fmt.Sprintf("mongodb+srv://%s:%s@%s",
			url.QueryEscape(c.config.Username),
			url.QueryEscape(c.config.Password),
			c.config.Host)
	}

	// Append database and query parameters
	return fmt.Sprintf("%s/%s?%s", baseURL, c.config.Database, query.Encode())
}

// Close closes the connection to the MongoDB database.
func (c *MongoConnector) Close(ctx context.Context) error {
	if c.client == nil {
		return errors.NewError(errors.ErrorTypeDatabaseConnection, "connection already closed", nil)
	}
	return c.client.Disconnect(ctx)
}

// Query executes a query and returns the results as a slice of maps.
func (c *MongoConnector) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, errors.NewError(errors.ErrorTypeQuery, "missing collection name", nil)
	}
	collection, ok := args[0].(string)
	if !ok {
		return nil, errors.NewError(errors.ErrorTypeQuery, "invalid collection name", nil)
	}

	coll := c.client.Database(c.config.Database).Collection(collection)

	var filter bson.M
	err := json.Unmarshal([]byte(query), &filter)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to parse query", err)
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to execute query", err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to decode query results", err)
	}

	return results, nil
}

// Execute executes a command and returns the number of affected documents.
func (c *MongoConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	if len(args) == 0 {
		return 0, errors.NewError(errors.ErrorTypeExecution, "missing collection name", nil)
	}
	collectionName, ok := args[0].(string)
	if !ok {
		return 0, errors.NewError(errors.ErrorTypeExecution, "invalid collection name", nil)
	}

	collection := c.client.Database(c.config.Database).Collection(collectionName)

	var result int64
	var err error

	switch command {
	case "insert":
		if len(args) < 2 {
			return 0, errors.NewError(errors.ErrorTypeExecution, "missing document to insert", nil)
		}
		doc, ok := args[1].(map[string]interface{})
		if !ok {
			return 0, errors.NewError(errors.ErrorTypeExecution, "invalid document format", nil)
		}
		_, err = collection.InsertOne(ctx, doc)
		if err == nil {
			result = 1
		}
	case "update":
		if len(args) < 3 {
			return 0, errors.NewError(errors.ErrorTypeExecution, "missing update parameters", nil)
		}
		filter, ok := args[1].(map[string]interface{})
		if !ok {
			return 0, errors.NewError(errors.ErrorTypeExecution, "invalid filter format", nil)
		}
		update, ok := args[2].(map[string]interface{})
		if !ok {
			return 0, errors.NewError(errors.ErrorTypeExecution, "invalid update format", nil)
		}
		updateResult, err := collection.UpdateMany(ctx, filter, bson.M{"$set": update})
		if err == nil {
			result = updateResult.ModifiedCount
		}
	case "delete":
		if len(args) < 2 {
			return 0, errors.NewError(errors.ErrorTypeExecution, "missing delete parameters", nil)
		}
		filter, ok := args[1].(map[string]interface{})
		if !ok {
			return 0, errors.NewError(errors.ErrorTypeExecution, "invalid filter format", nil)
		}
		deleteResult, err := collection.DeleteMany(ctx, filter)
		if err == nil {
			result = deleteResult.DeletedCount
		}
	default:
		return 0, errors.NewError(errors.ErrorTypeExecution, fmt.Sprintf("unsupported command: %s", command), nil)
	}

	if err != nil {
		return 0, errors.NewError(errors.ErrorTypeExecution, fmt.Sprintf("failed to execute command: %s", command), err)
	}

	return result, nil
}

// Ping checks if the database connection is still alive.
func (c *MongoConnector) Ping(ctx context.Context) error {
	if c.client == nil {
		return errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}
	return c.client.Ping(ctx, nil)
}

// Transaction starts a new database transaction and returns a TransactionConnector.
func (c *MongoConnector) Transaction(ctx context.Context) (TransactionConnector, error) {
	if c.client == nil {
		return nil, errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	session, err := c.client.StartSession()
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeTransaction, "failed to start session", err)
	}

	if err := session.StartTransaction(); err != nil {
		return nil, errors.NewError(errors.ErrorTypeTransaction, "failed to start transaction", err)
	}

	return &MongoTransactionConnector{session: session, client: c.client, config: c.config}, nil
}

// MongoTransactionConnector implements the TransactionConnector interface for MongoDB.
type MongoTransactionConnector struct {
	session mongo.Session
	client  *mongo.Client
	config  *Config
}

// Query executes a query within the transaction and returns the results.
func (c *MongoTransactionConnector) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	_, err := c.session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Parse the query string into a BSON document
		var filter bson.D
		err := bson.UnmarshalExtJSON([]byte(query), true, &filter)
		if err != nil {
			return nil, errors.NewError(errors.ErrorTypeQuery, "failed to parse query", err)
		}

		// Ensure we have at least one argument for the collection name
		if len(args) == 0 {
			return nil, errors.NewError(errors.ErrorTypeQuery, "missing collection name", nil)
		}
		collectionName, ok := args[0].(string)
		if !ok {
			return nil, errors.NewError(errors.ErrorTypeQuery, "invalid collection name", nil)
		}

		collection := c.client.Database(c.config.Database).Collection(collectionName)
		cursor, err := collection.Find(sessCtx, filter)
		if err != nil {
			return nil, errors.NewError(errors.ErrorTypeQuery, "failed to execute query", err)
		}
		defer cursor.Close(sessCtx)

		for cursor.Next(sessCtx) {
			var result map[string]interface{}
			if err := cursor.Decode(&result); err != nil {
				return nil, errors.NewError(errors.ErrorTypeQuery, "failed to decode result", err)
			}
			results = append(results, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, errors.NewError(errors.ErrorTypeQuery, "error during cursor iteration", err)
		}

		return results, nil
	})

	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeTransaction, "transaction failed", err)
	}

	return results, nil
}

// Execute executes a command within the transaction and returns the number of affected documents.
func (c *MongoTransactionConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	var modifiedCount int64

	_, err := c.session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var doc bson.D
		err := bson.UnmarshalExtJSON([]byte(command), true, &doc)
		if err != nil {
			return nil, errors.NewError(errors.ErrorTypeExecution, "failed to parse command", err)
		}

		// Ensure we have at least one argument for the collection name
		if len(args) == 0 {
			return nil, errors.NewError(errors.ErrorTypeExecution, "missing collection name", nil)
		}
		collectionName, ok := args[0].(string)
		if !ok {
			return nil, errors.NewError(errors.ErrorTypeExecution, "invalid collection name", nil)
		}

		collection := c.client.Database(c.config.Database).Collection(collectionName)

		// Assuming the first element is the filter and the second is the update
		if len(doc) < 2 {
			return nil, errors.NewError(errors.ErrorTypeExecution, "invalid command structure", nil)
		}

		result, err := collection.UpdateMany(sessCtx, doc[0].Value, doc[1].Value)
		if err != nil {
			return nil, errors.NewError(errors.ErrorTypeExecution, "failed to execute command", err)
		}

		modifiedCount = result.ModifiedCount
		return modifiedCount, nil
	})

	if err != nil {
		return 0, errors.NewError(errors.ErrorTypeTransaction, "transaction failed", err)
	}

	return modifiedCount, nil
}

// Commit commits the transaction.
func (c *MongoTransactionConnector) Commit(ctx context.Context) error {
	return c.session.CommitTransaction(ctx)
}

// Rollback rolls back the transaction.
func (c *MongoTransactionConnector) Rollback(ctx context.Context) error {
	return c.session.AbortTransaction(ctx)
}
