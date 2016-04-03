package schema

import (
	"strings"
)

func GenerateSchema(schemaDialect SchemaDialect, t *Table) string {
	lines := []string{}

	lines = append(lines, GenerateSchemaCreateTable(schemaDialect, t))

	for _, ix := range t.Indexes {
		lines = append(lines, GenerateSchemaCreateIndex(schemaDialect, t, ix))
	}

	return strings.Join(lines, "\n")
}
