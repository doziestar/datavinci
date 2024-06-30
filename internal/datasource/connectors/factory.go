package connectors

import (
	"fmt"
)

// ConnectorFactory creates and returns a Connector based on the provided Config.
func ConnectorFactory(config *Config) (Connector, error) {
	switch config.Type {
	case "sql":
		return NewSQLConnector(config), nil
	case "redis":
		return NewRedisConnector(config), nil
	case "mongo":
		return NewMongoConnector(config), nil
	case "api":
		// return NewAPIConnector(config), nil
	default:
		return nil, fmt.Errorf("unsupported connector type: %s", config.Type)
	}
	return nil, fmt.Errorf("connector type not implemented: %s", config.Type)
}