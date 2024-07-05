package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "datasource/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewDataSourceServiceClient(conn)

	// Use a timeout for our gRPC calls
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the data source
	connResp, err := c.Connect(ctx, &pb.ConnectRequest{ConnectorName: "mongo"})
	if err != nil {
		fmt.Printf("Could not connect: %v", err)
	}
	log.Printf("Connect Response: %t", connResp.GetSuccess())

	// Insert test data
	insertQuery := map[string]interface{}{
		"type":       "INSERT",
		"collection": "users",
		"data": map[string]interface{}{
			"name":  "John Doe",
			"age":   30,
			"email": "john@example.com",
			"address": map[string]interface{}{
				"city": "New York",
				"zip":  "10001",
			},
		},
	}
	insertQueryJSON, _ := json.Marshal(insertQuery)
	_, err = c.ExecuteCommand(ctx, &pb.CommandRequest{
		ConnectorName: "mongo",
		Command:       string(insertQueryJSON),
	})
	if err != nil {
		fmt.Printf("Could not insert data: %v", err)
	}
	fmt.Println("Inserted test data successfully")

	// Query data
	selectQuery := map[string]interface{}{
		"type":       "SELECT",
		"collection": "users",
		"conditions": map[string]interface{}{"name": "John Doe"},
	}
	selectQueryJSON, _ := json.Marshal(selectQuery)
	queryResp, err := c.ExecuteQuery(ctx, &pb.QueryRequest{
		ConnectorName: "mongo",
		Query:         string(selectQueryJSON),
	})
	if err != nil {
		fmt.Printf("Could not query data: %v", err)
	}

	fmt.Println("Query results:")
	for _, row := range queryResp.GetRows() {
		var result map[string]interface{}
		json.Unmarshal(row, &result)
		fmt.Printf("%+v\n", result)
	}

	// Update data
	updateQuery := map[string]interface{}{
		"type":       "UPDATE",
		"collection": "users",
		"conditions": map[string]interface{}{"name": "John Doe"},
		"data": map[string]interface{}{
			"$set": map[string]interface{}{"age": 31},
		},
	}
	updateQueryJSON, _ := json.Marshal(updateQuery)
	_, err = c.ExecuteCommand(ctx, &pb.CommandRequest{
		ConnectorName: "mongo",
		Command:       string(updateQueryJSON),
	})
	if err != nil {
		fmt.Printf("Could not update data: %v", err)
	}
	fmt.Println("Updated data successfully")

	// Query updated data
	queryResp, err = c.ExecuteQuery(ctx, &pb.QueryRequest{
		ConnectorName: "example_mongo",
		Query:         string(selectQueryJSON),
	})
	if err != nil {
		fmt.Printf("Could not query updated data: %v", err)
	}
	fmt.Println("Updated query results:")
	for _, row := range queryResp.GetRows() {
		var result map[string]interface{}
		json.Unmarshal(row, &result)
		fmt.Printf("%+v\n", result)
	}

	// Delete data
	deleteQuery := map[string]interface{}{
		"type":       "DELETE",
		"collection": "users",
		"conditions": map[string]interface{}{"name": "John Doe"},
	}
	deleteQueryJSON, _ := json.Marshal(deleteQuery)
	_, err = c.ExecuteCommand(ctx, &pb.CommandRequest{
		ConnectorName: "mongo",
		Command:       string(deleteQueryJSON),
	})
	if err != nil {
		fmt.Printf("Could not delete data: %v", err)
	}
	fmt.Println("Deleted data successfully")

	// Confirm deletion
	queryResp, err = c.ExecuteQuery(ctx, &pb.QueryRequest{
		ConnectorName: "mongo",
		Query:         string(selectQueryJSON),
	})
	if err != nil {
		fmt.Printf("Could not query data after deletion: %v", err)
	}
	fmt.Printf("Number of results after deletion: %d\n", len(queryResp.GetRows()))

	// Disconnect from the data source
	disconnResp, err := c.Disconnect(ctx, &pb.DisconnectRequest{ConnectorName: "example_mongo"})
	if err != nil {
		fmt.Printf("Could not disconnect: %v", err)
	}
	log.Printf("Disconnect Response: %t", disconnResp.GetSuccess())
}
