package types

type Data struct {
	Package  string
	Entities []*Entity
}

type Entity struct {
	ReceiverVarName    string
	Module             string
	Filename           string
	Name               string
	PackageName        string
	StructName         string
	VariableName       string
	Imports            []string
	Fields             []*Field
	FieldsCount        int
	PrimaryKeys        []*Field
	PrimaryKeyAutoIncr bool
}

type FieldTypeKind int8

const (
	FieldTypeKindUnknown FieldTypeKind = iota
	FieldTypeKindNumber
	FieldTypeKindString
	FieldTypeKindBool
	FieldTypeKindDate
	FieldTypeKindJson
)

type Field struct {
	Name                   string
	PropertyName           string
	VariableName           string
	TypeKind               FieldTypeKind
	Type                   string
	SQLType                string
	Nullable               bool
	NullableSQLAccessValue string
	DefaultValue           string
}

type DataEntity struct {
	Package string
	Module  string
	Entity  *Entity
}
