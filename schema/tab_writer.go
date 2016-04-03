package schema

import (
	"fmt"
	"strings"
)

func newTabWriter(indentString string, initialIndentLevel int) *tabWriter {
	return &tabWriter{
		indentString: indentString,
		indentLevel:  initialIndentLevel,
	}
}

type tabWriter struct {
	indentString string
	indentLevel  int
	lines        []string
}

func (t *tabWriter) LevelUp() *tabWriter {
	t.indentLevel++
	return t
}
func (t *tabWriter) LevelDown() *tabWriter {
	t.indentLevel--
	return t
}

func (t *tabWriter) AppendLine(format string, a ...interface{}) *tabWriter {
	prefix := ""
	if t.indentLevel > 0 {
		prefix = strings.Repeat(t.indentString, t.indentLevel)
	}
	t.lines = append(t.lines, prefix+fmt.Sprintf(format, a...))
	return t
}

func (t *tabWriter) CombineLines() string {
	return strings.Join(t.lines, "\n")
}
