package rest

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ArtemNik1tin/distributed-github/gateway/internal/usecase"
)

type RepositoryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type Handler struct {
	useCase *usecase.GetRepositoryUseCase
}

func NewHandler(useCase *usecase.GetRepositoryUseCase) *Handler {
	return &Handler{useCase: useCase}
}

// GetRepository godoc
// @Summary Get repository info
// @Description Get information about a GitHub repository
// @Tags repositories
// @Param owner path string true "Repository owner"
// @Param repo path string true "Repository name"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Repository not found"
// @Failure 502 {string} string "Upstream error"
// @Failure 500 {string} string "Internal Server Error"
// @Router /repos/{owner}/{repo} [get]
func (handler Handler) GetRepository(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	ownerName := httpRequest.PathValue("owner")
	repoName := httpRequest.PathValue("repo")
	repo, err := handler.useCase.Execute(httpRequest.Context(), ownerName, repoName)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				responseWriter.WriteHeader(http.StatusNotFound)
				responseWriter.Write([]byte("Repository not found"))
				return
			default:
				responseWriter.WriteHeader(http.StatusBadGateway)
				responseWriter.Write([]byte("Upstream error"))
				return
			}
		}
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Internal Server Error"))
		return
	}

	response := RepositoryResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(response)
}
