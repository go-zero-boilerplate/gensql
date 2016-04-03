package schema

import (
	"fmt"
)

// Some inspiration from: https://github.com/drone/sqlgen

type TableBuilder interface {
	Field(name string, fieldType FieldType, primary, auto bool, size int) TableBuilder
	FieldWithDefault(name string, fieldType FieldType, primary, auto bool, size int, defaultValue string) TableBuilder
	NullableField(name string, fieldType FieldType, size int) TableBuilder
	Primary(fieldNames ...string) TableBuilder
	Index(name string, unique bool, fieldNames ...string) TableBuilder

	Build() *Table
}

func NewTableBuilder(tableName string) TableBuilder {
	return &tableBuilder{
		t: &Table{
			Name: tableName,
		},
	}
}

type tableBuilder struct {
	t *Table
}

func (t *tableBuilder) getFieldsFromFieldNames(fieldNames []string) (fields []*Field, err error) {
	for _, fieldName := range fieldNames {
		foundField := false
		for _, f := range t.t.Fields {
			if f.Name == fieldName {
				fields = append(fields, f)
				foundField = true
				break
			}
		}

		if !foundField {
			fields = nil
			err = fmt.Errorf("Field '%s' is not added to the table yet", fieldName)
			return
		}
	}
	err = nil
	return
}

func (t *tableBuilder) Field(name string, fieldType FieldType, primary, auto bool, size int) TableBuilder {
	t.t.Fields = append(t.t.Fields, &Field{Name: name, Type: fieldType, Primary: primary, Auto: auto, Size: size})
	return t
}

func (t *tableBuilder) FieldWithDefault(name string, fieldType FieldType, primary, auto bool, size int, defaultValue string) TableBuilder {
	t.t.Fields = append(t.t.Fields, &Field{Name: name, Type: fieldType, Primary: primary, Auto: auto, Size: size, Default: defaultValue})
	return t
}

func (t *tableBuilder) NullableField(name string, fieldType FieldType, size int) TableBuilder {
	t.t.Fields = append(t.t.Fields, &Field{Name: name, Type: fieldType, Nullable: true, Size: size})
	return t
}

func (t *tableBuilder) Primary(fieldNames ...string) TableBuilder {
	if len(fieldNames) == 0 {
		panic("Primary requires at least one field")
	}

	fields, err := t.getFieldsFromFieldNames(fieldNames)
	if err != nil {
		panic(fmt.Errorf("Cannot set primary key, error: %s", err.Error()))
	}

	t.t.Primary = fields
	return t
}

func (t *tableBuilder) Index(name string, unique bool, fieldNames ...string) TableBuilder {
	if len(fieldNames) == 0 {
		panic("Index requires at least one field")
	}

	fields, err := t.getFieldsFromFieldNames(fieldNames)
	if err != nil {
		panic(fmt.Errorf("Cannot add index '%s', error: %s", name, err.Error()))
	}

	t.t.Indexes = append(t.t.Indexes, &Index{Name: name, Unique: unique, Fields: fields})
	return t
}

func (t *tableBuilder) Build() *Table {
	return t.t
}
