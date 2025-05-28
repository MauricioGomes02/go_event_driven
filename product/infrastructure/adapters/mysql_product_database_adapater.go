package adapters

import (
	"go_event_driven/product/domain/entities"
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

func (adapter *MySqlProductDatabaseAdapter) Add(entity *entities.Product) (*entities.Product, error) {
	transaction, _error := adapter.databaseAdapter.Database.Begin()
	defer func() {
		_panic := recover()
		if _panic != nil {
			transaction.Rollback()
			panic(_panic)
		} else if _error != nil {
			transaction.Rollback()
		}
	}()

	insertCommand := "INSERT INTO products (product_id, product_name, product_description, product_amount, product_created_at, product_active) VALUES (?, ?, ?, ?, ?, ?)"
	_, _error = transaction.Exec(insertCommand, entity.Id, entity.Name, entity.Description, entity.Amount, entity.CreatedAt, entity.Active)
	if _error != nil {
		return nil, _error
	}

	_error = transaction.Commit()
	if _error != nil {
		return nil, _error
	}

	return entity, nil
}
