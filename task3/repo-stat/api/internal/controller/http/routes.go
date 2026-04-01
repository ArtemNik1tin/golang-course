package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/docs"
	"repo-stat/api/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func AddRoutes(
	mux *http.ServeMux,
	log *slog.Logger,
	pingUC *usecase.Ping,
	repoUC *usecase.GetRepositoryUseCase,
) {
	docs.SwaggerInfo.Title = "Repo Stat API"
	docs.SwaggerInfo.Description = "API Gateway for repo-stat microservice"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:28080"
	docs.SwaggerInfo.BasePath = "/"

	mux.Handle("GET /api/ping", NewPingHandler(log, pingUC))
	mux.Handle("GET /api/repositories/info", NewGetRepositoryInfoHandler(log, repoUC))
	mux.Handle("/swagger/doc.json", httpSwagger.WrapHandler)
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"),
	))
}
