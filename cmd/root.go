package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dotnetmentor/rq/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type GlobalOptions struct {
}

var (
	opt GlobalOptions
)

var rootCmd = &cobra.Command{
	Use:          "rq",
	Short:        "rq - short description",
	Long:         `rq - long description`,
	SilenceUsage: true,
	Version:      fmt.Sprintf("%s (commit=%s)", version.Version, version.Commit),
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
