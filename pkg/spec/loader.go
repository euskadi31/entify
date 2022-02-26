package spec

import (
	"fmt"
	"os"

	"ariga.io/atlas/schema/schemaspec"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
)

type Loader struct {
	providers map[string]schemaspec.UnmarshalerFunc
}

func New() *Loader {
	return &Loader{
		providers: map[string]schemaspec.UnmarshalerFunc{
			"mysql":    mysql.UnmarshalHCL,
			"mariadb":  mysql.UnmarshalHCL,
			"postgres": postgres.UnmarshalHCL,
		},
	}
}

func (e *Loader) Parse(driver string, b []byte) (schema.Schema, error) {
	d, ok := e.providers[driver]
	if !ok {
		return schema.Schema{}, fmt.Errorf("driver %s not supported", driver)
	}

	var spec schema.Schema

	if err := d(b, &spec); err != nil {
		return schema.Schema{}, fmt.Errorf("unmarshal spec failed: %w", err)
	}

	return spec, nil
}

func (e *Loader) ParseFile(driver string, file string) (schema.Schema, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return schema.Schema{}, fmt.Errorf("read spec file failed: %w", err)
	}

	return e.Parse(driver, b)
}
