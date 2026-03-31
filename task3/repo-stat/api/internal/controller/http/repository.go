package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

var regularExpression = regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)

func NewGetRepositoryInfoHandler(log *slog.Logger, getRepositoryUseCase *usecase.GetRepositoryUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")

		matches := regularExpression.FindStringSubmatch(rawURL)
		if len(matches) < 3 {
			log.Warn("invalid url provided", "url", rawURL)
			http.Error(w, "Invalid GitHub URL", http.StatusBadRequest)
			return
		}

		ownerName, repoName := matches[1], matches[2]

		repository, err := getRepositoryUseCase.Execute(r.Context(), ownerName, repoName)
		if err != nil {
			log.Error("usecase error", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := dto.Repository{
			Name:        repository.Name,
			Description: repository.Description,
			Stars:       repository.Stars,
			Forks:       repository.Forks,
			CreatedAt:   repository.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
