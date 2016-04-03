package schema

import (
	"strings"

	"github.com/go-zero-boilerplate/databases"
)

type SchemaDialect interface {
	databases.Dialect
	Accept(DialectVisitor)
	IndexNameFromFieldNames(fieldNames ...string) (string, error)
}

var (
	MysqlSchemaDialect    = NewMysqlSchemaDialect()
	SqliteSchemaDialect   = NewSqliteSchemaDialect()
	PostgresSchemaDialect = NewPostgresSchemaDialect()
)

func ParseSchemaDialectFromString(dialectStr string) SchemaDialect {
	switch strings.ToLower(dialectStr) {
	case "mysql":
		return MysqlSchemaDialect
	case "sqlite":
		return SqliteSchemaDialect
	case "postgres":
		return PostgresSchemaDialect
	default:
		panic("Schema dialect '" + dialectStr + "' not supported")
	}
}
