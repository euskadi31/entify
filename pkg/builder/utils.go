//nolint: goconst
package builder

import (
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"ariga.io/atlas/sql/schema"
	"github.com/davecgh/go-spew/spew"
	"github.com/euskadi31/entify/pkg/types"
	"github.com/iancoleman/strcase"
)

const (
	pkgDatabaseSQL = "database/sql"
)

func TableNameToReceiver(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		parts[i] = strings.ToLower(part[0:1])
	}

	return strings.Join(parts, "")
}

func TableNameToStructName(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		parts[i] = pluralizeClient.Singular(part)
	}

	return strcase.ToCamel(strings.Join(parts, "_"))
}

func TableNameToFileName(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		parts[i] = pluralizeClient.Singular(part)
	}

	return strings.Join(parts, "_")
}

func TableNameToPackageName(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		parts[i] = pluralizeClient.Singular(part)
	}

	return strings.Join(parts, "")
}

func ColumnNameToPropertyName(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		if v, ok := signMap[strings.ToLower(part)]; ok {
			parts[i] = v
		} else {
			parts[i] = strcase.ToCamel(part)
		}
	}

	return strings.Join(parts, "")
}

func TableNameToVariableName(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		if v, ok := signMap[strings.ToLower(part)]; ok && i > 0 {
			parts[i] = v
		} else {
			parts[i] = strcase.ToCamel(part)
		}

		parts[i] = pluralizeClient.Singular(parts[i])
	}

	parts[0] = strings.ToLower(parts[0])

	return strings.Join(parts, "")
}

func ColumnNameToVariableName(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		if v, ok := signMap[strings.ToLower(part)]; ok && i > 0 {
			parts[i] = v
		} else {
			parts[i] = strcase.ToCamel(part)
		}
	}

	parts[0] = strings.ToLower(parts[0])

	vn := strings.Join(parts, "")

	if token.Lookup(vn).IsKeyword() {
		return vn[0:1]
	}

	return vn
}

func ColumnTypeToSQLType(colType *schema.ColumnType) (string, string, string) {
	switch t := colType.Type.(type) {
	case *schema.StringType:
		goType := "string"
		sqlType := "string"
		goPkg := ""

		if colType.Null {
			sqlType = "sql.NullString"
			goPkg = pkgDatabaseSQL
		}

		return goType, sqlType, goPkg
	case *schema.BoolType:
		goType := "bool"
		sqlType := "bool"
		goPkg := ""

		if colType.Null {
			sqlType = "sql.NullBool"
			goPkg = pkgDatabaseSQL
		}

		return goType, sqlType, goPkg
	case *schema.IntegerType:
		goType := ""
		sqlType := ""
		goPkg := ""

		if t.Unsigned {
			if v, ok := typeUintMap[t.T]; ok {
				goType = v
				sqlType = v
			}
		} else {
			if v, ok := typeIntMap[t.T]; ok {
				goType = v
				sqlType = v
			}
		}

		if colType.Null {
			if v, ok := typeNullableIntMap[t.T]; ok {
				sqlType = v
				goPkg = pkgDatabaseSQL
			}
		}

		return goType, sqlType, goPkg

	case *schema.FloatType:
		goType := "float64"
		sqlType := "float64"
		goPkg := ""

		if colType.Null {
			sqlType = "sql.NullFloat64"
			goPkg = pkgDatabaseSQL
		}

		return goType, sqlType, goPkg

	case *schema.TimeType:
		goType := "time.Time"
		sqlType := "time.Time"
		goPkg := "time"

		if colType.Null {
			sqlType = "sql.NullTime"
			goPkg = pkgDatabaseSQL
		}

		return goType, sqlType, goPkg
	}

	return "", "", ""
}

type ColumnType struct {
	Type                   string
	TypeKind               types.FieldTypeKind
	SQLType                string
	DefaultValue           string
	NullableSQLType        string
	Package                string
	NullablePackage        string
	NullableSQLAccessValue string
}

func ColumnTypeToType(colType *schema.ColumnType) *ColumnType {
	switch t := colType.Type.(type) {
	case *schema.StringType:
		return &ColumnType{
			Type:                   "string",
			TypeKind:               types.FieldTypeKindString,
			SQLType:                "string",
			DefaultValue:           `""`,
			NullableSQLType:        "sql.NullString",
			NullablePackage:        pkgDatabaseSQL,
			NullableSQLAccessValue: "String",
		}
	case *schema.BoolType:
		return &ColumnType{
			Type:                   "bool",
			TypeKind:               types.FieldTypeKindBool,
			SQLType:                "bool",
			DefaultValue:           `false`,
			NullableSQLType:        "sql.NullBool",
			NullablePackage:        pkgDatabaseSQL,
			NullableSQLAccessValue: "Bool",
		}
	case *schema.IntegerType:
		ct := &ColumnType{
			TypeKind:        types.FieldTypeKindNumber,
			DefaultValue:    `0`,
			NullablePackage: pkgDatabaseSQL,
		}

		if t.Unsigned {
			if v, ok := typeUintMap[t.T]; ok {
				ct.Type = v
				ct.SQLType = v
			}
		} else {
			if v, ok := typeIntMap[t.T]; ok {
				ct.Type = v
				ct.SQLType = v
			}
		}

		if v, ok := typeNullableIntMap[t.T]; ok {
			ct.NullableSQLType = v
			ct.NullableSQLAccessValue = strings.Replace(v, "sql.Null", "", 1)
		}

		return ct

	case *schema.FloatType:
		return &ColumnType{
			Type:                   "float64",
			TypeKind:               types.FieldTypeKindNumber,
			SQLType:                "float64",
			DefaultValue:           `0.0`,
			NullableSQLType:        "sql.NullFloat64",
			NullablePackage:        pkgDatabaseSQL,
			NullableSQLAccessValue: "Float64",
		}

	case *schema.TimeType:
		return &ColumnType{
			Type:                   "time.Time",
			TypeKind:               types.FieldTypeKindDate,
			SQLType:                "time.Time",
			DefaultValue:           `time.Time{}`,
			NullableSQLType:        "sql.NullTime",
			Package:                "time",
			NullablePackage:        pkgDatabaseSQL,
			NullableSQLAccessValue: "Time",
		}
	case *schema.JSONType:
		return &ColumnType{
			Type:                   "json.RawMessage",
			TypeKind:               types.FieldTypeKindJson,
			SQLType:                "string",
			DefaultValue:           `json.RawMessage{}`,
			NullableSQLType:        "sql.NullString",
			Package:                "encoding/json",
			NullablePackage:        pkgDatabaseSQL,
			NullableSQLAccessValue: "String",
		}
	default:
		spew.Dump(t)
	}

	return nil
}

func findModuleRoot(dir string) (root string) {
	if dir == "" {
		panic("dir not set")
	}

	dir = filepath.Clean(dir)

	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}

		d := filepath.Dir(dir)

		if d == dir {
			break
		}

		dir = d
	}

	return ""
}

func reverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func findCurrentModulePath(dir string) (root string) {
	if dir == "" {
		panic("dir not set")
	}

	dir = filepath.Clean(dir)

	parts := []string{}

	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return path.Join(reverseSlice(parts)...)
		}

		d := filepath.Dir(dir)

		if d == dir {
			break
		}

		parts = append(parts, path.Base(dir))

		dir = d
	}

	return path.Join(reverseSlice(parts)...)
}

//nolint:wrapcheck
func readModFile(dir string) ([]byte, error) {
	return os.ReadFile(path.Join(dir, "go.mod"))
}
