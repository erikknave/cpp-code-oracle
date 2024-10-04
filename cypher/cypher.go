package cypher

import (
	"context"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Define a struct to hold the response data
type QueryResult struct {
	Key   int         `json:"key"`
	Value interface{} `json:"value"`
}

func InjectCypher(query string) interface{} {
	uri := os.Getenv("NEO4J_URL")
	username := os.Getenv("NEO4J_USER")
	password := os.Getenv("NEO4J_PASSWORD")

	// Create a Neo4j driver
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}
	defer driver.Close(context.Background())

	// Open a new session
	session := driver.NewSession(context.TODO(), neo4j.SessionConfig{DatabaseName: os.Getenv("NEO4J_DB")})
	if session == nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close(context.Background())

	// Execute the query
	result, err := session.Run(context.TODO(), query, nil)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	// var results []QueryResult
	var results []interface{}

	// Process the result
	for result.Next(context.TODO()) {
		record := result.Record()
		results = append(results, record.Values...)
		// for key, value := range record.Values {
		// 	queryResult := QueryResult{
		// 		Key:   key,
		// 		Value: value,
		// 	}
		// 	results = append(results, queryResult)
		// }
	}

	// Check for errors at the end of the query
	if err = result.Err(); err != nil {
		log.Fatalf("Error in result processing: %v", err)
	}

	// Convert results to pretty-printed JSON
	// jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal results to JSON: %v", err)
	}
	// golog.Golog("Cypher response", results)

	// Print the pretty-printed JSON data
	// fmt.Println(string(jsonData))
	return results
}
