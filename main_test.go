package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/dotnetmentor/rq/cmd"
	"github.com/dotnetmentor/rq/internal/pkg/test"

	"github.com/andreyvit/diff"
	"github.com/spf13/cobra"
)

var cases = []struct {
	name string
	file string
	args []string
}{
	{"query_packages_noargs", "rq.yaml", []string{"query", "packages"}},
	{"query_packages_condition_1", "rq.yaml", []string{"query", "packages", "-p", "environment=local", "-d"}},
	{"query_packages_condition_2", "rq.yaml", []string{"query", "packages", "-p", "environment=dev", "-d"}},
	{"query_packages_property_match_1", "rq.yaml", []string{"query", "packages", "-p", "environment=local", "-p", "build=true", "-d"}},
	{"query_packages_property_match_2", "rq.yaml", []string{"query", "packages", "-p", "environment=local", "-p", "build=false"}},
	{"query_packages_property_sort", "rq.yaml", []string{"query", "packages", "--sort"}},

	{"query_tenants_noargs", "rq.yaml", []string{"query", "tenants", "--sort"}},
	{"query_tenants_property_match_1", "rq.yaml", []string{"query", "tenants", "--sort", "-p", "default=true"}},
	{"query_tenants_property_match_2", "rq.yaml", []string{"query", "tenants", "--sort", "-p", "tenancy=shared"}},
}

func TestCommands(t *testing.T) {
	cleanupTestOutput()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			expectedFile := fmt.Sprintf("./testdata/%s.expected", tt.name)
			actualFile := fmt.Sprintf("./testdata/%s.actual", tt.name)

			args := make([]string, 0)
			args = append(args, "--file", fmt.Sprintf("./testdata/%s", tt.file))
			args = append(args, tt.args...)

			out := test.CaptureOutput(func() {
				_, _, err := executeCommandC(cmd.RootCmd, args...)
				if err != nil {
					t.Error(err)
				}
			})

			f, _ := os.OpenFile(actualFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			defer f.Close()
			f.WriteString(out.Stdout)

			expected, _ := os.ReadFile(expectedFile)
			actual, _ := os.ReadFile(actualFile)
			if string(expected) != string(actual) {
				t.Errorf("\ncase: %s\ncommand:%s\n\n%s", tt.name, args, diff.LineDiff(string(expected), string(actual)))
			}
		})
	}
}

func cleanupTestOutput() {
	files, err := filepath.Glob("./testdata/*.actual")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
