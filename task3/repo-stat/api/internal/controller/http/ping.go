package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// Ping godoc
//
//	@Summary		Ping services
//	@Description	Check status of all backend services (Processor, Subscriber)
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	dto.PingResponse
//	@Success		503	{object}	dto.PingResponse
//	@Router			/api/ping [get]
func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := ping.Execute(r.Context())

		response := dto.PingResponse{
			Status:   status.Status,
			Services: status.Services,
		}

		w.Header().Set("Content-Type", "application/json")

		if status.Status == domain.StatusDegraded {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
