package adapters

import (
	"go_event_driven/product/configurations"
	"go_event_driven/product/domain/ports"
	"go_event_driven/product/infrastructure/adapters/kafka"
	"go_event_driven/product/infrastructure/adapters/mysql"
	"go_event_driven/product/infrastructure/adapters/zap"
)

func SetupKafkaAdapter(configuration *configurations.KafkaConfiguration, logger ports.Logger) (*kafka.KafkaEventBusAdapter, error) {
	return kafka.NewKafkaEventBusAdapter(configuration, logger)
}

func SetupMySqlAdapter(databaseConfiguration *configurations.DatabaseConfiguration) (*mysql.MySqlDatabaseAdapter, error) {
	return mysql.NewMySqlDatabaseAdapter(databaseConfiguration)
}

func SetupProductMysqlAdapter(adapter *mysql.MySqlDatabaseAdapter, logger ports.Logger) *mysql.MySqlProductDatabaseAdapter {
	return mysql.NewMySqlProductDatabaseAdapter(adapter, logger)
}

func SetupCriterionBuilderMysqlAdapter() *mysql.SqlCriterionBuilder {
	return mysql.NewSqlCriterionBuilder()
}

func SetupLoggerZapAdapter() *zap.ZapLogger {
	return zap.NewZapLogger()
}
