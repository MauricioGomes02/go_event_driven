package mysql

import (
	"fmt"
	"go_event_driven/product/domain/ports"
	"strings"
)

type SqlCriterion struct {
	query     string
	arguments []any
}

func (sql *SqlCriterion) ToSql() (string, []any) {
	return sql.query, sql.arguments
}

func (sql *SqlCriterion) And(others ...ports.CriterionPort) ports.CriterionPort {
	var parts []string
	var arguments []any

	parts = append(parts, "("+sql.query+")")
	arguments = append(arguments, sql.arguments...)

	for _, criterion := range others {
		sqlCriterion := criterion.(*SqlCriterion)
		parts = append(parts, "("+sqlCriterion.query+")")
		arguments = append(arguments, sqlCriterion.arguments...)
	}

	return &SqlCriterion{
		query:     strings.Join(parts, " AND "),
		arguments: arguments,
	}
}

func (sql *SqlCriterion) Or(others ...ports.CriterionPort) ports.CriterionPort {
	var parts []string
	var arguments []any

	parts = append(parts, "("+sql.query+")")
	arguments = append(arguments, sql.arguments...)

	for _, criterion := range others {
		sqlCriterion := criterion.(*SqlCriterion)
		parts = append(parts, "("+sqlCriterion.query+")")
		arguments = append(arguments, sqlCriterion.arguments...)
	}

	return &SqlCriterion{
		query:     strings.Join(parts, " OR "),
		arguments: arguments,
	}
}

type SqlCriterionBuilder struct{}

func NewSqlCriterionBuilder() *SqlCriterionBuilder {
	return &SqlCriterionBuilder{}
}

var (
	FieldMapping = map[string]map[string]string{
		"OutboxEvent": {
			"Status":       "status",
			"ErrorMessage": "error_message",
			"SentAt":       "sent_at",
			"Retries":      "retries",
		},
	}
)

func (builder *SqlCriterionBuilder) Where(entity string, property string, operator string, value any) ports.CriterionPort {
	field, _ := FieldMapping[entity][property]
	return &SqlCriterion{
		query:     fmt.Sprintf("%s %s ?", field, operator),
		arguments: []any{value},
	}
}

func (builder *SqlCriterionBuilder) And(conditions ...ports.CriterionPort) ports.CriterionPort {
	if len(conditions) == 0 {
		return &SqlCriterion{query: "1=1"} // neutro
	}
	criterion := conditions[0].(*SqlCriterion)
	return criterion.And(conditions[1:]...)
}

func (builder SqlCriterionBuilder) Or(conditions ...ports.CriterionPort) ports.CriterionPort {
	if len(conditions) == 0 {
		return &SqlCriterion{query: "1=0"} // neutro
	}
	criterion := conditions[0].(*SqlCriterion)
	return criterion.Or(conditions[1:]...)
}
