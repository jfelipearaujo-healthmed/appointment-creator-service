package events

import "github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/external/queue"

const (
	CreateAppointment queue.EventType = "create_appointment"
	UpdateAppointment queue.EventType = "update_appointment"
)
