package search

import (
	"encoding/json"

	"github.com/erikknave/go-code-oracle/types"
	"github.com/meilisearch/meilisearch-go"
)

func SearchDocuments(query string, limit int, filters ...string) ([]types.SearchableDocument, error) {
	searchRequest := &meilisearch.SearchRequest{
		Query:  query,
		Filter: filters,
		Limit:  int64(limit),
	}

	searchResponse, err := Index.Search(query, searchRequest)
	if err != nil {
		return nil, err
	}

	var results []types.SearchableDocument
	for _, hit := range searchResponse.Hits {
		hitJSON, err := json.Marshal(hit)
		if err != nil {
			return nil, err
		}

		var doc types.SearchableDocument
		err = json.Unmarshal(hitJSON, &doc)
		if err != nil {
			return nil, err
		}
		results = append(results, doc)
	}

	return results, nil
}
