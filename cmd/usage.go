package cmd

import (
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

const customUsageTemplate = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}

Resources:{{resourceTypes | trimTrailingWhitespaces}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

var templateFuncs = template.FuncMap{
	"trimTrailingWhitespaces": trimTrailingWhitespaces,
	"resourceTypes":           resourceTypes,
}

func customUsageFunc(cmd *cobra.Command) error {
	tryReadManifest()

	err := tmpl(cmd.OutOrStderr(), customUsageTemplate, cmd)
	if err != nil {
		cmd.PrintErrln(err)
	}
	return err
}

func resourceTypes() string {
	if len(manifest.Config.Resources) > 0 {
		return manifest.ResourceTypeNames()
	}
	return " unknown"
}

func trimTrailingWhitespaces(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}
