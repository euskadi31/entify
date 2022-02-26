package main

import (
	stdlog "log"
	"os"
	"path"

	"github.com/euskadi31/entify/pkg/builder"
	"github.com/euskadi31/entify/pkg/spec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	outputFlag   string
	providerFlag string
)

var rootCmd = &cobra.Command{
	Use:   "entify",
	Short: "Entify is a entity generator",
	RunE:  builderRun,
}

func init() {
	dest := path.Join(".", "entify", "entity")

	rootCmd.PersistentFlags().StringVarP(&outputFlag, "out", "o", dest, "out directory (default is ./entify/entity)")
	rootCmd.PersistentFlags().StringVarP(&providerFlag, "provider", "p", "", "out directory (mysql, postgres)")
}

func builderRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Usage()
	}

	filename := args[0]

	loader := spec.New()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Error().Err(err).Msg("")
	}

	data, err := loader.ParseFile(providerFlag, filename)
	if err != nil {
		log.Error().Err(err).Msg("open spec file failed")

		os.Exit(1)
	}

	dest := outputFlag

	if err := builder.New(data, dest).Build(); err != nil {
		log.Error().Err(err).Msg("generate entity files failed")

		os.Exit(1)
	}

	log.Info().Msg("done")

	return nil
}

func main() {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Stack().
		Logger()

	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	stdlog.SetFlags(stdlog.Llongfile)
	stdlog.SetOutput(logger)

	log.Logger = logger

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("")

		os.Exit(1)
	}
}
