package templates

import (
	"fmt"
	"strings"
	"time"

	"github.com/erikknave/go-code-oracle/types"
)

func FormatAuthorString(authors []string) string {
	if len(authors) == 0 {
		return "No recent authors"
	}
	if len(authors) == 1 {
		return authors[0]
	}
	return strings.Join(authors[:len(authors)-1], ", ") + " and " + authors[len(authors)-1]
}

func FormatDateString(date time.Time) string {
	dateStr := date.Format("2006-01-02")
	if dateStr == "0001-01-01" {
		return "No date available"
	}
	return dateStr
}

func FormatSummaryString(searchResult *types.SearchableDocument) string {
	if searchResult.ShortSummary == "" {
		return searchResult.Summary
	}
	return searchResult.ShortSummary
}

func FormatRepoLastName(repoName string) string {
	words := strings.Split(repoName, "/")
	return words[len(words)-1]
}

func FormatIsoDateString(dateStr string) string {
	if dateStr == "" {
		return "No date available"
	}
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "Unformatted date: " + dateStr
	}
	formattadDate := FormatDateString(date)
	return formattadDate
}

func GetRelativePackagePath(repoPath, packagePath string) (string, error) {
	if packagePath == "" {
		return "/", nil
	}
	if strings.HasPrefix(packagePath, repoPath) {
		relativePath := strings.TrimPrefix(packagePath, repoPath)
		relativePath = strings.TrimPrefix(relativePath, "/") // Remove leading slash if present
		return relativePath, nil
	}

	// If the packagePath doesn't start with repoPath, take whatever is after the second /
	parts := strings.SplitN(packagePath, "/", 3)
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid package path format")
	}
	return parts[2], nil
}

func GetAfterSecondSlash(path string) string {
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 3 {
		return "Unknown"
	}
	return parts[2]
}
