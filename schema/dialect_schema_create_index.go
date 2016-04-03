package schema

import "strings"

func GenerateSchemaCreateIndex(schemaDialect SchemaDialect, table *Table, index *Index) string {
	visitor := &dialectSchemaCreateIndexVisitor{table: table, index: index}
	schemaDialect.Accept(visitor)
	return visitor.result
}

type dialectSchemaCreateIndexVisitor struct {
	table *Table
	index *Index

	result string
}

func (d *dialectSchemaCreateIndexVisitor) VisitMysql(m *mysql) {
	var obj = m.Token(INDEX)
	if d.index.Unique {
		obj = m.Token(UNIQUE_INDEX)
	}

	tabWriter := newTabWriter("\t", 2)
	tabWriter.AppendLine("")
	tabWriter.AppendLine("CREATE %s %s ON %s (", obj, d.index.Name, d.table.Name)

	tabWriter.LevelUp()
	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		trimmedJoined := strings.TrimSpace(strings.Join(conditionalAppender.Slice(), " "))

		suffix := ""
		if i < len(d.table.Fields)-1 {
			suffix = ","
		}
		tabWriter.AppendLine(trimmedJoined + suffix)
	}
	tabWriter.LevelDown()

	tabWriter.AppendLine(");")

	d.result = tabWriter.CombineLines()
}

func (d *dialectSchemaCreateIndexVisitor) VisitSqlite(s *sqlite) {
	var obj = s.Token(INDEX)
	if d.index.Unique {
		obj = s.Token(UNIQUE_INDEX)
	}

	tabWriter := newTabWriter("\t", 2)
	tabWriter.AppendLine("")
	tabWriter.AppendLine("CREATE %s %s ON %s (", obj, d.index.Name, d.table.Name)

	tabWriter.LevelUp()
	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		trimmedJoined := strings.TrimSpace(strings.Join(conditionalAppender.Slice(), " "))

		suffix := ""
		if i < len(d.table.Fields)-1 {
			suffix = ","
		}
		tabWriter.AppendLine(trimmedJoined + suffix)
	}
	tabWriter.LevelDown()

	tabWriter.AppendLine(");")

	d.result = tabWriter.CombineLines()
}

func (d *dialectSchemaCreateIndexVisitor) VisitPostgres(p *postgres) {
	var obj = p.Token(INDEX)
	if d.index.Unique {
		obj = p.Token(UNIQUE_INDEX)
	}

	tabWriter := newTabWriter("\t", 2)
	tabWriter.AppendLine("")
	tabWriter.AppendLine("CREATE %s %s ON %s (", obj, d.index.Name, d.table.Name)

	tabWriter.LevelUp()
	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		trimmedJoined := strings.TrimSpace(strings.Join(conditionalAppender.Slice(), " "))

		suffix := ""
		if i < len(d.table.Fields)-1 {
			suffix = ","
		}
		tabWriter.AppendLine(trimmedJoined + suffix)
	}
	tabWriter.LevelDown()

	tabWriter.AppendLine(");")

	d.result = tabWriter.CombineLines()
}
