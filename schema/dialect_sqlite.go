package schema

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-zero-boilerplate/databases"
)

type sqlite struct {
	databases.Dialect
}

func NewSqliteSchemaDialect() SchemaDialect {
	d := &sqlite{}
	d.Dialect = databases.SqliteDialect
	return d
}

func (s *sqlite) Token(keyword SchemaKeyword) (_ string) {
	visitor := &sqliteTokenVisitor{}
	keyword.Accept(visitor)
	return visitor.keyword
}

func (s *sqlite) Column(f *Field) string {
	visitor := &sqliteFieldTypeVisitor{}
	f.Type.Accept(visitor)
	return visitor.typKeyword
}

func (s *sqlite) WrapDefaultValue(defaultValue string) string {
	if strings.EqualFold(defaultValue, "null") {
		return "NULL"
	}
	if _, err := strconv.ParseFloat(defaultValue, 32); err == nil {
		return "(" + defaultValue + ")"
	}
	return "('" + defaultValue + "')"
}

func (s *sqlite) IndexNameFromFieldNames(fieldNames ...string) (string, error) {
	indexName := strings.Join(fieldNames, "_")
	if len(indexName) > 64 {
		return "", fmt.Errorf("The combined field names exceed 64 characters")
	}
	return indexName, nil
}

type sqliteTokenVisitor struct {
	keyword string
}

func (s *sqliteTokenVisitor) VisitAutoIncrement(*AutoIncrementKeyword) { s.keyword = "AUTOINCREMENT" }
func (s *sqliteTokenVisitor) VisitPrimaryKey(*PrimaryKeyKeyword)       { s.keyword = "PRIMARY KEY" }
func (s *sqliteTokenVisitor) VisitIndex(*IndexKeyword)                 { s.keyword = "INDEX" }
func (s *sqliteTokenVisitor) VisitUniqueIndex(*UniqueIndexKeyword)     { s.keyword = "UNIQUE INDEX" }
func (s *sqliteTokenVisitor) VisitNull(*NullKeyword)                   { s.keyword = "NULL" }
func (s *sqliteTokenVisitor) VisitNotNull(*NotNullKeyword)             { s.keyword = "NOT NULL" }
func (s *sqliteTokenVisitor) VisitDefault(*DefaultKeyword)             { s.keyword = "DEFAULT" }

type sqliteFieldTypeVisitor struct {
	typKeyword string
}

func (s *sqliteFieldTypeVisitor) VisitInteger(*IntegerFieldType)     { s.typKeyword = "INTEGER" }
func (s *sqliteFieldTypeVisitor) VisitVarchar(*VarcharFieldType)     { s.typKeyword = "TEXT" }
func (s *sqliteFieldTypeVisitor) VisitBoolean(*BooleanFieldType)     { s.typKeyword = "INTEGER" }
func (s *sqliteFieldTypeVisitor) VisitReal(*RealFieldType)           { s.typKeyword = "REAL" }
func (s *sqliteFieldTypeVisitor) VisitBlob(*BlobFieldType)           { s.typKeyword = "BLOB" }
func (s *sqliteFieldTypeVisitor) VisitDateTime(*DateTimeFieldType)   { s.typKeyword = "DATETIME" }
func (s *sqliteFieldTypeVisitor) VisitTimeStamp(*TimeStampFieldType) { s.typKeyword = "DATETIME" }
