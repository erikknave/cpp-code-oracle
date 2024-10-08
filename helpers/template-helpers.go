package helpers

import (
	"bytes"
	"html/template"
)

func CreateStringFromTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("prompt").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func SafeSubstring(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}
