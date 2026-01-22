package file

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ValidateWorkspacePath(path string, workspace string) error {
	if !filepath.IsAbs(path) {
		return errors.New("path must be absolute")
	}

	cleanWorkspace := filepath.Clean(workspace)
	cleanPath := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		cleanWorkspace = strings.ToLower(cleanWorkspace)
		cleanPath = strings.ToLower(cleanPath)
	}

	rel, err := filepath.Rel(cleanWorkspace, cleanPath)
	if err != nil {
		return errors.New("path must be within workspace")
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return errors.New("path must be within workspace")
	}
	return nil
}

func Read(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Write(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func List(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	results := make([]string, 0, len(entries))
	for _, entry := range entries {
		results = append(results, entry.Name())
	}
	return results, nil
}

func Search(path string, query string) ([]string, error) {
	if query == "" {
		return nil, errors.New("query must not be empty")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	var matches []string
	for _, line := range lines {
		if strings.Contains(line, query) {
			matches = append(matches, line)
		}
	}
	return matches, nil
}

func Replace(path string, search string, replace string) (int, error) {
	if search == "" {
		return 0, errors.New("search must not be empty")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	updated := strings.ReplaceAll(string(content), search, replace)
	count := strings.Count(string(content), search)
	if err := os.WriteFile(path, []byte(updated), fs.FileMode(0644)); err != nil {
		return 0, err
	}
	return count, nil
}
