package manager

import (
	"context"
	"fmt"
	"sync"

	"datasource/connectors"
	"pkg/common/errors"
)

// ConnectorManager manages multiple data source connectors.
type ConnectorManager struct {
	connectors map[string]connectors.Connector
	mu         sync.RWMutex
}

// NewConnectorManager creates a new ConnectorManager.
func NewConnectorManager() *ConnectorManager {
	return &ConnectorManager{
		connectors: make(map[string]connectors.Connector),
	}
}

// AddConnector adds a new connector with the given name.
func (m *ConnectorManager) AddConnector(name string, config *connectors.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.connectors[name]; exists {
		return errors.NewError(errors.ErrorTypeConfiguration, "connector with this name already exists", nil)
	}

	connector, err := connectors.ConnectorFactory(config)
	if err != nil {
		return err
	}

	m.connectors[name] = connector
	return nil
}

// GetConnector retrieves a connector by name.
func (m *ConnectorManager) GetConnector(name string) (connectors.Connector, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connector, exists := m.connectors[name]
	if !exists {
		return nil, errors.NewError(errors.ErrorTypeNotFound, fmt.Sprintf("connector '%s' not found", name), nil)
	}

	return connector, nil
}

// RemoveConnector removes a connector by name.
func (m *ConnectorManager) RemoveConnector(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.connectors[name]; !exists {
		return errors.NewError(errors.ErrorTypeNotFound, fmt.Sprintf("connector '%s' not found", name), nil)
	}

	delete(m.connectors, name)
	return nil
}

// ConnectAll connects all managed connectors.
func (m *ConnectorManager) ConnectAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, connector := range m.connectors {
		if err := connector.Connect(ctx); err != nil {
			return errors.NewError(errors.ErrorTypeConnection, fmt.Sprintf("failed to connect '%s'", name), err)
		}
	}

	return nil
}

// CloseAll closes all managed connectors.
func (m *ConnectorManager) CloseAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, connector := range m.connectors {
		if err := connector.Close(ctx); err != nil {
			return errors.NewError(errors.ErrorTypeConnection, fmt.Sprintf("failed to close '%s'", name), err)
		}
	}

	return nil
}
