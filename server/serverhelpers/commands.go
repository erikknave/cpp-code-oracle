package serverhelpers

import (
	"fmt"
	"log"
	"strings"

	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/server/chromaclient"
	"github.com/erikknave/go-code-oracle/server/pgvector"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
)

func PerformWSSearch(searchStr string, c *websocket.Conn) {
	words := strings.Fields(searchStr)

	// Check if there are at least three words
	if len(words) >= 3 {
		switchStr := words[1]
		thirdWordOnwards := words[2:]
		queryString := strings.Join(thirdWordOnwards, " ")
		responseStr := ""
		var err error
		switch switchStr {
		case "repositories":
			responseStr, err = search.SearchWSReporisitories(queryString)
			if err != nil {
				responseStr = fmt.Sprintf("Error searching repositories: %v", err)
				fmt.Println(responseStr)
			}
			SendServerMessage(c, responseStr)
		case "packages":
			responseStr, err = search.SearchWSPackages(queryString)
			if err != nil {
				responseStr = fmt.Sprintf("Error searching packages: %v", err)
				fmt.Println(responseStr)
			}
			SendServerMessage(c, responseStr)
		case "modules":
			responseStr, err = search.SearchWSModules(queryString)
			if err != nil {
				responseStr = fmt.Sprintf("Error searching modules: %v", err)
				fmt.Println(responseStr)
			}
			SendServerMessage(c, responseStr)
		case "files":
			responseStr, err = search.SearchWSFiles(queryString)
			if err != nil {
				responseStr = fmt.Sprintf("Error searching files: %v", err)
				fmt.Println(responseStr)
			}
			SendServerMessage(c, responseStr)
		case "entities":
			responseStr, err = search.SearchWSEntities(queryString)
			if err != nil {
				responseStr = fmt.Sprintf("Error searching entities: %v", err)
				fmt.Println(responseStr)
			}
			SendServerMessage(c, responseStr)
		case "all":
			responseStr, err = search.SearchWSAllDocuments(queryString)
			if err != nil {
				responseStr = fmt.Sprintf("Error searching documents: %v", err)
				fmt.Println("Error searching documents:", err)
			}
			SendServerMessage(c, responseStr)
		default:
			responseStr := fmt.Sprintf("Unknown search type: %v\nTip: /search <type> <query>\nTypes are repository, module, package, file, entity", switchStr)
			SendServerMessage(c, responseStr)
		}

	} else {
		responseStr := "Not enough words in the search string\nTip: /search <type> <query>"
		SendServerMessage(c, responseStr)
	}
}

func PerformSearch(searchStr string) ([]types.SearchableDocument, error) {
	words := strings.Fields(searchStr)

	// Check if there are at least three words

	var response []types.SearchableDocument
	var err error
	if len(words) >= 2 {
		switchStr := words[0]
		thirdWordOnwards := words[1:]
		queryString := strings.Join(thirdWordOnwards, " ")
		switch switchStr {
		case "/repositories":
			response, err = search.SearchReporisitories(queryString, 100)
			if err != nil {
				log.Fatalf("Error searching repositories: %v", err)
				return nil, err
			}
			return response, nil
		case "/directories":
			response, err = search.SearchDirectories(queryString, "", 100)
			if err != nil {
				log.Fatalf("Error searching packages: %v", err)
				return nil, err
			}
			return response, nil
		case "/containers":
			response, err = search.SearchContainers(queryString, "", 100)
			if err != nil {
				log.Fatalf("Error searching modules: %v", err)
				return nil, err
			}
			return response, nil
		case "/files":
			response, err = search.SearchFiles(queryString, "", 100)
			if err != nil {
				log.Fatalf("Error searching files: %v", err)
				return nil, err
			}
			return response, nil
		case "/entities":
			response, err = search.SearchEntities(queryString, "", 100)
			if err != nil {
				log.Fatalf("Error searching entities: %v", err)
				return nil, err
			}
			return response, nil
		case "/all":
			response, err = search.SearchAllDocuments(queryString, 200)
			if err != nil {
				log.Fatalf("Error searching documents: %v", err)
				return nil, err
			}
			return response, nil
		case "/embeddings":
			secondWord := words[1]

			switch secondWord {
			case "packages":
				queryString := strings.Join(words[2:], " ")
				response = chromaclient.PerformPackageQuery(queryString, 10)
				return response, nil
			case "files":
				queryString := strings.Join(words[2:], " ")
				response = chromaclient.PerformFileQuery(queryString, 10)
				return response, nil
			default:
				// response = chromaclient.PerformRepositoryQuery(queryString, 10)
				response = pgvector.PerformSearch(queryString, 10)
				return response, nil
			}
		default:
			response, err = search.SearchReporisitories(queryString, 100)
			if err != nil {
				log.Fatalf("Error searching documents: %v", err)
				return nil, err
			}
			return response, nil
		}

	} else {
		queryString := words[0]
		response, err = search.SearchReporisitories(queryString, 100)
		if err != nil {
			log.Fatalf("Error searching documents: %v", err)
			return nil, err
		}
		return response, nil

	}
}
