package builder

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path"
	"regexp"
	"strings"

	"ariga.io/atlas/sql/schema"
	tmpl "github.com/euskadi31/entify/pkg/template"
	"github.com/euskadi31/entify/pkg/types"
	"github.com/gertd/go-pluralize"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/modfile"
)

var pattern = regexp.MustCompile(`__([a-z0-9\-]+)__`)

var pluralizeClient *pluralize.Client

var signMap map[string]string

var typeIntMap map[string]string

var typeUintMap map[string]string

var typeNullableIntMap map[string]string

func init() {
	pluralizeClient = pluralize.NewClient()

	signMap = map[string]string{
		"id":  "ID",
		"url": "URL",
		"uri": "URI",
	}

	typeIntMap = map[string]string{
		"tinyint":   "int8",
		"smallint":  "int16",
		"mediumint": "int32",
		"int":       "int32",
		"bigint":    "int64",
	}

	typeUintMap = map[string]string{
		"tinyint":   "uint8",
		"smallint":  "uint16",
		"mediumint": "uint32",
		"int":       "uint32",
		"bigint":    "uint64",
	}

	typeNullableIntMap = map[string]string{
		"tinyint":   "sql.NullInt16",
		"smallint":  "sql.NullInt16",
		"mediumint": "sql.NullInt32",
		"int":       "sql.NullInt32",
		"bigint":    "sql.NullInt64",
	}
}

type Builder struct {
	dest string
	spec schema.Schema
	data *types.Data
	tpl  *tmpl.Engine
}

func New(spec schema.Schema, dest string) *Builder {
	return &Builder{
		dest: dest,
		spec: spec,
		data: &types.Data{
			Package: "entity",
		},
		tpl: tmpl.Must(tmpl.New()),
	}
}

func (b *Builder) processSpec() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	moddir := findModuleRoot(dir)

	mod, err := readModFile(moddir)
	if err != nil {
		panic(err)
	}

	for _, t := range b.spec.Tables {
		importsMap := map[string]struct{}{}
		imports := []string{}

		pks := []*types.Field{}
		pksMap := map[string]struct{}{}

		for _, pk := range t.PrimaryKey.Parts {
			pksMap[pk.C.Name] = struct{}{}
		}

		fields := []*types.Field{}

		for _, col := range t.Columns {
			ct := ColumnTypeToType(col.Type)

			if ct.Package != "" {
				if _, ok := importsMap[ct.Package]; !ok {
					imports = append(imports, ct.Package)
					importsMap[ct.Package] = struct{}{}
				}
			}

			field := &types.Field{
				PropertyName:           ColumnNameToPropertyName(col.Name),
				VariableName:           ColumnNameToVariableName(col.Name),
				Name:                   col.Name,
				Type:                   ct.Type,
				SQLType:                ct.NullableSQLType,
				Nullable:               col.Type.Null,
				TypeKind:               ct.TypeKind,
				NullableSQLAccessValue: ct.NullableSQLAccessValue,
				DefaultValue:           ct.DefaultValue,
			}

			// set field to primary keys
			if _, ok := pksMap[col.Name]; ok {
				pks = append(pks, field)
			}

			fields = append(fields, field)
		}

		autoIncr := false

		if len(pks) == 1 && pks[0].TypeKind == types.FieldTypeKindNumber {
			// @TODO: check if auto_increment is set
			autoIncr = true
		}

		b.data.Entities = append(b.data.Entities, &types.Entity{
			ReceiverVarName:    TableNameToReceiver(t.Name),
			Module:             modfile.ModulePath(mod),
			Name:               t.Name,
			VariableName:       TableNameToVariableName(t.Name),
			Filename:           TableNameToFileName(t.Name),
			StructName:         TableNameToStructName(t.Name),
			PackageName:        TableNameToPackageName(t.Name),
			Imports:            imports,
			Fields:             fields,
			FieldsCount:        len(fields),
			PrimaryKeys:        pks,
			PrimaryKeyAutoIncr: autoIncr,
		})
	}
}

// name = predicate/predicate.go.tmpl => predicate/predicate.go, placeholder = ""
// name = __entity__/__entity__.go => user/user.go, placeholder = "user".
func (b *Builder) getDestFilename(name string, placeholders map[string]string) string {
	filename := strings.Replace(name, ".go.tmpl", ".go", 1)

	if len(placeholders) > 0 {
		matches := pattern.FindAllStringSubmatch(name, 10)

		for _, match := range matches {
			filename = strings.Replace(filename, match[0], placeholders[match[1]], 1)
		}
	}

	filename = path.Join(strings.Split(filename, "/")...)

	return path.Join(b.dest, filename)
}

func (b *Builder) render(name string, placeholders map[string]string, data interface{}) error {
	if data == nil {
		data = b.data
	}

	dest := b.getDestFilename(name, placeholders)

	log.Debug().Msgf("create %s", dest)

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create %s file failed: %w", dest, err)
	}

	defer f.Close()

	buf := bytes.NewBuffer(nil)

	if err := b.tpl.ExecuteTemplate(buf, name, data); err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	content, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format %s source failed: %w", name, err)
	}

	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("write %s file failed: %w", dest, err)
	}

	return nil
}

func (b *Builder) generatePredicate() error {
	return b.render("predicate/predicate.go.tmpl", nil, nil)
}

func (b *Builder) generateClient() error {
	return b.render("client.go.tmpl", nil, nil)
}

func (b *Builder) generateEntity(entity *types.Entity) error {
	if err := b.render("__entity-package__/__entity-file__.go.tmpl", map[string]string{
		"entity-package": entity.PackageName,
		"entity-file":    entity.Filename,
	}, entity); err != nil {
		return fmt.Errorf("generate entity file: %w", err)
	}

	if err := b.render("__entity-package__/where.go.tmpl", map[string]string{
		"entity-package": entity.PackageName,
	}, entity); err != nil {
		return fmt.Errorf("generate where file: %w", err)
	}

	de := &types.DataEntity{
		Package: b.data.Package,
		Module:  entity.Module,
		Entity:  entity,
	}

	for _, f := range []string{
		"__entity-file___client.go.tmpl",
		"__entity-file___create.go.tmpl",
		"__entity-file___delete.go.tmpl",
		"__entity-file___mutation.go.tmpl",
		"__entity-file___query.go.tmpl",
		"__entity-file___update.go.tmpl",
		"__entity-file__.go.tmpl",
	} {
		if err := b.render(f, map[string]string{
			"entity-file": entity.Filename,
		}, de); err != nil {
			return fmt.Errorf("generate entity file: %w", err)
		}
	}

	return nil
}

func (b *Builder) createFolders() error {
	log.Debug().Msgf("create destination dir: %s", b.dest)

	if err := os.MkdirAll(b.dest, 0755); err != nil {
		return fmt.Errorf("create destination dir failed: %w", err)
	}

	for _, e := range b.data.Entities {
		dest := path.Join(b.dest, e.PackageName)

		log.Debug().Msgf("create entity dir: %s", dest)

		if err := os.MkdirAll(dest, 0755); err != nil {
			return fmt.Errorf("create %s dir failed: %w", e.PackageName, err)
		}
	}

	dest := path.Join(b.dest, "predicate")

	log.Debug().Msgf("create predicate dir: %s", dest)

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("create predicate dir failed: %w", err)
	}

	return nil
}

func (b *Builder) Build() error {
	b.processSpec()

	if err := b.createFolders(); err != nil {
		return fmt.Errorf("create destination structure folders: %w", err)
	}

	if err := b.generateClient(); err != nil {
		return fmt.Errorf("generate client: %w", err)
	}

	if err := b.generatePredicate(); err != nil {
		return fmt.Errorf("generate predicate: %w", err)
	}

	for _, entity := range b.data.Entities {
		if err := b.generateEntity(entity); err != nil {
			return fmt.Errorf("generate entity: %w", err)
		}
	}

	return nil
}
