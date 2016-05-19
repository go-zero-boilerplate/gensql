package main

func dbGoFieldTypeFromGoType(goType string, isNullable bool) (importPackagePath, dbGoType, dbDotSuffix string, mustCast bool) {
	if !isNullable {
		return "", goType, "", false
	}

	switch goType {
	case "string":
		return "gopkg.in/guregu/null.v3/zero", "zero.String", ".String", false
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "gopkg.in/guregu/null.v3/zero", "zero.Int", ".Int64", goType != "int64"
	case "float32", "float64":
		return "gopkg.in/guregu/null.v3/zero", "zero.Float", ".Float64", goType != "float64"
	case "bool":
		return "gopkg.in/guregu/null.v3/zero", "zero.Bool", ".Bool", false
	case "time.Time":
		return "gopkg.in/guregu/null.v3/zero", "zero.Time", ".Time", false
	}

	return "", goType, "", false
}
