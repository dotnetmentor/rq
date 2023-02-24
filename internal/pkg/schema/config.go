package schema

import (
	"fmt"
	"strings"
)

type Config struct {
	Resources CustomResourceList `yaml:"resources,omitempty"`
}

type CustomResourceList []CustomResource

type CustomResource struct {
	Kind  string              `yaml:"kind,omitempty"`
	Names CustomResourceNames `yaml:"names,omitempty"`
}

type CustomResourceNames struct {
	Plural   string `yaml:"plural,omitempty"`
	Singular string `yaml:"singular,omitempty"`
	Short    string `yaml:"short,omitempty"`
}

func (n CustomResourceNames) Match(t string) bool {
	return n.Plural == t || n.Singular == t || n.Short == t
}

func (n CustomResourceNames) String() string {
	out := ""
	if n.Plural != "" && n.Singular != "" && strings.HasPrefix(n.Plural, n.Singular) {
		pe := strings.ReplaceAll(n.Plural, n.Singular, "")
		out = fmt.Sprintf("  - %s(%s)", n.Singular, pe)
	} else {
		v := []string{n.Plural}
		if n.Singular != "" {
			v = append(v, n.Singular)
		}
		out = fmt.Sprintf("  - %s", strings.Join(v, "/"))
	}

	if n.Short != "" {
		out = fmt.Sprintf("%s (short: %s)", out, n.Short)
	}
	return out
}
