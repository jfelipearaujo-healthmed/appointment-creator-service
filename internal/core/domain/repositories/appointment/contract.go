package appointment_repository_contract

import (
	"context"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/core/domain/entities"
)

type Repository interface {
	GetByIDsAndDateTime(ctx context.Context, scheduleID uint, doctorID uint, dateTime time.Time) (*entities.Appointment, error)
	Create(ctx context.Context, appointment *entities.Appointment) (*entities.Appointment, error)
	Update(ctx context.Context, userID uint, appointment *entities.Appointment) (*entities.Appointment, error)
	Delete(ctx context.Context, userID uint, appointmentID uint) error
}
