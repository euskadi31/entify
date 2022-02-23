package main

import (
	stdlog "log"
	"os"
	"path"

	"github.com/euskadi31/entify/pkg/builder"
	"github.com/euskadi31/entify/pkg/spec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EntityContext struct {
	Package    string
	Imports    []string
	StructName string
	Fields     []*EntityField
}

type EntityField struct {
	PropertyName string
	Name         string
	Type         string
}

func main() {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Stack().
		// Caller().
		Logger()

	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	stdlog.SetFlags(stdlog.Llongfile)
	stdlog.SetOutput(logger)

	log.Logger = logger

	//@TODO: use option --out for config this
	dest := path.Join(".", "entify", "entity")

	loader := spec.New()

	//@TODO: use option --provider for config this and check if args[1] is not empty and file readable
	data, err := loader.ParseFile("mysql", os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("open spec file failed")

		os.Exit(1)
	}

	if err := builder.New(data, dest).Build(); err != nil {
		log.Error().Err(err).Msg("generate entity files failed")

		os.Exit(1)
	}

	log.Info().Msg("done")
}
