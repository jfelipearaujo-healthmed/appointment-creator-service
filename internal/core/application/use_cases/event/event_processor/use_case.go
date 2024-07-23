package event_processor_uc

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/core/domain/events"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/core/domain/repositories/appointment"
	event_repository_contract "github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/core/domain/repositories/event"
	event_processor_contract "github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/core/domain/use_cases/event/event_processor"
	"github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-creator-service/internal/external/queue"
)

var (
	OUTCOME_APPOINTMENT_CREATED        string = "appointment created successfully"
	OUTCOME_APPOINTMENT_ALREADY_BOOKED string = "appointment already booked"
)

type useCase struct {
	eventRepository       event_repository_contract.Repository
	appointmentRepository appointment_repository_contract.Repository
}

func NewUseCase(
	eventRepository event_repository_contract.Repository,
	appointmentRepository appointment_repository_contract.Repository,
) event_processor_contract.UseCase {
	return &useCase{
		eventRepository:       eventRepository,
		appointmentRepository: appointmentRepository,
	}
}

func (uc *useCase) Handle(ctx context.Context, messageID string, message queue.Message) error {
	slog.InfoContext(ctx, "event received", "message_id", messageID)

	event := new(entities.Event)

	messageJson, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(messageJson, event); err != nil {
		return err
	}

	eventMap := map[queue.EventType]func(ctx context.Context, messageID string, appointment *entities.Appointment) error{
		events.CreateAppointment: uc.HandleCreateAppointment,
	}

	slog.InfoContext(ctx, "checking handler for event type", "event_type", event.EventType)

	if handler, ok := eventMap[event.EventType]; ok {
		appointment := new(entities.Appointment)

		if err := json.Unmarshal([]byte(event.Data), appointment); err != nil {
			return err
		}

		err = handler(ctx, messageID, appointment)
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "event processed successfully", "message_id", messageID)

		return nil
	}

	slog.ErrorContext(ctx, "event handler not found", "message_id", messageID)

	return nil
}

func (uc *useCase) HandleCreateAppointment(ctx context.Context, messageID string, appointment *entities.Appointment) error {
	slog.InfoContext(ctx, "handling appointment creation", "message_id", messageID)

	slog.InfoContext(ctx, "loading event for message received", "message_id", messageID)

	event, err := uc.eventRepository.GetByMessageID(ctx, messageID)
	if err != nil {
		slog.ErrorContext(ctx, "error loading event", "message_id", messageID, "error", err)
		return err
	}

	if event.Outcome != nil {
		slog.WarnContext(ctx, "event already processed", "message_id", messageID)
		return nil
	}

	slog.InfoContext(ctx, "checking if appointment already exists", "message_id", messageID)

	existingAppointment, err := uc.appointmentRepository.GetByIDsAndDateTime(ctx, appointment.ScheduleID, appointment.DoctorID, appointment.DateTime)
	if err != nil && !app_error.IsAppError(err) {
		slog.ErrorContext(ctx, "error checking if appointment already exists", "message_id", messageID, "error", err)
		return err
	}

	if existingAppointment == nil {
		slog.InfoContext(ctx, "appointment does not exist, creating...", "message_id", messageID)

		appointment.EventID = event.ID
		appointment.Status = entities.WaitingForConfirmation

		if _, err := uc.appointmentRepository.Create(ctx, appointment); err != nil {
			slog.ErrorContext(ctx, "error creating appointment", "message_id", messageID, "error", err)
			return err
		}

		event.Outcome = &OUTCOME_APPOINTMENT_CREATED

		return uc.updateEvent(ctx, messageID, event)
	}

	slog.WarnContext(ctx, "appointment date and time already booked", "message_id", messageID)

	event.Outcome = &OUTCOME_APPOINTMENT_ALREADY_BOOKED

	return uc.updateEvent(ctx, messageID, event)
}

func (uc *useCase) updateEvent(ctx context.Context, messageID string, event *entities.Event) error {
	slog.InfoContext(ctx, "updating event...", "message_id", messageID)

	if _, err := uc.eventRepository.Update(ctx, event); err != nil {
		slog.ErrorContext(ctx, "error updating event", "message_id", messageID, "error", err)
		return err
	}

	slog.InfoContext(ctx, "event updated successfully", "message_id", messageID)

	return nil
}
