package chromaclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

func PerformQuery(queryString string, numDocuments int, store chroma.Store) []types.SearchableDocument {
	ctx := context.Background()
	var searchDocs []types.SearchableDocument

	docs, err := store.SimilaritySearch(ctx, queryString, numDocuments)
	if err != nil {
		log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
	}
	// var dbids []int
	for _, doc := range docs {
		dbidInterface := doc.Metadata["dbid"]
		dbid, ok := dbidInterface.(float64)
		if !ok {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		idInterface := doc.Metadata["id"]
		id, ok := idInterface.(string)
		if !ok {
			log.Printf("error while performing the following: PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		latestCommitStrInterface := doc.Metadata["latest_commit"]
		latestCommitStr, ok := latestCommitStrInterface.(string)
		if !ok {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		pathInterface := doc.Metadata["path"]
		pathStr, ok := pathInterface.(string)
		if !ok {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		shortSummaryInterface := doc.Metadata["short_summary"]
		shortSummary, ok := shortSummaryInterface.(string)
		if !ok {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		nameInterface := doc.Metadata["name"]
		name, ok := nameInterface.(string)
		if !ok {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		if strings.Contains(name, "DUMMY_REPO") {
			continue
		}
		summaryInterface := doc.Metadata["summary"]
		summary, ok := summaryInterface.(string)
		if !ok {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		typeInterface := doc.Metadata["type"]
		docType, ok := typeInterface.(string)
		if !ok {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		const layout = time.RFC3339 // Define the layout according to the date string format

		latestCommit, err := time.Parse(time.RFC3339, latestCommitStr)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return searchDocs
		}
		authorsStringInterface := doc.Metadata["authors"]
		authorsString, ok := authorsStringInterface.(string)
		if !ok {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		var authors []string
		err = json.Unmarshal([]byte(authorsString), &authors)
		if err != nil {
			log.Printf("error while performing PerformQuery: %v\nChroma Query Store: %v\n", err, store)
			continue
		}
		searchDoc := types.SearchableDocument{
			Dbid:         int(dbid),
			ID:           id,
			LatestCommit: latestCommit,
			Path:         pathStr,
			ShortSummary: shortSummary,
			Summary:      summary,
			Type:         docType,
			Name:         name,
			Authors:      authors,
		}
		searchDocs = append(searchDocs, searchDoc)

		// dbids = append(dbids, int(dbid))
	}
	return searchDocs
}

func PerformRepositoryQuery(queryString string, numDocuments int) []types.SearchableDocument {
	return PerformQuery(queryString, numDocuments, RepositoryStore)
}

func PerformModuleQuery(queryString string, numDocuments int) []types.SearchableDocument {
	return PerformQuery(queryString, numDocuments, ModuleStore)
}

func PerformPackageQuery(queryString string, numDocuments int) []types.SearchableDocument {
	return PerformQuery(queryString, numDocuments, PackageStore)
}

func PerformFileQuery(queryString string, numDocuments int) []types.SearchableDocument {
	return PerformQuery(queryString, numDocuments, FileStore)
}

func PerformEntityQuery(queryString string, numDocuments int) []types.SearchableDocument {
	return PerformQuery(queryString, numDocuments, EntityStore)
}
