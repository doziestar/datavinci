package connectors

import (
	"context"
	"fmt"

	"pkg/common/errors"
	"pkg/common/retry"

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
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		c.config.Username, c.config.Password, c.config.Host, c.config.Port)

	clientOptions := options.Client().ApplyURI(uri)

	var client *mongo.Client
	err := retry.Retry(ctx, func() error {
		var err error
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			return errors.NewError(errors.ErrorTypeDatabaseConnection, "failed to connect to MongoDB", err)
		}
		return client.Ping(ctx, nil)
	}, retry.DefaultConfig())

	if err != nil {
		return err
	}

	c.client = client
	return nil
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
	if c.client == nil {
		return nil, errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	// Parse the query string into a BSON document
	var filter bson.D
	err := bson.UnmarshalExtJSON([]byte(query), true, &filter)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to parse query", err)
	}

	collection := c.client.Database(c.config.Database).Collection(args[0].(string))
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to execute query", err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	for cursor.Next(ctx) {
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
}

// Execute executes a command and returns the number of affected documents.
func (c *MongoConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	if c.client == nil {
		return 0, errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	var doc bson.D
	err := bson.UnmarshalExtJSON([]byte(command), true, &doc)
	if err != nil {
		return 0, errors.NewError(errors.ErrorTypeExecution, "failed to parse command", err)
	}

	collection := c.client.Database(c.config.Database).Collection(args[0].(string))
	result, err := collection.UpdateMany(ctx, doc[0].Value, doc[1].Value)
	if err != nil {
		return 0, errors.NewError(errors.ErrorTypeExecution, "failed to execute command", err)
	}

	return result.ModifiedCount, nil
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