package pgvector

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func PerformSearch(queryString string, numDocuments int) []types.SearchableDocument {
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	store, err := pgvector.New(
		ctx,
		pgvector.WithConnectionURL(os.Getenv("PGVECTOR_URL")),
		pgvector.WithEmbedder(e),
	)
	if err != nil {
		log.Fatal(err)
	}
	docs, err := store.SimilaritySearch(context.Background(), queryString, numDocuments)
	if err != nil {
		log.Fatal(err)
	}
	searchDocs := []types.SearchableDocument{}
	for _, doc := range docs {
		metadata := doc.Metadata
		doc_json := metadata["container"]
		var searchableDoc types.SearchableDocument

		// Convert map[string]interface{} to JSON
		docBytes, err := json.Marshal(doc_json)
		if err != nil {
			log.Fatalf("Error marshalling doc_json: %v", err)
		}

		// Unmarshal JSON into types.SearchableDocument
		err = json.Unmarshal(docBytes, &searchableDoc)
		if err != nil {
			log.Fatalf("Error unmarshalling doc_json: %v", err)
		}
		searchDocs = append(searchDocs, searchableDoc)
		// fmt.Println(searchableDoc)
	}
	// fmt.Println(docs)

	return searchDocs
}
