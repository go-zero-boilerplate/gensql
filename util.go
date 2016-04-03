package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/codemodus/kace"
)

func handleDeferAndSetError(errPtr *error) {
	if r := recover(); r != nil {
		switch t := r.(type) {
		case error:
			*errPtr = t
		case string:
			*errPtr = fmt.Errorf("%s", t)
		default:
			*errPtr = fmt.Errorf("%#v", t)
		}
	}
}

type IFieldName interface {
	FieldName() string
}
type IGoType interface {
	GoType() string
}
type IAsGoStructDef interface {
	AsGoStructDef() string
}
type IAsVariableName interface {
	AsVariableName() string
}
type IAsSqlColumn interface {
	AsSqlColumn() string
}

func execTemplateToString(templateString string, data interface{}) (string, error) {
	funcMap := template.FuncMap{
		"AsSqlSelectColumns": func(fields []*GeneratorField) string {
			strs := []string{}
			for _, f := range fields {
				strs = append(strs, f.SqlColumn)
			}
			return strings.Join(strs, ", ")
		},
		"AsSqlParameterizedWhereColumns": func(fields []*GeneratorField) string {
			strs := []string{}
			for _, f := range fields {
				//TODO: This ? symbol is dialect specific
				strs = append(strs, f.SqlColumn+" = ?")
			}
			return strings.Join(strs, " AND ")
		},

		"FieldName": func(s interface{}) string {
			return s.(IFieldName).FieldName()
		},
		"GoType": func(s interface{}) string {
			return s.(IGoType).GoType()
		},
		"AsGoStructDef": func(s interface{}) string {
			return s.(IAsGoStructDef).AsGoStructDef()
		},
		"AsVariableName": func(s interface{}) string {
			return s.(IAsVariableName).AsVariableName()
		},
		"AsSqlColumn": func(s interface{}) string {
			return s.(IAsSqlColumn).AsSqlColumn()
		},

		"CamelFirstUpper": func(s string) string {
			return kace.Camel(s, true)
		},
		"CamelFirstLower": func(s string) string {
			return kace.Camel(s, false)
		},
		"Snake": func(s string) string {
			return kace.Snake(s)
		},

		"CombineForSqlColumns": func(s []yamlField) string {
			strs := []string{}
			for _, y := range s {
				strs = append(strs, y.AsSqlColumn())
			}
			return strings.Join(strs, ", ")
		},
	}

	t, err := template.New("").Funcs(funcMap).Parse(templateString)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	err = t.Execute(&doc, data)
	if err != nil {
		return "", err
	}

	return doc.String(), nil
}

func mustExecTemplateToString(templateString string, data interface{}) string {
	s, err := execTemplateToString(templateString, data)
	if err != nil {
		panic(err)
	}
	return s
}
