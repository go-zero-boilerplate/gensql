package main

import (
	"strings"

	"github.com/codemodus/kace"
)

type yamlField string

func (y yamlField) Split() (name, goType string, args []string) {
	str := string(y)
	for strings.Contains(str, "  ") {
		str = strings.Replace(str, "  ", " ", -1)
	}

	strs := strings.Split(str, " ")
	return strs[0], strs[1], strs[2:]
}

func (y yamlField) AsGoStructDef() string {
	_, goType, _ := y.Split()
	return y.FieldName() + " " + goType
}

func (y yamlField) FieldName() string {
	fieldName, _, _ := y.Split()
	return kace.Camel(fieldName, true)
}

func (y yamlField) AsVariableName() string {
	fieldName, _, _ := y.Split()
	return kace.Camel(fieldName, false)
}

func (y yamlField) GoType() string {
	_, goType, _ := y.Split()
	return goType
}

func (y yamlField) AsSqlColumn() string {
	return kace.Snake(string(y.FieldName()))
}

type predefinedField string

func (p predefinedField) String() string { return string(p) }

func (p predefinedField) GoType() string {
	switch p {
	case "created":
		return "time.Time"
	case "updated":
		return "time.Time"
	default:
		panic("GoType not supported for PREDEFINED field name: " + string(p))
	}
}

func (p predefinedField) FieldName() string {
	return kace.Camel(string(p), true)
}

func (p predefinedField) AsVariableName() string {
	return kace.Camel(string(p), false)
}

func (p predefinedField) AsGoStructDef() string {
	goType := p.GoType()
	return p.FieldName() + " " + goType
}
