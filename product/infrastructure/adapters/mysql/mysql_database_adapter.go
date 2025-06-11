package mysql

import (
	"database/sql"
	"fmt"
	"go_event_driven/product/configurations"

	_ "github.com/go-sql-driver/mysql" // Register driver
)

type MySqlDatabaseAdapter struct {
	Database *sql.DB
}

const driver = "mysql"

func NewMySqlDatabaseAdapter(databaseConfiguration *configurations.DatabaseConfiguration) (*MySqlDatabaseAdapter, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		databaseConfiguration.UserName,
		databaseConfiguration.Password,
		databaseConfiguration.Host,
		databaseConfiguration.Port,
		databaseConfiguration.Name,
	)

	_database, _error := sql.Open(driver, connectionString)

	if _error != nil {
		return nil, _error
	}

	return &MySqlDatabaseAdapter{
		Database: _database,
	}, nil
}
