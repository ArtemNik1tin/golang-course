package dto

import (
	"repo-stat/api/internal/domain"
)

type PingResponse struct {
	Status   string                 `json:"status"`
	Services []domain.ServiceStatus `json:"services"`
}
