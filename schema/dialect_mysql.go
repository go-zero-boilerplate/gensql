package schema

import (
	"strings"

	"github.com/go-zero-boilerplate/databases"

	"fmt"
)

type mysql struct {
	databases.Dialect
}

func NewMysqlSchemaDialect() SchemaDialect {
	d := &mysql{}
	d.Dialect = databases.MysqlDialect
	return d
}

func (m *mysql) Token(keyword SchemaKeyword) (_ string) {
	visitor := &mysqlTokenVisitor{}
	keyword.Accept(visitor)
	return visitor.keyword
}

func (m *mysql) Column(f *Field) string {
	size := f.Size
	if size == 0 {
		size = DEFAULT_VARCHAR_SIZE
	}

	visitor := &mysqlFieldTypeVisitor{fieldSize: size}
	f.Type.Accept(visitor)
	return visitor.typKeyword
}

func (m *mysql) WrapDefaultValue(defaultValue string) string {
	if strings.EqualFold(defaultValue, "null") {
		return "NULL"
	}
	return "'" + defaultValue + "'"
}

type mysqlTokenVisitor struct {
	keyword string
}

func (m *mysqlTokenVisitor) VisitAutoIncrement(*AutoIncrementKeyword) { m.keyword = "AUTO_INCREMENT" }
func (m *mysqlTokenVisitor) VisitPrimaryKey(*PrimaryKeyKeyword)       { m.keyword = "PRIMARY KEY" }
func (m *mysqlTokenVisitor) VisitIndex(*IndexKeyword)                 { m.keyword = "INDEX" }
func (m *mysqlTokenVisitor) VisitUniqueIndex(*UniqueIndexKeyword)     { m.keyword = "UNIQUE INDEX" }
func (m *mysqlTokenVisitor) VisitNull(*NullKeyword)                   { m.keyword = "NULL" }
func (m *mysqlTokenVisitor) VisitNotNull(*NotNullKeyword)             { m.keyword = "NOT NULL" }
func (m *mysqlTokenVisitor) VisitDefault(*DefaultKeyword)             { m.keyword = "DEFAULT" }

type mysqlFieldTypeVisitor struct {
	fieldSize int

	typKeyword string
}

func (m *mysqlFieldTypeVisitor) VisitInteger(*IntegerFieldType) { m.typKeyword = "INTEGER" }
func (m *mysqlFieldTypeVisitor) VisitVarchar(*VarcharFieldType) {
	m.typKeyword = fmt.Sprintf("VARCHAR(%d)", m.fieldSize)
}
func (m *mysqlFieldTypeVisitor) VisitBoolean(*BooleanFieldType)     { m.typKeyword = "BOOLEAN" }
func (m *mysqlFieldTypeVisitor) VisitReal(*RealFieldType)           { m.typKeyword = "DOUBLE" }
func (m *mysqlFieldTypeVisitor) VisitBlob(*BlobFieldType)           { m.typKeyword = "MEDIUMBLOB" }
func (m *mysqlFieldTypeVisitor) VisitDateTime(*DateTimeFieldType)   { m.typKeyword = "DATETIME" }
func (m *mysqlFieldTypeVisitor) VisitTimeStamp(*TimeStampFieldType) { m.typKeyword = "TIMESTAMP" }
