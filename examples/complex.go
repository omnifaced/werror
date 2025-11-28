package examples

import (
	"fmt"
	"github.com/omnifaced/werror"
	"strings"
)

func readFile(filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("empty filename")
	}

	return "user: admin\nrole: superuser", nil
}

func validateContent(content string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("empty content")
	}

	if !strings.Contains(content, "user:") {
		return "", fmt.Errorf("invalid format")
	}

	return content, nil
}

func normalizeContent(content string) (string, error) {
	normalized := strings.TrimSpace(content)

	if len(normalized) == 0 {
		return "", fmt.Errorf("content is empty after normalization")
	}

	return normalized, nil
}

func complexExample() {
	result := werror.Wrap(readFile("config.yaml")).
		OnSuccess(func(content string) {
			fmt.Println("file read successfully")
		}).
		ThenFn(validateContent).
		OnSuccess(func(content string) {
			fmt.Println("content validated")
		}).
		ThenFn(normalizeContent).
		OnError(func(err error) {
			fmt.Println("processing failed:", err)
		}).
		Always(func() {
			fmt.Println("operation completed")
		})

	if result.IsOk() {
		content := result.Value()
		fmt.Println("final content:", content)
	} else {
		defaultContent := "user: guest\nrole: viewer"
		fmt.Println("using default:", defaultContent)
	}
}
