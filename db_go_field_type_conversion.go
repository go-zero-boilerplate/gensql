package main

func dbGoFieldTypeFromGoType(goType string, isNullable bool) (importPackagePath, dbGoType, dbDotSuffix string, mustCast bool) {
	if !isNullable {
		return "", goType, "", false
	}

	switch goType {
	case "string":
		return "gopkg.in/guregu/null.v3", "null.String", ".String", false
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "gopkg.in/guregu/null.v3", "null.Int", ".Int64", goType != "int64"
	case "float32", "float64":
		return "gopkg.in/guregu/null.v3", "null.Float", ".Float64", goType != "float64"
	case "bool":
		return "gopkg.in/guregu/null.v3", "null.Bool", ".Bool", false
	case "time.Time":
		return "gopkg.in/guregu/null.v3", "null.Time", ".Time", false
	}

	return "", goType, "", false
}
