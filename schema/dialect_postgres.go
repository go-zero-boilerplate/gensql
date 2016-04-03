package schema

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-zero-boilerplate/databases"
)

type postgres struct {
	databases.Dialect
}

func NewPostgresSchemaDialect() SchemaDialect {
	d := &postgres{}
	d.Dialect = databases.PostgresDialect
	return d
}

func (p *postgres) Token(keyword SchemaKeyword) (_ string) {
	visitor := &postgresTokenVisitor{}
	keyword.Accept(visitor)
	return visitor.keyword
}

func (p *postgres) Column(f *Field) string {
	if f.Auto {
		return "SERIAL"
	}

	visitor := &postgresFieldTypeVisitor{fieldSize: f.Size}
	f.Type.Accept(visitor)
	return visitor.typKeyword
}

func (p *postgres) WrapDefaultValue(defaultValue string) string {
	if strings.EqualFold(defaultValue, "null") {
		return "NULL"
	}
	if _, err := strconv.ParseFloat(defaultValue, 32); err == nil {
		return "" + defaultValue + ""
	}
	return "'" + defaultValue + "'"
}

func (p *postgres) IndexNameFromFieldNames(fieldNames ...string) (string, error) {
	indexName := strings.Join(fieldNames, "_")
	if len(indexName) > 64 {
		return "", fmt.Errorf("The combined field names exceed 64 characters")
	}
	return indexName, nil
}

type postgresTokenVisitor struct {
	keyword string
}

func (p *postgresTokenVisitor) VisitAutoIncrement(*AutoIncrementKeyword) {} // postgres does not have this keyword but uses SERIAL field type
func (p *postgresTokenVisitor) VisitPrimaryKey(*PrimaryKeyKeyword)       { p.keyword = "PRIMARY KEY" }
func (p *postgresTokenVisitor) VisitIndex(*IndexKeyword)                 { p.keyword = "INDEX" }
func (p *postgresTokenVisitor) VisitUniqueIndex(*UniqueIndexKeyword)     { p.keyword = "UNIQUE INDEX" }
func (p *postgresTokenVisitor) VisitNull(*NullKeyword)                   { p.keyword = "NULL" }
func (p *postgresTokenVisitor) VisitNotNull(*NotNullKeyword)             { p.keyword = "NOT NULL" }
func (p *postgresTokenVisitor) VisitDefault(*DefaultKeyword)             { p.keyword = "DEFAULT" }

type postgresFieldTypeVisitor struct {
	fieldSize int

	typKeyword string
}

func (p *postgresFieldTypeVisitor) VisitInteger(*IntegerFieldType) { p.typKeyword = "INTEGER" }
func (p *postgresFieldTypeVisitor) VisitVarchar(*VarcharFieldType) {
	if p.fieldSize > 0 {
		p.typKeyword = fmt.Sprintf("VARCHAR(%d)", p.fieldSize)
	} else {
		p.typKeyword = "TEXT"
	}
}
func (p *postgresFieldTypeVisitor) VisitBoolean(*BooleanFieldType) { p.typKeyword = "BOOLEAN" }
func (p *postgresFieldTypeVisitor) VisitReal(*RealFieldType)       { p.typKeyword = "REAL" }
func (p *postgresFieldTypeVisitor) VisitBlob(*BlobFieldType)       { p.typKeyword = "BYTEA" }
func (p *postgresFieldTypeVisitor) VisitDateTime(*DateTimeFieldType) {
	p.typKeyword = "timestamp without time zone"
}
func (p *postgresFieldTypeVisitor) VisitTimeStamp(*TimeStampFieldType) {
	p.typKeyword = "timestamp without time zone"
}
