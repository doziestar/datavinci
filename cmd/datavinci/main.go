package main

import (
	"context"
	"datasource/connectors"
	"datasource/managers/query"
	"datasource/managers/transform"
	"fmt"
	"log"
)

func main() {
    config := &connectors.Config{
        
    }
    connector, err := connectors.ConnectorFactory(config)
    if err != nil {
        log.Fatalf("Failed to create MongoDB connector: %v", err)
    }


    // Connect to MongoDB
    ctx := context.Background()
    err = connector.Connect(ctx)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer connector.Close(ctx)

	log.Println("Connected to MongoDB successfully")

    // Create query executor
    executor := query.NewQueryExecutor(connector)

    // Create transformer
    transformer := transform.NewTransformer()

    // Insert test data
    insertQuery := query.Query{
        Type:       query.Insert,
        Collection: "users",
        Data: map[string]interface{}{
            "name": "John Doe",
            "age":  30,
            "email": "john@example.com",
            "address": map[string]interface{}{
                "city": "New York",
                "zip":  "10001",
            },
        },
    }
    _, err = executor.Execute(ctx, insertQuery)
    if err != nil {
        log.Fatalf("Failed to insert data: %v", err)
    }
    fmt.Println("Inserted test data successfully")

    // Query data
    selectQuery := query.Query{
        Type:       query.Select,
        Collection: "users",
        Conditions: map[string]interface{}{"name": "John Doe"},
    }
    results, err := executor.Execute(ctx, selectQuery)
    if err != nil {
        log.Fatalf("Failed to query data: %v", err)
    }

    fmt.Println("Query results:")
    for _, result := range results {
        fmt.Printf("%+v\n", result)
    }

    // Transform data
    if len(results) > 0 {
        flattenedData, err := transformer.TransformData(results[0], "map")
        if err != nil {
            log.Fatalf("Failed to transform data: %v", err)
        }
        flattenedData = transformer.FlattenMap(flattenedData.(map[string]interface{}), "")
        fmt.Println("Flattened data:")
        for k, v := range flattenedData.(map[string]interface{}) {
            fmt.Printf("%s: %v\n", k, v)
        }

        // Extract a specific field
        city, err := transform.ExtractField(results[0], "address.city")
        if err != nil {
            log.Fatalf("Failed to extract field: %v", err)
        }
        fmt.Printf("Extracted city: %v\n", city)

        // Convert age to string
        ageStr, err := transformer.ConvertType(results[0]["age"], "string")
        if err != nil {
            log.Fatalf("Failed to convert age to string: %v", err)
        }
        fmt.Printf("Age as string: %s\n", ageStr)
    }

    // Update data
    updateQuery := query.Query{
        Type:       query.Update,
        Collection: "users",
        Conditions: map[string]interface{}{"name": "John Doe"},
        Data: map[string]interface{}{
            "$set": map[string]interface{}{"age": 31},
        },
    }
    _, err = executor.Execute(ctx, updateQuery)
    if err != nil {
        log.Fatalf("Failed to update data: %v", err)
    }
    fmt.Println("Updated data successfully")

    // Query updated data
    results, err = executor.Execute(ctx, selectQuery)
    if err != nil {
        log.Fatalf("Failed to query updated data: %v", err)
    }
    fmt.Println("Updated query results:")
    for _, result := range results {
        fmt.Printf("%+v\n", result)
    }

    // Delete data
    deleteQuery := query.Query{
        Type:       query.Delete,
        Collection: "users",
        Conditions: map[string]interface{}{"name": "John Doe"},
    }
    _, err = executor.Execute(ctx, deleteQuery)
    if err != nil {
        log.Fatalf("Failed to delete data: %v", err)
    }
    fmt.Println("Deleted data successfully")

    // Confirm deletion
    results, err = executor.Execute(ctx, selectQuery)
    if err != nil {
        log.Fatalf("Failed to query after deletion: %v", err)
    }
    fmt.Printf("Number of results after deletion: %d\n", len(results))
}