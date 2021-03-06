package main

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/go-zero-boilerplate/gensql/schema"
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

	for indexGroupNum, indexGroup := range entity.Indexes {
		fieldNames := []string{}
		for _, ig := range indexGroup {
			fieldNames = append(fieldNames, ig.SqlColumn)
		}
		nameOfIndex, err := entity.Dialect.Dialect.IndexNameFromFieldNames(fieldNames...)
		if err != nil {
			nameOfIndex = fmt.Sprintf("entity_index_%d", indexGroupNum)
		}
		builder = builder.Index(nameOfIndex, false, fieldNames...)
	}

	for indexGroupNum, uniqueGroup := range entity.Uniques {
		fieldNames := []string{}
		for _, uf := range uniqueGroup {
			fieldNames = append(fieldNames, uf.SqlColumn)
		}
		nameOfUniqueIndex, err := entity.Dialect.Dialect.IndexNameFromFieldNames(fieldNames...)
		if err != nil {
			nameOfUniqueIndex = fmt.Sprintf("entity_unique_%d", indexGroupNum)
		}
		builder = builder.Index(nameOfUniqueIndex, true, fieldNames...)
	}

	table := builder.Build()

	schemaText := schema.GenerateSchema(entity.Dialect.Dialect, table)
	a.appendTemplate(
		`
			import (
				"github.com/go-zero-boilerplate/databases"
			)

			func {{.Entity.Dialect.Name | CamelFirstUpper}}CreateSchema_{{.Entity.EntityName | CamelFirstUpper}}(db databases.Database) error {
				_, err := db.Exec(`+"`{{.SchemaText}}`"+`)
				return err
			}
			`,
		map[string]interface{}{
			"Entity":     entity,
			"SchemaText": schemaText,
		},
	)
	return a
}

func (a *Appender) AppendEntityStructs(entity *GeneratorEntity) *Appender {
	imports := []string{}

	for _, importPath := range entity.AdditionalImports.EntityStruct {
		imports = append(imports, importPath)
	}

	a.appendLines(fmt.Sprintf(`import (
		%s
	)`, JoinImportsForGoFile(imports, "\n")))

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

func (a *Appender) AppendEntityHelpers(entity *GeneratorEntity) *Appender {
	a.appendTemplate(
		`
			{{$entity := .Entity}}
			{{$entityChar := .Entity.VariableNameFirstLetter}}
			{{$sqlPlaceholderIndex0 := .SqlPlaceholderIndex0}}

			var {{.Entity.StructName}}Helpers = struct {
				SqlColumnNames *{{.Entity.VariableName}}SqlColumnNames
			}{
				SqlColumnNames: &{{.Entity.VariableName}}SqlColumnNames{
					{{range .Entity.AllFields}}
						{{- .FieldName}}: "{{.SqlColumn}}",
					{{end}}
				},
			}

			type {{.Entity.VariableName}}SqlColumnNames struct {
				{{range .Entity.AllFields}}
					{{- .FieldName}} string
				{{end}}
			}
			`,
		map[string]interface{}{
			"Entity":               entity,
			"SqlPlaceholderIndex0": entity.Dialect.Dialect.Placeholder(0),
		},
	)

	return a
}

func (a *Appender) AppendEntityIterators(entity *GeneratorEntity) *Appender {
	imports := []string{
		`github.com/go-zero-boilerplate/databases/sql_statement`,
		`github.com/go-zero-boilerplate/databases/paginator`,
	}

	for _, importPath := range entity.AdditionalImports.EntityIterator {
		imports = append(imports, importPath)
	}

	a.appendLines(fmt.Sprintf(`import (
		%s
	)`, JoinImportsForGoFile(imports, "\n")))

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
				{{.Entity.VariableName}}s []*{{.Entity.StructName}}
			}

			type db{{.Entity.VariableName}} struct {
				{{range .Entity.AllFields}}
					{{- .GoStructDbDef}} `+"`db:\"{{.SqlColumn}}\"`"+`
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
					{{range .Entity.AllFields -}}
						{{- if .IsNullable}}
							{{if .MustCastDbField}} {{.VariableName}} := {{.GoType}}(t.{{- .FieldName}}.Ptr())
							{{- else}} {{.VariableName}} := t.{{- .FieldName}}.Ptr()
							{{- end}}
						{{- end}}
					{{end -}}

					d.{{.Entity.VariableName}}s = append(d.{{.Entity.VariableName}}s, &{{.Entity.StructName}}{
						{{range .Entity.AllFields}}
							{{- if .IsNullable}}
								{{- .FieldName}}: {{- .VariableName}},
							{{- else}}
								{{- if .MustCastDbField}} {{- .FieldName}}: {{- .GoType}}(t.{{- .FieldName}}{{- .DbFieldNameDotSuffix}}),
								{{- else}} {{- .FieldName}}: t.{{- .FieldName}}{{- .DbFieldNameDotSuffix}},
								{{- end}}
							{{- end}}
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
	imports := []string{
		"github.com/go-zero-boilerplate/databases",
		"github.com/go-zero-boilerplate/databases/sql_statement",
	}

	if len(entity.NullableFields) > 0 {
		imports = append(imports, "gopkg.in/guregu/null.v3/zero")
	}

	a.appendLines(fmt.Sprintf(`import (
		%s
	)`, JoinImportsForGoFile(imports, "\n")))

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

func (a *Appender) AppendStatementBuilderFactory(entity *GeneratorEntity) *Appender {
	a.appendLines(`import (
		"github.com/go-zero-boilerplate/databases"
		"github.com/go-zero-boilerplate/databases/sql_statement"
	)`)

	a.appendTemplate(
		`
			{{$entityChar := .Entity.VariableNameFirstLetter}}

			func New{{.Entity.StructName}}StatementBuilderFactory() {{.Entity.StructName}}StatementBuilderFactory {
			    return &{{.Entity.VariableName}}StatementBuilderFactory{
			        dialect:   {{.Entity.Dialect.GoVariablePart}},
			        tableName: "{{.Entity.SqlTable}}",
			    }
			}

			type {{.Entity.StructName}}StatementBuilderFactory interface {
			    Insert(db databases.Database) sql_statement.InsertBuilder
			    Select(db databases.Database) sql_statement.SelectBuilder
			    Update(db databases.Database) sql_statement.UpdateBuilder
			    Delete(db databases.Database) sql_statement.DeleteBuilder
			}

			type {{.Entity.VariableName}}StatementBuilderFactory struct {
			    dialect   databases.Dialect
			    tableName string
			}

			func ({{$entityChar}} *{{.Entity.VariableName}}StatementBuilderFactory) Insert(db databases.Database) sql_statement.InsertBuilder {
			    return sql_statement.NewInsertBuilderFactory(db).FromDialect({{$entityChar}}.dialect, {{$entityChar}}.tableName)
			}

			func ({{$entityChar}} *{{.Entity.VariableName}}StatementBuilderFactory) Select(db databases.Database) sql_statement.SelectBuilder {
			    return sql_statement.NewSelectBuilderFactory(db).FromDialect({{$entityChar}}.dialect, {{$entityChar}}.tableName)
			}

			func ({{$entityChar}} *{{.Entity.VariableName}}StatementBuilderFactory) Update(db databases.Database) sql_statement.UpdateBuilder {
			    return sql_statement.NewUpdateBuilderFactory(db).FromDialect({{$entityChar}}.dialect, {{$entityChar}}.tableName)
			}

			func ({{$entityChar}} *{{.Entity.VariableName}}StatementBuilderFactory) Delete(db databases.Database) sql_statement.DeleteBuilder {
			    return sql_statement.NewDeleteBuilderFactory(db).FromDialect({{$entityChar}}.dialect, {{$entityChar}}.tableName)
			}
		`,
		map[string]interface{}{
			"Entity": entity,
		},
	)

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

				{{range .Entity.NullableFields -}}
				db{{- .FieldName}} := {{.DbGoType}}{}
				{{end -}}

				err := r.db.QueryRow("SELECT {{.Entity.NonPkFields | AsSqlSelectColumns}} FROM " + r.tableName + " WHERE {{.Entity.PkFields | AsSqlParameterizedWhereColumns}}",
					{{range .Entity.PkFields -}}
						{{- .VariableName -}},
					{{end -}}
				).
					Scan(
						{{range .Entity.NonPkFields -}}
							{{- if .IsNullable}} &db{{- .FieldName}},
							{{- else}} &{{$outerScope.Entity.VariableName}}.{{- .FieldName -}},
							{{- end}}
						{{end -}}
					)
				if err != nil {
					return nil, err
				}
				
				{{range .Entity.NullableFields -}}
					{{- if .MustCastDbField}} {{$outerScope.Entity.VariableName}}.{{- .FieldName -}} = {{- .GoType}}(db{{- .FieldName}}.Ptr())
					{{- else}} {{$outerScope.Entity.VariableName}}.{{- .FieldName -}} = db{{- .FieldName}}.Ptr()
					{{- end}}
				{{end -}}
				
				return {{.Entity.VariableName}}, err
			}

			func (r *{{.Entity.VariableName}}Repository) List() ({{.Entity.StructName}}Iterator, error) {
				selectBuilder := sql_statement.NewSelectBuilder(r.dialect, r.db, r.tableName)
				pageSize := 100
				return New{{.Entity.StructName}}Iterator(selectBuilder, pageSize)
			}

			func (r *{{.Entity.VariableName}}Repository) Add({{.Entity.VariableName}} *{{.Entity.StructName}}) error {
				var lastInsertId int64
				err := sql_statement.NewInsertBuilder(r.dialect, r.db, r.tableName).
					{{range .Entity.InsertableFields -}}
						Set("{{- .SqlColumn -}}", {{$outerScope.Entity.VariableName}}.{{- .FieldName}}).
					{{end -}}
					{{if .Entity.HasSingleIntPkField -}}
						LastInsertIdDest(&lastInsertId).
					{{end -}}
					Build().
					Execute()
				if err == nil {
					{{$outerScope.Entity.VariableName}}.{{.Entity.IntPkField.FieldName}} = {{.Entity.IntPkField.GoType}}(lastInsertId)
				}
				return err
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
					{{if .Entity.MustSetUpdated}}
						Set("{{- .Entity.UpdatedField.SqlColumn -}}", time.Now().UTC()).
					{{end}}
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
	a.appendLines(`import (
		"github.com/go-zero-boilerplate/databases"
	)`)

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

func (a *Appender) AppendStatementBuilderFactories(generatorSetup *GeneratorSetup) *Appender {
	a.appendTemplate(
		`
			func NewStatementBuilderFactories() StatementBuilderFactories {
				return &statementBuilderFactories{}
			}

			type StatementBuilderFactories interface {
				{{range .Entities}}
				{{.StructName}}() {{.StructName}}StatementBuilderFactory
				{{- end}}
			}

			type statementBuilderFactories struct{}

			{{range .Entities}}
			func (s *statementBuilderFactories) {{.StructName}}() {{.StructName}}StatementBuilderFactory {
				return New{{.StructName}}StatementBuilderFactory()
			}
			{{- end}}
		`,
		map[string]interface{}{
			"Entities": generatorSetup.Entities,
		},
	)

	return a
}

func (a *Appender) AsGoFile(packageName string) []byte {
	if a.lines == nil {
		return []byte{}
	}

	header := strings.Join([]string{
		"/*",
		"  DO NOT CHANGE THIS FILE as it was auto-generated by https://github.com/go-zero-boilerplate/gensql",
		"  VERSION " + Version,
		// "  ON " + time.Now().Format("2006-01-02 15-04-05"),
		"*/",
	}, "\n")

	combinedLines := fmt.Sprintf("%s\n\npackage %s\n\n%s", header, packageName, strings.Join(a.lines, "\n"))

	prettyCombined, err := format.Source([]byte(combinedLines))
	if err != nil {
		panic(fmt.Errorf("Cannot format, error: %s. Source was:\n%s", err.Error(), string(combinedLines)))
	}

	return prettyCombined
	/*
		TODO: Disabling this because imports are processed in the context of the environment executing it
		For example we import package `gopkg.in/guregu/null.v3/zero` for nullable fields. If that package is not
		installed but rather `github.com/guregu/null` it will get swapped for that

		processedImports, err := imports.Process(packageName, prettyCombined, nil)
		if err != nil {
			panic(fmt.Errorf("Cannot format with imports, error: %s", err.Error()))
		}

		return processedImports*/
}
