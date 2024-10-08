package search

import (
	"log"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

var Client *meilisearch.Client
var Index *meilisearch.Index

// Init initializes the Meilisearch client and ensures the index exists
func Init() {
	Client = meilisearch.NewClient(meilisearch.ClientConfig{
		Host: os.Getenv("MEILI_URL"), // Replace with your Meilisearch host
	})

	indexName := "searchable_documents"

	// Create the index
	// _, err := Client.CreateIndex(&meilisearch.IndexConfig{
	// 	Uid: indexName,
	// })
	// if err != nil {
	// 	log.Fatalf("Error creating index: %v", err)
	// }

	Index = Client.Index(indexName)

	// Update filterable attributes only if necessary
	_, err := Index.UpdateFilterableAttributes(&[]string{"doc_type", "dbid", "repository_id", "directory_id", "file_id", "name"})
	if err != nil {
		log.Fatalf("Error setting filterable attributes: %v", err)
	}
}
