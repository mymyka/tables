package builder

import (
	"strings"
	"tables/pkg/schema"
	"unicode"
)

func Build(tables []schema.Table) map[string]string {
	result := make(map[string]string)

	for _, t := range tables {
		block := "package " + t.Name + "\n\n"

		// Add necessary imports
		imports := buildImports(t)
		if imports != "" {
			block += imports + "\n"
		}

		// Build type aliases
		typeAliases := buildTable(t)
		block += typeAliases + "\n"

		// Build column names struct and variables
		columnStruct := buildColumnNamesStruct(t)
		block += columnStruct + "\n"

		result[t.Name] = block
	}

	return result
}

func buildImports(t schema.Table) string {
	needsTime := false
	needsUUID := false
	needsJSON := false
	needsDecimal := false

	for _, c := range t.Columns {
		switch strings.ToLower(c.Type) {
		case "timestamp", "timestamp with time zone", "timestamp without time zone", "date", "time", "time with time zone", "time without time zone":
			needsTime = true
		case "uuid":
			needsUUID = true
		case "json", "jsonb":
			needsJSON = true
		case "numeric", "decimal":
			needsDecimal = true
		}
	}

	if !needsTime && !needsUUID && !needsJSON && !needsDecimal {
		return ""
	}

	var imports []string
	if needsJSON {
		imports = append(imports, "\"encoding/json\"")
	}
	if needsDecimal {
		imports = append(imports, "\"github.com/shopspring/decimal\"")
	}
	if needsTime {
		imports = append(imports, "\"time\"")
	}
	if needsUUID {
		imports = append(imports, "\"github.com/google/uuid\"")
	}

	return "import (\n\t" + strings.Join(imports, "\n\t") + "\n)"
}

func buildTable(t schema.Table) string {
	block := "\n"

	for _, c := range t.Columns {
		line := buildType(c)
		block += line + "\n"
	}

	return block
}

func buildType(c schema.Column) string {
	line := "type " + toPascalCase(c.Name) + " = "

	if c.Nullable {
		line += "*"
	}

	goType := postgresTypeToGoType(c.Type)
	line += goType

	return line
}

func postgresTypeToGoType(pgType string) string {
	// Normalize the type (remove length specifications, etc.)
	normalizedType := strings.ToLower(strings.TrimSpace(pgType))

	// Handle types with parentheses (e.g., "varchar(255)" -> "varchar")
	if idx := strings.Index(normalizedType, "("); idx != -1 {
		normalizedType = normalizedType[:idx]
	}

	switch normalizedType {
	// Integer types
	case "smallint", "int2":
		return "int16"
	case "integer", "int", "int4":
		return "int32"
	case "bigint", "int8":
		return "int64"
	case "serial", "serial4":
		return "int32"
	case "bigserial", "serial8":
		return "int64"
	case "smallserial", "serial2":
		return "int16"

	// Floating point types
	case "real", "float4":
		return "float32"
	case "double precision", "float8":
		return "float64"

	// Decimal types
	case "numeric", "decimal":
		return "decimal.Decimal"

	// String types
	case "character varying", "varchar":
		return "string"
	case "character", "char":
		return "string"
	case "text":
		return "string"

	// Boolean type
	case "boolean", "bool":
		return "bool"

	// Date/Time types
	case "timestamp", "timestamp with time zone", "timestamptz":
		return "time.Time"
	case "timestamp without time zone":
		return "time.Time"
	case "date":
		return "time.Time"
	case "time", "time with time zone", "timetz":
		return "time.Time"
	case "time without time zone":
		return "time.Time"
	case "interval":
		return "time.Duration"

	// UUID type
	case "uuid":
		return "uuid.UUID"

	// JSON types
	case "json":
		return "json.RawMessage"
	case "jsonb":
		return "json.RawMessage"

	// Binary types
	case "bytea":
		return "[]byte"

	// Network types
	case "inet":
		return "string" // Could use net.IP but string is more common
	case "cidr":
		return "string"
	case "macaddr":
		return "string"
	case "macaddr8":
		return "string"

	// Geometric types
	case "point":
		return "string" // Could create custom types but string is simpler
	case "line":
		return "string"
	case "lseg":
		return "string"
	case "box":
		return "string"
	case "path":
		return "string"
	case "polygon":
		return "string"
	case "circle":
		return "string"

	// Range types
	case "int4range":
		return "string"
	case "int8range":
		return "string"
	case "numrange":
		return "string"
	case "tsrange":
		return "string"
	case "tstzrange":
		return "string"
	case "daterange":
		return "string"

	// Array types (basic handling)
	case "text[]", "varchar[]", "character varying[]":
		return "[]string"
	case "integer[]", "int4[]":
		return "[]int32"
	case "bigint[]", "int8[]":
		return "[]int64"
	case "smallint[]", "int2[]":
		return "[]int16"
	case "boolean[]", "bool[]":
		return "[]bool"
	case "real[]", "float4[]":
		return "[]float32"
	case "double precision[]", "float8[]":
		return "[]float64"

	// Money type
	case "money":
		return "string" // Could use decimal.Decimal but string is safer

	// Enum types (generic handling)
	case "enum":
		return "string"

	// XML type
	case "xml":
		return "string"

	// Bit string types
	case "bit":
		return "string"
	case "bit varying", "varbit":
		return "string"

	// PostgreSQL specific types
	case "tsvector":
		return "string"
	case "tsquery":
		return "string"
	case "pg_lsn":
		return "string"
	case "pg_snapshot":
		return "string"
	case "txid_snapshot":
		return "string"

	// Default fallback
	default:
		// Handle array types that weren't caught above
		if strings.HasSuffix(normalizedType, "[]") {
			return "[]interface{}"
		}
		// Unknown type, default to string
		return "string"
	}
}

func buildColumnNamesStruct(t schema.Table) string {
	var block strings.Builder

	// Build struct type
	structName := t.Name + "ColumnNames"
	block.WriteString("type " + structName + " struct {\n")

	for _, c := range t.Columns {
		fieldName := toPascalCase(c.Name)
		block.WriteString("\t" + fieldName + " string\n")
	}

	block.WriteString("}\n\n")

	// Build C variable with column names
	block.WriteString("var C = " + structName + "{\n")

	for _, c := range t.Columns {
		fieldName := toPascalCase(c.Name)
		block.WriteString("\t" + fieldName + ": \"" + c.Name + "\",\n")
	}

	block.WriteString("}\n\n")

	// Build Table variable
	block.WriteString("var Table = \"" + t.Name + "\"\n")

	return block.String()
}

// Helper function to capitalize the first letter
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	// Convert first character to uppercase, keep rest as is
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Helper function to convert snake_case to PascalCase
func toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}

	// Split by underscore and capitalize each part
	parts := strings.Split(s, "_")
	var result strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(capitalizeFirst(part))
		}
	}

	return result.String()
}
