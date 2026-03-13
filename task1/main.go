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

	printRepoInfo(repoInfo)

	return nil
}

func fetchRepository(path string) (*Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: %s", resp.Status)
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

func printRepoInfo(repo *Repository) {
	if repo == nil {
		return
	}
	fmt.Println("--------------------------------------------------")
	fmt.Printf("%sRepository:%s  %s\n", colorCyan, colorReset, repo.Name)
	fmt.Printf("%sDescription:%s %s\n", colorCyan, colorReset, repo.Description)
	fmt.Printf("%sStars:%s       %d\n", colorCyan, colorYellow, repo.Stars)
	fmt.Printf("%sForks:%s       %d\n", colorCyan, colorReset, repo.Forks)
	fmt.Printf("%sCreated at:%s  %s\n", colorCyan, colorReset, repo.CreatedAt)
	fmt.Println("--------------------------------------------------")
}
