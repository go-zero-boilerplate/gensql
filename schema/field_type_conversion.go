package schema

import (
	"fmt"
)

var goToFieldType = map[string]FieldType{
	"bool":      BOOLEAN,
	"int":       INTEGER,
	"int8":      INTEGER,
	"int16":     INTEGER,
	"int32":     INTEGER,
	"int64":     BIGINT,
	"uint":      INTEGER,
	"uint8":     INTEGER,
	"uint16":    INTEGER,
	"uint32":    INTEGER,
	"uint64":    BIGINT,
	"float32":   REAL,
	"float64":   REAL,
	"[]byte":    BLOB,
	"string":    VARCHAR,
	"time.Time": DATETIME,
}

func GoToFieldType(goType string) (FieldType, error) {
	if ft, ok := goToFieldType[goType]; ok {
		return ft, nil
	} else {
		return nil, fmt.Errorf("Go type '%s' cannot be converted to SQL FieldType", goType)
	}
}
