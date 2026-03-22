// Package main provides a CLI tool for analyzing GitHub repositories.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func (r *Repository) String() string {
	if r == nil {
		return "nil repository"
	}

	return fmt.Sprintf(
		"--------------------------------------------------\n"+
			"%sRepository:%s  %s\n"+
			"%sDescription:%s %s\n"+
			"%sStars:%s       %d\n"+
			"%sForks:%s       %d\n"+
			"%sCreated at:%s  %s\n"+
			"--------------------------------------------------",
		colorCyan, colorReset, r.Name,
		colorCyan, colorReset, r.Description,
		colorCyan, colorYellow, r.Stars,
		colorCyan, colorReset, r.Forks,
		colorCyan, colorReset, r.CreatedAt,
	)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: ./repo-info <owner/repository>")
	}
	repoPath := os.Args[1]

	repoInfo, err := fetchRepository(repoPath)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	fmt.Println(repoInfo)

	return nil
}

func fetchRepository(path string) (*Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close body: %w", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: %s", resp.Status)
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}
