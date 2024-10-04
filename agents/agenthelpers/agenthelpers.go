package agenthelpers

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/erikknave/go-code-oracle/maps"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v2"
)

//go:embed oracle-agents.xlsx
var embeddedFile []byte

func InitAgentDescriptions() {

	fileBytes := embeddedFile
	reader := bytes.NewReader(fileBytes)

	f, err := excelize.OpenReader(reader)
	if err != nil {
		log.Fatalf("Failed to open the Excel file: %v", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Failed to get rows from the sheet: %v", err)
	}

	for idx, row := range rows {
		// Skip the header row
		if idx == 0 {
			continue
		}
		if len(row) < 5 {
			log.Printf("Skipping incomplete row: %v", row)
			continue
		}
		name := row[0]
		maps.AgentDescriptions.Store(name, types.AgentDescription{
			Name:           row[0],
			Caller:         row[1],
			SystemMessage:  row[2],
			PromptTemplate: row[3],
			Model:          row[4],
		})
		embeddedFile = nil
	}
}

// ParseJSON takes a string as input and attempts to parse JSON data from it.
func ParseJSON(input string, v interface{}) error {
	startDelimiters := []string{"```json", "```"}
	endDelimiter := "```"

	scanner := bufio.NewScanner(strings.NewReader(input))
	var jsonStrBuilder strings.Builder
	inCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if inCodeBlock {
			if strings.TrimSpace(line) == endDelimiter {
				break
			}
			jsonStrBuilder.WriteString(line)
			jsonStrBuilder.WriteRune('\n')
		} else {
			for _, startDelimiter := range startDelimiters {
				if strings.HasPrefix(strings.TrimSpace(line), startDelimiter) {
					inCodeBlock = true
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	if !inCodeBlock {
		// No code block found, try to parse the entire input as JSON
		return json.Unmarshal([]byte(input), v)
	}

	return json.Unmarshal([]byte(jsonStrBuilder.String()), v)
}

func ParseYAML(input string, v interface{}) error {
	startDelimiters := []string{"```yaml", "```"}
	endDelimiter := "```"

	scanner := bufio.NewScanner(strings.NewReader(input))
	var yamlStrBuilder strings.Builder
	inCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if inCodeBlock {
			if strings.TrimSpace(line) == endDelimiter {
				break
			}
			yamlStrBuilder.WriteString(line)
			yamlStrBuilder.WriteRune('\n')
		} else {
			for _, startDelimiter := range startDelimiters {
				if strings.HasPrefix(strings.TrimSpace(line), startDelimiter) {
					inCodeBlock = true
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	if !inCodeBlock {
		// No code block found, try to parse the entire input as YAML
		return yaml.Unmarshal([]byte(input), v)
	}

	return yaml.Unmarshal([]byte(yamlStrBuilder.String()), v)
}
