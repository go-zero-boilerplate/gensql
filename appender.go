package main

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/go-zero-boilerplate/gensql/schema"

	"golang.org/x/tools/imports"
)

func NewAppender() *Appender {
	return &Appender{}
}

type Appender struct {
	lines []string
}

func (a *Appender) appendLines(lines ...string) {
	a.lines = append(a.lines, lines...)
}

func (a *Appender) appendTemplate(templateString string, data interface{}) {
	tpl := mustExecTemplateToString(templateString, data)
	a.appendLines(tpl)
}

func (a *Appender) AppendSchemaCreate(entity *GeneratorEntity) *Appender {
	primaryNames := []string{}
	builder := schema.NewTableBuilder(entity.SqlTable)
	for _, field := range entity.AllFields {
		if field.IsNullable {
			builder = builder.NullableField(field.SqlColumn, field.SchemaType, field.SqlSize)
		} else if field.DefaultValue != "" {
			builder = builder.FieldWithDefault(field.SqlColumn, field.SchemaType, field.IsPkField, field.IsAuto, field.SqlSize, field.DefaultValue)
		} else {
			builder = builder.Field(field.SqlColumn, field.SchemaType, field.IsPkField, field.IsAuto, field.SqlSize)
		}
		if field.IsPkField {
			primaryNames = append(primaryNames, field.SqlColumn)
		}
	}
	builder = builder.Primary(primaryNames...)
	//builder.Index(name, unique, ...)

	table := builder.Build()

	schemaText := schema.GenerateSchema(entity.Dialect.Dialect, table)
	a.appendLines(fmt.Sprintf(`const ( 
			%s__CREATE_SQL_%s = `+"`%s`"+`
		)`,
		strings.ToUpper(entity.Dialect.Name),
		strings.ToUpper(entity.EntityName),
		schemaText,
	))
	return a
}

func (a *Appender) AppendEntityStructs(entity *GeneratorEntity) *Appender {
	a.appendTemplate(
		`
			type {{.Entity.StructName}} struct {
				{{range .Entity.AllFields}}
					{{- .GoStructDef}}
				{{end}}
			}
			`,
		map[string]interface{}{
			"Entity": entity,
		},
	)

	return a
}

func (a *Appender) AppendEntityIterators(entity *GeneratorEntity) *Appender {
	a.appendLines(`import (
		"github.com/go-zero-boilerplate/databases/sql_statement"
		"github.com/go-zero-boilerplate/databases/paginator"
	)`)

	a.appendTemplate(
		`
			type {{.Entity.StructName}}Iterator interface {
				HasMore() bool
				Next() (*{{.Entity.StructName}}, error)
			}

			func New{{.Entity.StructName}}Iterator(selectBuilder sql_statement.SelectBuilder, pageSize int) ({{.Entity.StructName}}Iterator, error) {
				iterator := &db{{.Entity.VariableName}}Iterator{}
				paginator, err := paginator.NewDBPaginator(selectBuilder, pageSize, iterator)
				if err != nil {
					return nil, err
				}
				iterator.paginator = paginator
				return iterator, nil
			}

			type db{{.Entity.VariableName}}Iterator struct {
				paginator paginator.DBPaginator

				tmpDestinationSlice []*db{{.Entity.VariableName}}
				{{.Entity.VariableName}}s               []*{{.Entity.StructName}}
			}

			type db{{.Entity.VariableName}} struct {
				{{range .Entity.AllFields}}
					{{- .GoStructDef}}
				{{end}}
			}

			func (d *db{{.Entity.VariableName}}Iterator) HasMore() bool {
				return d.paginator.HasMore()
			}

			func (d *db{{.Entity.VariableName}}Iterator) Next() (*{{.Entity.StructName}}, error) {
				nextIndex, err := d.paginator.GetNextIndex()
				if err != nil {
					return nil, err
				}
				return d.{{.Entity.VariableName}}s[nextIndex], nil
			}

			func (d *db{{.Entity.VariableName}}Iterator) Count() int {
				return len(d.{{.Entity.VariableName}}s)
			}

			func (d *db{{.Entity.VariableName}}Iterator) SlicePointer() interface{} {
				d.tmpDestinationSlice = []*db{{.Entity.VariableName}}{}
				return &d.tmpDestinationSlice
			}

			func (d *db{{.Entity.VariableName}}Iterator) Clear() {
				d.{{.Entity.VariableName}}s = nil
			}

			func (d *db{{.Entity.VariableName}}Iterator) AfterSliceLoaded() {
				for _, t := range d.tmpDestinationSlice {
					d.{{.Entity.VariableName}}s = append(d.{{.Entity.VariableName}}s, &{{.Entity.StructName}}{
						{{range .Entity.AllFields}}
							{{- .FieldName}}: t.{{- .FieldName}},
						{{end}}
					})
				}
				d.tmpDestinationSlice = nil
			}
			`,
		map[string]interface{}{
			"Entity": entity,
		},
	)

	return a
}

func (a *Appender) AppendRepoInterface(entity *GeneratorEntity) *Appender {
	a.appendTemplate(
		`
			type {{.Entity.StructName}}Repository interface {
				GetByPk({{- range .Entity.PkFields}} {{- .VariableName}} {{.GoType -}}, {{end -}}) (*{{.Entity.StructName}}, error)
				List() ({{.Entity.StructName}}Iterator, error)
				// ListFiltered(filter func(*{{.Entity.StructName}}) bool) ({{.Entity.StructName}}Iterator, error)  How to do this?
				Add({{.Entity.VariableName}} *{{.Entity.StructName}}) error
				Delete({{.Entity.VariableName}} *{{.Entity.StructName}}) error
				Save({{.Entity.VariableName}} *{{.Entity.StructName}}) error
			}
		`,
		map[string]interface{}{
			"Entity": entity,
		},
	)

	a.appendRepoDBImplementation(entity)

	return a
}

func (a *Appender) appendRepoDBImplementation(entity *GeneratorEntity) *Appender {
	//TODO: In the generated code in the `Where(` we use the ? symbol which is dialect specific

	a.appendTemplate(
		`
			{{$outerScope := .}}

			type {{.Entity.VariableName}}Repository struct {
				db        databases.Database
				dialect   databases.Dialect
				tableName string
			}

			func (r *{{.Entity.VariableName}}Repository) GetByPk({{- range .Entity.PkFields}} {{- .VariableName}} {{.GoType -}}, {{end -}}) (*{{.Entity.StructName}}, error) {
				{{.Entity.VariableName}} := &{{.Entity.StructName}}{
					{{range .Entity.PkFields -}}
						{{- .FieldName -}}: {{- .VariableName -}},
					{{end -}}
				}

				err := r.db.QueryRow("SELECT {{.Entity.NonPkFields | AsSqlSelectColumns}} FROM " + r.tableName + " WHERE {{.Entity.PkFields | AsSqlParameterizedWhereColumns}}",
					{{range .Entity.PkFields -}}
						{{- .VariableName -}},
					{{end -}}
				).
					Scan(
						{{range .Entity.NonPkFields -}}
							&{{$outerScope.Entity.VariableName}}.{{- .FieldName -}},
						{{end -}}
					)
				if err != nil {
					return nil, err
				}
				return {{.Entity.VariableName}}, err
			}

			func (r *{{.Entity.VariableName}}Repository) List() ({{.Entity.StructName}}Iterator, error) {
				selectBuilder := sql_statement.NewSelectBuilder(r.dialect, r.db, r.tableName)
				pageSize := 100
				return New{{.Entity.StructName}}Iterator(selectBuilder, pageSize)
			}

			func (r *{{.Entity.VariableName}}Repository) Add({{.Entity.VariableName}} *{{.Entity.StructName}}) error {
				return sql_statement.NewInsertBuilder(r.dialect, r.db, r.tableName).
					{{range .Entity.InsertableFields -}}
						Set("{{- .SqlColumn -}}", {{$outerScope.Entity.VariableName}}.{{- .FieldName}}).
					{{end -}}
					Build().
					Execute()
			}

			func (r *{{.Entity.VariableName}}Repository) Delete({{.Entity.VariableName}} *{{.Entity.StructName}}) error {
				return sql_statement.NewDeleteBuilder(r.dialect, r.db, r.tableName).
					{{range .Entity.PkFields -}}
						Where("{{- .VariableName -}} = ?", {{$outerScope.Entity.VariableName}}.{{.FieldName}}).
					{{end -}}
					Build().
					Execute()
			}

			func (r *{{.Entity.VariableName}}Repository) Save({{.Entity.VariableName}} *{{.Entity.StructName}}) error {
				return sql_statement.NewUpdateBuilder(r.dialect, r.db, r.tableName).
					{{range .Entity.InsertableFields -}}
						Set("{{- .SqlColumn -}}", {{$outerScope.Entity.VariableName}}.{{- .FieldName}}).
					{{end -}}
					{{range .Entity.PkFields -}}
						Where("{{- .VariableName -}} = ?", {{$outerScope.Entity.VariableName}}.{{.FieldName}}).
					{{end -}}
					Build().
					Execute()
			}
		`,
		map[string]interface{}{
			"Entity": entity,
		},
	)

	return a
}

func (a *Appender) AppendRepositoryFactories(generatorSetup *GeneratorSetup) *Appender {
	a.appendTemplate(
		`
			func NewRepositoryFactory() RepositoryFactory {
				return &repositoryFactory{}
			}

			type RepositoryFactory interface {
				{{range .Entities}}
				{{.StructName}}(db databases.Database) {{.StructName}}Repository
				{{- end}}
			}

			type repositoryFactory struct{}

			{{range .Entities}}
			func (r *repositoryFactory) {{.StructName}}(db databases.Database) {{.StructName}}Repository {
				return &{{.VariableName}}Repository{
					db: db,
					dialect: {{.Dialect.GoVariablePart}},
					tableName: "{{.SqlTable}}",
				}
			}
			{{end}}
		`,
		map[string]interface{}{
			"Entities": generatorSetup.Entities,
		},
	)

	return a
}

func (a *Appender) Bytes(packageName string) []byte {
	if a.lines == nil {
		return []byte{}
	}

	combinedLines := fmt.Sprintf("package %s\n\n%s", packageName, strings.Join(a.lines, "\n"))

	prettyCombined, err := format.Source([]byte(combinedLines))
	if err != nil {
		panic(fmt.Errorf("Cannot format, error: %s", err.Error()))
	}

	processedImports, err := imports.Process(packageName, prettyCombined, nil)
	if err != nil {
		panic(fmt.Errorf("Cannot format with imports, error: %s", err.Error()))
	}

	return processedImports
}
