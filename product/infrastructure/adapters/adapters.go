package adapters

import (
	"go_event_driven/product/configurations"
	"go_event_driven/product/infrastructure/adapters/kafka"
	"go_event_driven/product/infrastructure/adapters/mysql"
	"go_event_driven/product/infrastructure/adapters/zap"
)

func SetupKafkaAdapter(configuration *configurations.KafkaConfiguration) *kafka.KafkaEventBusAdapter {
	adapter, _error := kafka.NewKafkaEventBusAdapter(configuration)

	if _error != nil {
		// TODO
	}

	return adapter
}

func SetupMySqlAdapter(databaseConfiguration *configurations.DatabaseConfiguration) *mysql.MySqlDatabaseAdapter {
	return mysql.NewMySqlDatabaseAdapter(databaseConfiguration)
}

func SetupProductMysqlAdapter(adapter *mysql.MySqlDatabaseAdapter) *mysql.MySqlProductDatabaseAdapter {
	return mysql.NewMySqlProductDatabaseAdapter(adapter)
}

func SetupCriterionBuilderMysqlAdapter() *mysql.SqlCriterionBuilder {
	return mysql.NewSqlCriterionBuilder()
}

func SetupLoggerZapAdapter() *zap.ZapLogger {
	return zap.NewZapLogger()
}
