package rest

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemNik1tin/distributed-github/gateway/internal/usecase"
)

type Handler struct {
	useCase *usecase.GetRepositoryUseCase
}

func NewHandler(useCase *usecase.GetRepositoryUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (handler Handler) GetRepository(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	ownerName := httpRequest.PathValue("owner")
	repoName := httpRequest.PathValue("repo")
	repo, err := handler.useCase.Execute(httpRequest.Context(), ownerName, repoName)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Internal Server Error"))
		return
	}

	response := map[string]interface{}{
		"name":        repo.Name,
		"description": repo.Description,
		"stars":       repo.Stars,
		"forks":       repo.Forks,
		"created_at":  repo.CreatedAt,
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(response)
}
