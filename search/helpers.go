package search

import (
	"fmt"
	"strings"

	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/types"
)

func SearchReporisitories(query string, limit int) ([]types.SearchableDocument, error) {
	results, err := SearchDocuments(query, limit, "doc_type = repository")
	if err != nil {
		return []types.SearchableDocument{}, err
	}
	return results, nil
}

func SearchDirectories(query string, dbid string, limit int) ([]types.SearchableDocument, error) {
	filterStr := "doc_type = directory"
	if dbid != "" {
		filterStr = "doc_type = directory AND repository_id = " + dbid
	}
	results, err := SearchDocuments(query, limit, filterStr)
	if err != nil {
		return []types.SearchableDocument{}, err
	}
	return results, nil
}

// func SearchModules(query string, dbid string, limit int) ([]types.SearchableDocument, error) {
// 	filterStr := "doc_type = module"
// 	if dbid != "" {
// 		filterStr = "doc_type = module AND repository_id = " + dbid
// 	}
// 	results, err := SearchDocuments(query, limit, filterStr)
// 	if err != nil {
// 		return []types.SearchableDocument{}, err
// 	}
// 	return results, nil
// }

func SearchFiles(query string, dbid string, limit int) ([]types.SearchableDocument, error) {
	filterStr := "doc_type = file"
	if dbid != "" {
		filterStr = "doc_type = file AND name != 'NON_EXISTING_FILE.go' AND package_id = " + dbid
	}
	results, err := SearchDocuments(query, limit, filterStr)
	if err != nil {
		return []types.SearchableDocument{}, err
	}
	return results, nil
}

func SearchContainers(query string, dbid string, limit int) ([]types.SearchableDocument, error) {
	filterStr := "doc_type = container"
	if dbid != "" {
		filterStr = "doc_type = container AND name != 'NON_EXISTING_FILE.go' AND repository_id = " + dbid
	}
	results, err := SearchDocuments(query, limit, filterStr)
	if err != nil {
		return []types.SearchableDocument{}, err
	}
	return results, nil
}

func SearchAllFiles(query string, limit int) ([]types.SearchableDocument, error) {
	filterStr := "doc_type = file AND name != 'NON_EXISTING_FILE.go'"
	results, err := SearchDocuments(query, limit, filterStr)
	if err != nil {
		return []types.SearchableDocument{}, err
	}
	return results, nil
}

func SearchFilesWithinRepository(query string, dbid string, limit int) ([]types.SearchableDocument, error) {
	filterStr := "doc_type = file"
	if dbid != "" {
		filterStr = "doc_type = file AND name != 'NON_EXISTING_FILE.go' AND repository_id = " + dbid
	}
	results, err := SearchDocuments(query, limit, filterStr)
	if err != nil {
		return []types.SearchableDocument{}, err
	}
	return results, nil
}

func SearchEntities(query string, dbid string, limit int) ([]types.SearchableDocument, error) {
	filterStr := "doc_type = entity"
	if dbid != "" {
		filterStr = "doc_type = entity AND file_id = " + dbid
	}
	results, err := SearchDocuments(query, limit, filterStr)
	if err != nil {
		return []types.SearchableDocument{}, err
	}

	return results, nil
}

func SearchAllDocuments(query string, limit int) ([]types.SearchableDocument, error) {
	results, err := SearchDocuments(query, limit)
	if err != nil {
		return []types.SearchableDocument{}, err
	}
	return results, nil
}

func SearchWSReporisitories(query string) (string, error) {
	results, err := SearchDocuments(query, 10, "doc_type = repository")
	if err != nil {
		return "", err
	}
	resultStr, err := helpers.PrettyPrintYAMLInterface(results)
	if err != nil {
		return "", err
	}
	return resultStr, nil
}

func SearchWSPackages(query string) (string, error) {
	results, err := SearchDocuments(query, 10, "doc_type = package")
	if err != nil {
		return "", err
	}
	resultStr, err := helpers.PrettyPrintYAMLInterface(results)
	if err != nil {
		return "", err
	}
	return resultStr, nil
}

func SearchWSModules(query string) (string, error) {
	results, err := SearchDocuments(query, 10, "doc_type = module")
	if err != nil {
		return "", err
	}
	resultStr, err := helpers.PrettyPrintYAMLInterface(results)
	if err != nil {
		return "", err
	}
	return resultStr, nil
}

func SearchWSFiles(query string) (string, error) {
	results, err := SearchDocuments(query, 10, "doc_type = file")
	if err != nil {
		return "", err
	}
	resultStr, err := helpers.PrettyPrintYAMLInterface(results)
	if err != nil {
		return "", err
	}
	return resultStr, nil
}

func SearchWSEntities(query string) (string, error) {
	results, err := SearchDocuments(query, 10, "doc_type = entity")
	if err != nil {
		return "", err
	}
	resultStr, err := helpers.PrettyPrintYAMLInterface(results)
	if err != nil {
		return "", err
	}
	return resultStr, nil
}

func SearchWSAllDocuments(query string) (string, error) {
	results, err := SearchDocuments(query, 10)
	if err != nil {
		return "", err
	}
	resultStr, err := helpers.PrettyPrintYAMLInterface(results)
	if err != nil {
		return "", err
	}
	return resultStr, nil
}

func GetTypeFromDbid(dbid string) string {
	docs, err := SearchAllDocuments(dbid, 1)
	if err != nil {
		return "error"
	}
	return docs[0].Type
}

func GetTypeFromSearchId(searchId string) string {
	parts := strings.Split(searchId, "-")
	if len(parts) < 2 {
		return "error"
	}
	return GetTypeFromDbid(parts[1])
}

func GetDbidFromSearchId(searchId string) string {
	parts := strings.Split(searchId, "-")
	if len(parts) < 2 {
		return "error"
	}
	return parts[1]
}

func GetRepoSearchIdFromSearchId(searchId string) (string, error) {
	docs, err := SearchAllDocuments(searchId, 1)
	if err != nil {
		return "", err
	}
	if len(docs) == 0 {
		return "", fmt.Errorf("no documents found for search id %s", searchId)
	}
	return fmt.Sprintf("repository-%d", docs[0].RepositoryID), nil
}
