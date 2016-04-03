package schema

import "strings"

func GenerateSchemaCreateTable(schemaDialect SchemaDialect, table *Table) string {
	visitor := &dialectSchemaCreateTableVisitor{table: table}
	schemaDialect.Accept(visitor)
	return visitor.result
}

type dialectSchemaCreateTableVisitor struct {
	table *Table

	result string
}

func (d *dialectSchemaCreateTableVisitor) VisitMysql(m *mysql) {
	tabWriter := newTabWriter("\t", 2).AppendLine("")
	tabWriter.AppendLine("CREATE TABLE IF NOT EXISTS %s (", d.table.Name)

	tabWriter.LevelUp()

	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		conditionalAppender.Append(m.Column(f))
		conditionalAppender.AppendWithCondition(f.Primary, m.Token(PRIMARY_KEY))
		conditionalAppender.AppendWithCondition(f.Auto, m.Token(AUTO_INCREMENT))
		conditionalAppender.AppendWithCondition(!f.Nullable, m.Token(NOT_NULL))
		conditionalAppender.AppendWithCondition(f.Default != "", m.Token(DEFAULT)+" "+m.WrapDefaultValue(f.Default))

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

func (d *dialectSchemaCreateTableVisitor) VisitSqlite(s *sqlite) {
	tabWriter := newTabWriter("\t", 2).AppendLine("")
	tabWriter.AppendLine("CREATE TABLE IF NOT EXISTS %s (", d.table.Name)

	tabWriter.LevelUp()

	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		conditionalAppender.Append(s.Column(f))
		conditionalAppender.AppendWithCondition(f.Primary, s.Token(PRIMARY_KEY))
		conditionalAppender.AppendWithCondition(f.Auto, s.Token(AUTO_INCREMENT))
		conditionalAppender.AppendWithCondition(!f.Nullable, s.Token(NOT_NULL))
		conditionalAppender.AppendWithCondition(f.Default != "", s.Token(DEFAULT)+" "+s.WrapDefaultValue(f.Default))

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

func (d *dialectSchemaCreateTableVisitor) VisitPostgres(p *postgres) {
	tabWriter := newTabWriter("\t", 2).AppendLine("")
	tabWriter.AppendLine("CREATE TABLE IF NOT EXISTS %s (", d.table.Name)

	tabWriter.LevelUp()

	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		conditionalAppender.Append(p.Column(f))
		conditionalAppender.AppendWithCondition(f.Primary, p.Token(PRIMARY_KEY))
		conditionalAppender.AppendWithCondition(f.Auto, p.Token(AUTO_INCREMENT))
		conditionalAppender.AppendWithCondition(!f.Nullable, p.Token(NOT_NULL))
		conditionalAppender.AppendWithCondition(f.Default != "", p.Token(DEFAULT)+" "+p.WrapDefaultValue(f.Default))

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
