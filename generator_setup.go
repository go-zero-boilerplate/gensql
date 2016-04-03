package main

import (
	"fmt"
	"strconv"

	"github.com/go-zero-boilerplate/gensql/schema"

	"github.com/codemodus/kace"

	"strings"
)

type GeneratorSetup struct {
	Entities []*GeneratorEntity
}

type GeneratorEntity struct {
	Dialect      *GeneratorDialect
	EntityName   string
	StructName   string
	VariableName string
	SqlTable     string

	AllFields        []*GeneratorField
	PkFields         []*GeneratorField
	NonPkFields      []*GeneratorField
	InsertableFields []*GeneratorField
	EditableFields   []*GeneratorField
	TriggerFields    []*GeneratorField
}

type GeneratorField struct {
	IsAuto         bool
	IsPkField      bool
	IsTriggerField bool
	IsNullable     bool

	FieldName    string
	VariableName string
	GoType       string
	GoStructDef  string
	SqlColumn    string
	SchemaType   schema.FieldType
	SqlSize      int
	DefaultValue string
	ExtraArgs    []string
}

type GeneratorDialect struct {
	Name           string
	Dialect        schema.SchemaDialect
	GoVariablePart string
}

func splitStringRemoveEmpty(s, separator string) (strs []string) {
	for _, split := range strings.Split(s, separator) {
		trimmed := strings.TrimSpace(split)
		if trimmed != "" {
			strs = append(strs, trimmed)
		}
	}
	return
}

func schemaFieldTypeFromGoType(goType string, isCreatedField, isUpdatedField bool) schema.FieldType {
	if isCreatedField {
		return schema.DATETIME
	} else if isUpdatedField {
		return schema.TIMESTAMP
	}

	fieldType, err := schema.GoToFieldType(goType)
	if err != nil {
		panic(fmt.Errorf("Cannot append schema, unable to get field type for go type '%s', error: %s", goType, err.Error()))
	}
	return fieldType
}

func generatorFieldFromString(s string) *GeneratorField {
	splitted := splitStringRemoveEmpty(s, " ")
	if len(splitted) == 0 {
		panic("Cannot generate field from string '" + s + "'")
	}

	var (
		fieldName      string = splitted[0]
		fieldGoType    string
		fieldArgs      []string
		isTriggerField bool             = false
		isCreatedField bool             = false
		isUpdatedField bool             = false
		isAuto         bool             = false
		isPkField      bool             = false
		isNullable     bool             = false
		schemaType     schema.FieldType = schema.VARCHAR
		sqlSize        int              = 0
		defaultValue   string           = ""
	)

	if len(splitted) == 1 {
		fieldArgs = []string{}
		switch splitted[0] {
		case "created":
			fieldGoType = "time.Time"
			isTriggerField = true
			isCreatedField = true
		case "updated":
			fieldGoType = "time.Time"
			isTriggerField = true
			isUpdatedField = true
		default:
			panic("GoType not supported for 'trigger' field name: " + splitted[0])
		}
	} else {
		fieldGoType = splitted[1]
		fieldArgs = splitted[2:]
	}

	for _, arg := range fieldArgs {
		switch strings.ToLower(arg) {
		case "pk":
			isPkField = true
		case "auto":
			isAuto = true
		case "nullable":
			isNullable = true
		}

		if num, err := strconv.ParseInt(arg, 10, 32); err == nil {
			sqlSize = int(num)
		}

		if strings.HasPrefix(strings.ToLower(arg), "default:") {
			defaultValue = strings.Split(arg, ":")[1]
		}
	}

	schemaType = schemaFieldTypeFromGoType(fieldGoType, isCreatedField, isUpdatedField)

	return &GeneratorField{
		IsAuto:         isAuto,
		IsPkField:      isPkField,
		IsTriggerField: isTriggerField,
		IsNullable:     isNullable,
		FieldName:      kace.Camel(fieldName, true),
		VariableName:   kace.Camel(fieldName, false),
		GoType:         fieldGoType,
		GoStructDef:    kace.Camel(fieldName, true) + " " + fieldGoType,
		SqlColumn:      kace.Snake(fieldName),
		SchemaType:     schemaType,
		SqlSize:        sqlSize,
		DefaultValue:   defaultValue,
		ExtraArgs:      fieldArgs,
	}
}

func getSchemaGoVariablePart(dialectString string) string {
	switch strings.ToLower(dialectString) {
	case "mysql":
		return "database.MysqlDialect"
	case "sqlite":
		return "database.SqliteDialect"
	case "postgres":
		return "database.PostgresDialect"
	default:
		panic("Dialect '" + dialectString + "' not supported for the generation 'variable' part")
	}
}

func GeneratorSetupFromYamlSetup(orderedEntityNames []string, y *YamlSetup) (g *GeneratorSetup) {
	g = &GeneratorSetup{}

	for _, entityName := range orderedEntityNames {
		entitySetup := (*y)[entityName]

		allFields := []*GeneratorField{}
		pkFields := []*GeneratorField{}
		nonPkFields := []*GeneratorField{}
		editableFields := []*GeneratorField{}
		insertableFields := []*GeneratorField{}
		triggerFields := []*GeneratorField{}

		for _, fieldString := range entitySetup.Fields {
			field := generatorFieldFromString(fieldString)
			allFields = append(allFields, field)

			if field.IsPkField {
				pkFields = append(pkFields, field)
			} else {
				nonPkFields = append(nonPkFields, field)
			}

			if field.IsTriggerField {
				triggerFields = append(triggerFields, field)
			}

			if !field.IsPkField && !field.IsTriggerField {
				editableFields = append(editableFields, field)
				insertableFields = append(insertableFields, field)
			}
		}

		dialectName := kace.Camel(entitySetup.Dialect, true)
		generatorDialect := &GeneratorDialect{
			Name:           dialectName,
			GoVariablePart: getSchemaGoVariablePart(entitySetup.Dialect),
			Dialect:        schema.ParseSchemaDialectFromString(entitySetup.Dialect),
		}

		g.Entities = append(g.Entities, &GeneratorEntity{
			Dialect: generatorDialect,

			EntityName:   entityName,
			StructName:   kace.Camel(entityName, true),
			VariableName: kace.Camel(entityName, false),
			SqlTable:     kace.Snake(entityName),

			AllFields:        allFields,
			PkFields:         pkFields,
			NonPkFields:      nonPkFields,
			EditableFields:   editableFields,
			InsertableFields: insertableFields,
			TriggerFields:    triggerFields,
		})
	}

	return g
}
