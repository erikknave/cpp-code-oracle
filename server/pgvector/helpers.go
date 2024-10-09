package pgvector

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func performSearchWithFilter(queryString string, numDocuments int, filter map[string]interface{}) []types.SearchableDocument {
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

	var docs []schema.Document
	if filter != nil {
		// Apply filter if provided
		docs, err = store.SimilaritySearch(ctx, queryString, numDocuments, vectorstores.WithFilters(filter))
	} else {
		// No filter, normal search
		docs, err = store.SimilaritySearch(ctx, queryString, numDocuments)
	}
	if err != nil {
		log.Fatal(err)
	}

	searchDocs := []types.SearchableDocument{}
	for _, doc := range docs {
		metadata := doc.Metadata
		doc_json := metadata["searchable_document"]
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
	}

	return searchDocs
}

func PerformQuery(queryString string, numDocuments int) []types.SearchableDocument {
	// Call performSearchWithFilter without any filter
	return performSearchWithFilter(queryString, numDocuments, nil)
}

func PerformRepositoryQuery(queryString string, numDocuments int) []types.SearchableDocument {
	// Call performSearchWithFilter with a filter for repositories
	filter := map[string]interface{}{"doc_type": "repository"}
	return performSearchWithFilter(queryString, numDocuments, filter)
}

func PerformDirectoryQuery(queryString string, numDocuments int) []types.SearchableDocument {
	// Call performSearchWithFilter with a filter for repositories
	filter := map[string]interface{}{"doc_type": "directory"}
	return performSearchWithFilter(queryString, numDocuments, filter)
}

func PerformFileQuery(queryString string, numDocuments int) []types.SearchableDocument {
	// Call performSearchWithFilter with a filter for repositories
	filter := map[string]interface{}{"doc_type": "file"}
	return performSearchWithFilter(queryString, numDocuments, filter)
}

func PerformContainerQuery(queryString string, numDocuments int) []types.SearchableDocument {
	// Call performSearchWithFilter with a filter for repositories
	filter := map[string]interface{}{"doc_type": "container"}
	return performSearchWithFilter(queryString, numDocuments, filter)
}
