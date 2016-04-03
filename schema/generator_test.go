package schema

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func cleanSQLStatement(s string) string {
	r := s
	r = strings.Replace(r, "\t", " ", -1)
	r = strings.Replace(r, "\n", " ", -1)
	for strings.Contains(r, "  ") {
		r = strings.Replace(r, "  ", " ", -1)
	}
	return r
}

func TestGenerateShema(t *testing.T) {
	Convey("Testing GenerateShema method", t, func() {
		table := NewTableBuilder("blog_tab1").
			Field("id", INTEGER, true, true, 0).
			Field("name", VARCHAR, false, false, 200).
			Primary("id").
			Index("unique_name", true, "name").
			Build()
		var schema string

		schema = GenerateSchema(NewMysqlSchemaDialect(), table)
		So(cleanSQLStatement(schema), ShouldEqual, cleanSQLStatement(`
			CREATE TABLE IF NOT EXISTS blog_tab1 (                     
		        id INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL,    
		        name VARCHAR(200) NOT NULL                         
			);                                                         
                                                           
			CREATE UNIQUE INDEX unique_name ON blog_tab1 ( name );`))

		schema = GenerateSchema(NewPostgresSchemaDialect(), table)
		So(cleanSQLStatement(schema), ShouldEqual, cleanSQLStatement(`
			CREATE TABLE IF NOT EXISTS blog_tab1 (
		        id SERIAL PRIMARY KEY  NOT NULL,
		        name VARCHAR(200) NOT NULL
			);

			CREATE UNIQUE INDEX unique_name ON blog_tab1 ( name );`))

		schema = GenerateSchema(NewSqliteSchemaDialect(), table)
		So(cleanSQLStatement(schema), ShouldEqual, cleanSQLStatement(`
			CREATE TABLE IF NOT EXISTS blog_tab1 (
			    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			    name TEXT NOT NULL
			);

			CREATE UNIQUE INDEX unique_name ON blog_tab1 ( name );`))
	})
}
