package cypherhelpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/helpers"
)

func PrintRepositoryShortSummaries(dbids []int) string {
	// Convert the slice to a string that Cypher can understand
	var quotedDbids []string
	for _, dbid := range dbids {
		quotedDbids = append(quotedDbids, fmt.Sprintf("\"%d\"", dbid))
	}

	// Join the quoted dbids with commas
	dbidString := strings.Join(quotedDbids, ", ")

	queryToExecute := fmt.Sprintf(
		`MATCH (r:Repository)-[]-(m:Module)-[:PART_OF_MODULE]-(p:Package)-[:CONTAINS]-(f:File)-[]-(fc:FileCommit)
	WHERE r.dbid IN [%s]
	WITH r.name AS RepositoryName, collect({CommitSummary: fc.summary, CommitAuthor: fc.authorName, Date: fc.commitDate}) AS Commits
RETURN {RepositoryName: RepositoryName, Commits: Commits} AS repository`,
		dbidString)

	result := cypher.InjectCypher(queryToExecute)
	responseStr, err := helpers.PrettyPrintYAMLInterface(result)
	if err != nil {
		return fmt.Sprintf("Error while performing PrintModuleShortSummaries: %s", err)
	}
	return responseStr
}

func PrintModuleShortSummaries(dbids []int) string {
	// Convert the slice to a string that Cypher can understand
	var quotedDbids []string
	for _, dbid := range dbids {
		quotedDbids = append(quotedDbids, fmt.Sprintf("\"%d\"", dbid))
	}

	// Join the quoted dbids with commas
	dbidString := strings.Join(quotedDbids, ", ")

	queryToExecute := fmt.Sprintf(
		`MATCH (m:Module)
	WHERE m.dbid IN [%s]
	RETURN {name: m.name , dbid: m.dbid, shortsummary: m.shortsummary} AS module`,
		dbidString)

	result := cypher.InjectCypher(queryToExecute)
	responseStr, err := helpers.PrettyPrintYAMLInterface(result)
	if err != nil {
		return fmt.Sprintf("Error while performing PrintModuleShortSummaries: %s", err)
	}
	return responseStr
}

func PrintPackageShortSummaries(dbids []int) string {
	// Convert the slice to a string that Cypher can understand
	var quotedDbids []string
	for _, dbid := range dbids {
		quotedDbids = append(quotedDbids, fmt.Sprintf("\"%d\"", dbid))
	}

	// Join the quoted dbids with commas
	dbidString := strings.Join(quotedDbids, ", ")

	queryToExecute := fmt.Sprintf(
		`MATCH (p:Package)-[:PART_OF_MODULE]-(m:Module)-[:HAS_MODULE]-(r:Repository)
        WHERE p.dbid IN [%s]
        RETURN {
            importpath: p.importpath, 
            dbid: p.dbid, 
            shortsummary: p.shortsummary,
            module: {name: m.name, shortsummary: m.shortsummary},
            repository: {name: r.name, shortsummary: r.shortsummary, dbid: r.dbid}
        } AS package`,
		dbidString)

	result := cypher.InjectCypher(queryToExecute)
	responseStr, err := helpers.PrettyPrintYAMLInterface(result)
	if err != nil {
		return fmt.Sprintf("Error while performing PrintModuleShortSummaries: %s", err)
	}
	return responseStr
}

func ExecuteQuery(queryString string) ([]interface{}, error) {
	cypherResult := cypher.InjectCypher(queryString)
	cypherResultJson, err := json.Marshal(cypherResult)
	if err != nil {
		return nil, fmt.Errorf("error in json.Marshal: %v", err)
	}

	var typedResult []interface{}
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return nil, fmt.Errorf("error in json.Unmarshal: %v", err)
	}

	if len(typedResult) == 0 {
		return nil, fmt.Errorf("unexpected result format: %v", typedResult)
	}

	endResult := typedResult
	return endResult, nil
}
