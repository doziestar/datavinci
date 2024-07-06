package query

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"datasource/connectors"
	"pkg/common/errors"
)

// QueryType represents the type of query operation
type QueryType string

const (
	Select QueryType = "SELECT"
	Insert QueryType = "INSERT"
	Update QueryType = "UPDATE"
	Delete QueryType = "DELETE"
)

// Query represents a unified query structure
type Query struct {
	Type       QueryType              `json:"type"`
	Collection string                 `json:"collection"`
	Fields     []string               `json:"fields,omitempty"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Limit      int                    `json:"limit,omitempty"`
	Offset     int                    `json:"offset,omitempty"`
}

// QueryExecutor handles query execution across different connector types
type QueryExecutor struct {
	connector connectors.Connector
}

// NewQueryExecutor creates a new QueryExecutor
func NewQueryExecutor(connector connectors.Connector) *QueryExecutor {
	return &QueryExecutor{connector: connector}
}

// Execute executes the given query on the appropriate connector
func (qe *QueryExecutor) Execute(ctx context.Context, query Query) ([]map[string]interface{}, error) {
	switch c := qe.connector.(type) {
	case *connectors.SQLConnector:
		return qe.executeSQL(ctx, c, query)
	case *connectors.MongoConnector:
		return qe.executeMongo(ctx, c, query)
	case *connectors.RedisConnector:
		return qe.executeRedis(ctx, c, query)
	case *connectors.FileConnector:
		return qe.executeFile(ctx, c, query)
	default:
		return nil, errors.NewError(errors.ErrorTypeUnsupported, "unsupported connector type", nil)
	}
}

func (qe *QueryExecutor) executeSQL(ctx context.Context, connector *connectors.SQLConnector, query Query) ([]map[string]interface{}, error) {
	sqlQuery, args := buildSQLQuery(query)
	if query.Type == Select {
		return connector.Query(ctx, sqlQuery, args...)
	}
	affected, err := connector.Execute(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	return []map[string]interface{}{{"affected_rows": affected}}, nil
}

func (qe *QueryExecutor) executeMongo(ctx context.Context, connector *connectors.MongoConnector, query Query) ([]map[string]interface{}, error) {
	switch query.Type {
	case Select:
		filterJSON, err := json.Marshal(query.Conditions)
		if err != nil {
			return nil, errors.NewError(errors.ErrorTypeQuery, "failed to marshal query conditions", err)
		}
		return connector.Query(ctx, string(filterJSON), query.Collection)
	case Insert:
		affected, err := connector.Execute(ctx, "insert", query.Collection, query.Data)
		if err != nil {
			return nil, err
		}
		return []map[string]interface{}{{"affected_documents": affected}}, nil
	case Update:
		affected, err := connector.Execute(ctx, "update", query.Collection, query.Conditions, query.Data)
		if err != nil {
			return nil, err
		}
		return []map[string]interface{}{{"affected_documents": affected}}, nil
	case Delete:
		affected, err := connector.Execute(ctx, "delete", query.Collection, query.Conditions)
		if err != nil {
			return nil, err
		}
		return []map[string]interface{}{{"affected_documents": affected}}, nil
	default:
		return nil, errors.NewError(errors.ErrorTypeUnsupported, "unsupported query type for MongoDB", nil)
	}
}

func (qe *QueryExecutor) executeRedis(ctx context.Context, connector *connectors.RedisConnector, query Query) ([]map[string]interface{}, error) {
	switch query.Type {
	case Select:
		return connector.Query(ctx, query.Collection)
	case Insert, Update:
		value, err := json.Marshal(query.Data)
		if err != nil {
			return nil, errors.NewError(errors.ErrorTypeQuery, "failed to marshal Redis data", err)
		}
		affected, err := connector.Execute(ctx, "SET", query.Collection, string(value))
		if err != nil {
			return nil, err
		}
		return []map[string]interface{}{{"affected_keys": affected}}, nil
	case Delete:
		affected, err := connector.Execute(ctx, "DEL", query.Collection)
		if err != nil {
			return nil, err
		}
		return []map[string]interface{}{{"affected_keys": affected}}, nil
	default:
		return nil, errors.NewError(errors.ErrorTypeUnsupported, "unsupported query type for Redis", nil)
	}
}

func (qe *QueryExecutor) executeFile(ctx context.Context, connector *connectors.FileConnector, query Query) ([]map[string]interface{}, error) {
	if query.Type != Select {
		return nil, errors.NewError(errors.ErrorTypeUnsupported, "only SELECT queries are supported for file connector", nil)
	}
	return connector.Query(ctx, query.Collection)
}

func buildSQLQuery(query Query) (string, []interface{}) {
	var sqlQuery strings.Builder
	var args []interface{}

	switch query.Type {
	case Select:
		sqlQuery.WriteString("SELECT ")
		if len(query.Fields) > 0 {
			sqlQuery.WriteString(strings.Join(query.Fields, ", "))
		} else {
			sqlQuery.WriteString("*")
		}
		sqlQuery.WriteString(" FROM ")
		sqlQuery.WriteString(query.Collection)
	case Insert:
		sqlQuery.WriteString("INSERT INTO ")
		sqlQuery.WriteString(query.Collection)
		sqlQuery.WriteString(" (")
		var columns []string
		var values []string
		for k, v := range query.Data {
			columns = append(columns, k)
			values = append(values, "?")
			args = append(args, v)
		}
		sqlQuery.WriteString(strings.Join(columns, ", "))
		sqlQuery.WriteString(") VALUES (")
		sqlQuery.WriteString(strings.Join(values, ", "))
		sqlQuery.WriteString(")")
	case Update:
		sqlQuery.WriteString("UPDATE ")
		sqlQuery.WriteString(query.Collection)
		sqlQuery.WriteString(" SET ")
		var sets []string
		for k, v := range query.Data {
			sets = append(sets, k+" = ?")
			args = append(args, v)
		}
		sqlQuery.WriteString(strings.Join(sets, ", "))
	case Delete:
		sqlQuery.WriteString("DELETE FROM ")
		sqlQuery.WriteString(query.Collection)
	}

	if len(query.Conditions) > 0 && (query.Type == Select || query.Type == Update || query.Type == Delete) {
		sqlQuery.WriteString(" WHERE ")
		var conditions []string
		for k, v := range query.Conditions {
			conditions = append(conditions, k+" = ?")
			args = append(args, v)
		}
		sqlQuery.WriteString(strings.Join(conditions, " AND "))
	}

	if query.Type == Select {
		if query.Limit > 0 {
			sqlQuery.WriteString(fmt.Sprintf(" LIMIT %d", query.Limit))
		}
		if query.Offset > 0 {
			sqlQuery.WriteString(fmt.Sprintf(" OFFSET %d", query.Offset))
		}
	}

	return sqlQuery.String(), args
}
