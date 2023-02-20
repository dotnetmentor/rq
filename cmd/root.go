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
	FilePath string
	LogLevel string
}

var (
	opt GlobalOptions
)

var RootCmd = &cobra.Command{
	Use:          "rq",
	Short:        "rq - query resources",
	Long:         `rq - for querying resources`,
	SilenceUsage: true,
	Version:      fmt.Sprintf("%s (commit=%s)", version.Version, version.Commit),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level, err := zerolog.ParseLevel(opt.LogLevel)
		if err != nil {
			fmt.Printf("invalid log level %s, err: %s", opt.LogLevel, err)
			os.Exit(1)
		}
		log.Logger = log.Level(level).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	},
}

func Execute() {
	RootCmd.SetOut(os.Stdout)
	RootCmd.SetErr(os.Stderr)

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func environmentOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&opt.FilePath, "file", "f", environmentOrDefault("RQ_DEFAULT_FILE", "rq.yaml"), "Manifest file (eg. rq.yaml)")
	RootCmd.PersistentFlags().StringVarP(&opt.LogLevel, "level", "l", "info", "Log level (eg. trace/debug/warn/error)")
}
