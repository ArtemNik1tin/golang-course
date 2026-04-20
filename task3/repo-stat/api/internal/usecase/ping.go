package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Ping struct {
	pingers []Pinger
}

func NewPing(pingers []Pinger) *Ping {
	return &Ping{
		pingers: pingers,
	}
}

func (u *Ping) Execute(ctx context.Context) domain.HealthStatus {
	services := make([]domain.ServiceStatus, len(u.pingers))
	healthStatus := domain.HealthStatus{
		Status:   domain.StatusOk,
		Services: services,
	}

	for i := 0; i < len(u.pingers); i++ {
		pingStatus := u.pingers[i].Ping(ctx)
		serviceName := u.pingers[i].Name()

		services[i] = domain.ServiceStatus{
			Status: pingStatus,
			Name:   serviceName,
		}

		if pingStatus == domain.PingStatusDown {
			healthStatus.Status = domain.StatusDegraded
		}
	}

	return healthStatus
}
