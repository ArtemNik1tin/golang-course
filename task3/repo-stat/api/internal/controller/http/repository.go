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

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewGetRepositoryInfoHandler(log *slog.Logger, getRepositoryUseCase *usecase.GetRepositoryUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")

		if rawURL == "" {
			log.Warn("url parameter is missing")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "url parameter is required"})
			return
		}

		matches := regularExpression.FindStringSubmatch(rawURL)
		if len(matches) < 3 {
			log.Warn("invalid url provided", "url", rawURL)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid GitHub URL"})
			return
		}

		ownerName, repoName := matches[1], matches[2]

		repository, err := getRepositoryUseCase.Execute(r.Context(), ownerName, repoName)
		if err != nil {
			log.Error("usecase error", "err", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal Server Error"})
			return
		}

		response := dto.Repository{
			FullName:    repository.Name,
			Description: repository.Description,
			Stars:       repository.Stars,
			Forks:       repository.Forks,
			CreatedAt:   repository.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
