package entities

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	Id            uuid.UUID
	AggregateId   uuid.UUID
	AggregateType string
	EventType     string
	Payload       []byte
	Status        string
	Retries       int
	ErrorMessage  *string
	CreatedAt     time.Time
	SentAt        *time.Time
	updatedFields map[string]any
}

func NewOutboxEvent(
	id uuid.UUID,
	aggregateId uuid.UUID,
	aggregateType string,
	eventType string,
	payload []byte,
	status string,
	retries int,
	errorMessage *string,
	createdAt time.Time,
	sentAt *time.Time) *OutboxEvent {

	return &OutboxEvent{
		Id:            id,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		EventType:     eventType,
		Payload:       payload,
		Status:        status,
		Retries:       retries,
		ErrorMessage:  errorMessage,
		CreatedAt:     createdAt,
		SentAt:        sentAt,
		updatedFields: make(map[string]any),
	}
}

func (entity *OutboxEvent) UpdateStatus(newStatus string) {
	entity.Status = newStatus
	entity.updatedFields["Status"] = newStatus
}

func (entity *OutboxEvent) UpdateErrorMessage(newErrorMessage string) {
	entity.ErrorMessage = &newErrorMessage
	entity.updatedFields["ErrorMessage"] = newErrorMessage
}

func (entity *OutboxEvent) UpdateSentAt(newSentAt time.Time) {
	entity.SentAt = &newSentAt
	entity.updatedFields["SentAt"] = newSentAt
}

func (entity *OutboxEvent) UpdateRetries(newRetries int) {
	entity.Retries = newRetries
	entity.updatedFields["Retries"] = newRetries
}

func (entity *OutboxEvent) GetUpdatedFields() map[string]any {
	return entity.updatedFields
}
