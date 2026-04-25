package domain

type PingStatus string

const (
	PingStatusUp   PingStatus = "up"
	PingStatusDown PingStatus = "down"

	StatusOk       string = "ok"
	StatusDegraded string = "degraded"
)

type ServiceStatus struct {
	Name   string
	Status PingStatus
}

type HealthStatus struct {
	Status   string
	Services []ServiceStatus
}
