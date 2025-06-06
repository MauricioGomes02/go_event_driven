package mysql

import (
	"errors"
	"fmt"
	"go_event_driven/product/domain/entities"
	"go_event_driven/product/domain/ports"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql" // Register driver
)

type MySqlProductDatabaseAdapter struct {
	databaseAdapter *MySqlDatabaseAdapter
}

var (
	once sync.Once
)

func NewMySqlProductDatabaseAdapter(databaseAdapter *MySqlDatabaseAdapter) *MySqlProductDatabaseAdapter {

	return &MySqlProductDatabaseAdapter{
		databaseAdapter: databaseAdapter,
	}
}

func (adapter *MySqlProductDatabaseAdapter) AddWithOutboxEvent(entity *entities.Product, event *entities.OutboxEvent) (*entities.Product, error) {
	transaction, _error := adapter.databaseAdapter.Database.Begin()
	if _error != nil {
		return nil, _error
	}
	defer func() {
		_panic := recover()
		if _panic != nil {
			transaction.Rollback()
			panic(_panic)
		} else if _error != nil {
			transaction.Rollback()
		}
	}()

	insertProductCommand := `
		INSERT INTO products (
			product_id, 
			product_name, 
			product_description, 
			product_amount, 
			product_created_at, 
			product_active
		) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, _error = transaction.Exec(
		insertProductCommand,
		entity.Id,
		entity.Name,
		entity.Description,
		entity.Amount,
		entity.CreatedAt,
		entity.Active)

	if _error != nil {
		transaction.Rollback()
		return nil, _error
	}

	insertOutboxEventCommand := `
		INSERT INTO outbox_events (
			id, 
			aggregate_id, 
			aggregate_type, 
			event_type, 
			payload,
			status, 
			retries, 
			error_message, 
			created_at, 
			sent_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, _error = transaction.Exec(
		insertOutboxEventCommand,
		event.Id,
		event.AggregateId,
		event.AggregateType,
		event.EventType,
		event.Payload,
		event.Status,
		event.Retries,
		event.ErrorMessage,
		event.CreatedAt,
		event.SentAt)

	if _error != nil {
		transaction.Rollback()
		return nil, _error
	}

	_error = transaction.Commit()

	return entity, _error
}

func (adapter *MySqlProductDatabaseAdapter) GetOutboxEvents(criterion ports.CriterionPort) ([]entities.OutboxEvent, error) {
	sqlCriterion, ok := criterion.(*SqlCriterion)

	if !ok {
		return nil, errors.New("Could not convert to SqlCriterion")
	}

	query, arguments := sqlCriterion.ToSql()
	formattedQuery := fmt.Sprintf("SELECT * FROM outbox_events WHERE %s", query)
	rows, _error := adapter.databaseAdapter.Database.Query(formattedQuery, arguments...)

	if _error != nil {
		return nil, _error
	}

	defer rows.Close()

	var _entities []entities.OutboxEvent
	for rows.Next() {
		var entity entities.OutboxEvent
		_error = rows.Scan(
			&entity.Id,
			&entity.AggregateId,
			&entity.AggregateType,
			&entity.EventType,
			&entity.Payload,
			&entity.Status,
			&entity.Retries,
			&entity.ErrorMessage,
			&entity.CreatedAt,
			&entity.SentAt,
		)

		if _error != nil {
			return nil, _error
		}

		_entities = append(_entities, *entities.NewOutboxEvent(
			entity.Id,
			entity.AggregateId,
			entity.AggregateType,
			entity.EventType,
			entity.Payload,
			entity.Status,
			entity.Retries,
			entity.ErrorMessage,
			entity.CreatedAt,
			entity.SentAt,
		))
	}

	return _entities, nil
}

func (adapter *MySqlProductDatabaseAdapter) UpdateOutboxEvent(entity entities.OutboxEvent) error {
	transaction, _error := adapter.databaseAdapter.Database.Begin()
	if _error != nil {
		return _error
	}
	defer func() {
		_panic := recover()
		if _panic != nil {
			transaction.Rollback()
			panic(_panic)
		} else if _error != nil {
			transaction.Rollback()
		}
	}()

	updateCommand := `
		UPDATE outbox_events 
		SET %s
		WHERE id = ?
	`

	_set := []string{}
	_values := []any{}

	for property, value := range entity.GetUpdatedFields() {
		_set = append(_set, fmt.Sprintf("%s = ?", FieldMapping["OutboxEvent"][property]))
		_values = append(_values, value)
	}

	updateCommand = fmt.Sprintf(updateCommand, strings.Join(_set, ", "))
	_values = append(_values, entity.Id)

	_, _error = transaction.Exec(
		updateCommand,
		_values...,
	)

	if _error != nil {
		transaction.Rollback()
		return _error
	}

	return transaction.Commit()
}
