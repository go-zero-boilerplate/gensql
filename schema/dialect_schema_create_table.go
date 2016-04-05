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
	tabWriter := newTabWriter("    ", 2).AppendLine("")
	tabWriter.AppendLine("CREATE TABLE IF NOT EXISTS %s (", d.table.Name)

	tabWriter.LevelUp()

	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		conditionalAppender.Append(m.ColumnType(f))
		conditionalAppender.AppendWithCondition(f.Primary, m.Token(PRIMARY_KEY))
		conditionalAppender.AppendWithCondition(f.Auto, m.Token(AUTO_INCREMENT))
		conditionalAppender.AppendWithCondition(!f.Nullable, m.Token(NOT_NULL))
		if f.Type == CREATED {
			conditionalAppender.Append(m.Token(DEFAULT) + " CURRENT_TIMESTAMP")
		} else if f.Type == UPDATED {
			conditionalAppender.Append(m.Token(DEFAULT) + " CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
		} else {
			conditionalAppender.AppendWithCondition(f.Default != "", m.Token(DEFAULT)+" "+m.WrapDefaultValue(f.Default))
		}

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
	tabWriter := newTabWriter("    ", 2).AppendLine("")
	tabWriter.AppendLine("CREATE TABLE IF NOT EXISTS %s (", d.table.Name)

	tabWriter.LevelUp()

	var updatedField *Field = nil
	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		conditionalAppender.Append(s.ColumnType(f))
		conditionalAppender.AppendWithCondition(f.Primary, s.Token(PRIMARY_KEY))
		conditionalAppender.AppendWithCondition(f.Auto, s.Token(AUTO_INCREMENT))
		conditionalAppender.AppendWithCondition(!f.Nullable, s.Token(NOT_NULL))
		if f.Type == CREATED {
			conditionalAppender.Append(s.Token(DEFAULT) + " CURRENT_TIMESTAMP")
		} else if f.Type == UPDATED {
			conditionalAppender.Append(s.Token(DEFAULT) + " CURRENT_TIMESTAMP")
			updatedField = f
		} else {
			conditionalAppender.AppendWithCondition(f.Default != "", s.Token(DEFAULT)+" "+s.WrapDefaultValue(f.Default))
		}

		trimmedJoined := strings.TrimSpace(strings.Join(conditionalAppender.Slice(), " "))

		suffix := ""
		if i < len(d.table.Fields)-1 {
			suffix = ","
		}
		tabWriter.AppendLine(trimmedJoined + suffix)
	}

	tabWriter.LevelDown()

	tabWriter.AppendLine(");")

	if updatedField != nil {
		//TODO: [SQLITE UPDATE TRIGGER] Sort out sqlite triggers, if using this trigger we get error:
		// An error occurred while commiting the data: too many levels of trigger recursion

		/*tabWriter.AppendLine(`
		  CREATE TRIGGER update_`+updatedField.Name+`
		  AFTER UPDATE ON `+d.table.Name+`
		  FOR EACH ROW
		  BEGIN
		    UPDATE `+d.table.Name+`
		      SET `+updatedField.Name+` = current_timestamp
		      WHERE rowid = old.rowid;
		  END;`)*/
	}

	d.result = tabWriter.CombineLines()
}

func (d *dialectSchemaCreateTableVisitor) VisitPostgres(p *postgres) {
	tabWriter := newTabWriter("    ", 2).AppendLine("")
	tabWriter.AppendLine("CREATE TABLE IF NOT EXISTS %s (", d.table.Name)

	tabWriter.LevelUp()

	var updatedField *Field = nil
	for i, f := range d.table.Fields {
		conditionalAppender := &ConditionalStringSliceAppender{}
		conditionalAppender.Append(f.Name)

		conditionalAppender.Append(p.ColumnType(f))
		conditionalAppender.AppendWithCondition(f.Primary, p.Token(PRIMARY_KEY))
		conditionalAppender.AppendWithCondition(f.Auto, p.Token(AUTO_INCREMENT))
		conditionalAppender.AppendWithCondition(!f.Nullable, p.Token(NOT_NULL))
		if f.Type == CREATED {
			conditionalAppender.Append(p.Token(DEFAULT) + " CURRENT_TIMESTAMP")
		} else if f.Type == UPDATED {
			conditionalAppender.Append(p.Token(DEFAULT) + " CURRENT_TIMESTAMP")
			updatedField = f
		} else {
			conditionalAppender.AppendWithCondition(f.Default != "", p.Token(DEFAULT)+" "+p.WrapDefaultValue(f.Default))
		}

		trimmedJoined := strings.TrimSpace(strings.Join(conditionalAppender.Slice(), " "))

		suffix := ""
		if i < len(d.table.Fields)-1 {
			suffix = ","
		}
		tabWriter.AppendLine(trimmedJoined + suffix)
	}

	tabWriter.LevelDown()

	tabWriter.AppendLine(");")

	if updatedField != nil {
		//TODO: Perhaps make more sense to try and not have duplicate triggers if multiple tables want a trigger to update same field name?
		tabWriter.AppendLine(`
        CREATE OR REPLACE FUNCTION update_` + d.table.Name + `_` + updatedField.Name + `_column()
            RETURNS TRIGGER AS '
          BEGIN
            NEW.` + updatedField.Name + ` = NOW();
            RETURN NEW;
          END;
        ' LANGUAGE 'plpgsql';

        CREATE TRIGGER update_` + updatedField.Name + ` BEFORE UPDATE
          ON ` + d.table.Name + ` FOR EACH ROW EXECUTE PROCEDURE
          update_` + d.table.Name + `_` + updatedField.Name + `_column();`)
	}

	d.result = tabWriter.CombineLines()
}
