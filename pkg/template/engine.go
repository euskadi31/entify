package template

import (
	"fmt"
	"io"
	"io/fs"
	"text/template"
)

type Engine struct {
	templates map[string]*template.Template
}

func Must(e *Engine, err error) *Engine {
	if err != nil {
		panic(err)
	}

	return e
}

func New() (*Engine, error) {
	e := &Engine{
		templates: make(map[string]*template.Template),
	}

	if err := e.load(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Engine) load() error {
	err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		pt, serr := template.ParseFS(files, path)
		if serr != nil {
			return fmt.Errorf("template parse fs: %w", serr)
		}

		e.templates[path] = pt

		return nil
	})
	if err != nil {
		return fmt.Errorf("read templates dirs failed: %w", err)
	}

	return nil
}

func (e *Engine) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	if t, ok := e.templates[name]; ok {
		//nolint: wrapcheck
		return t.Execute(wr, data)
	}

	return fmt.Errorf("template %s not found", name)
}
