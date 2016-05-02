package main

import "strings"

func JoinImportsForGoFile(imports []string, separator string) string {
	formatted := []string{}
	for _, i := range imports {
		formatted = append(formatted, `"`+i+`"`)
	}
	return strings.Join(formatted, separator)
}

func getImportPathForGoType(goType string) string {
	switch goType {
	case "time.Time":
		return "time"
	default:
		return ""
	}
}
