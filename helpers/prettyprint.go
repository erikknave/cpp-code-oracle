package helpers

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

func PrettyPrintJSONString(rawJSON string) (string, error) {
	var prettyJSON map[string]interface{}
	err := json.Unmarshal([]byte(rawJSON), &prettyJSON)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling raw JSON: %v", err)
	}

	prettyJSONBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling pretty JSON: %v", err)
	}

	return string(prettyJSONBytes), nil
}

func PrettyPrintJSONInterface(rawJSON interface{}) (string, error) {
	prettyJSONBytes, err := json.MarshalIndent(rawJSON, "", "  ")
	if err != nil {
		fmt.Printf("error marshaling pretty JSON: %v", err)
		return "", err
	}
	return string(prettyJSONBytes), nil
}

func PrettyPrintYAMLString(rawYAML string) (string, error) {
	var prettyYAML map[string]interface{}
	err := yaml.Unmarshal([]byte(rawYAML), &prettyYAML)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling raw YAML: %v", err)
	}

	prettyYAMLBytes, err := yaml.Marshal(prettyYAML)
	if err != nil {
		return "", fmt.Errorf("error marshaling pretty YAML: %v", err)
	}

	return string(prettyYAMLBytes), nil
}

func PrettyPrintYAMLInterface(rawInterface interface{}) (string, error) {
	prettyYAMLBytes, err := yaml.Marshal(rawInterface)
	if err != nil {
		fmt.Printf("error marshaling pretty YAML: %v", err)
		return "", err

	}
	return string(prettyYAMLBytes), nil
}
