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
	if strings.HasPrefix(n.Plural, n.Singular) {
		pe := strings.ReplaceAll(n.Plural, n.Singular, "")
		return fmt.Sprintf("  - %s(%s) (short: %s)", n.Singular, pe, n.Short)
	}
	return fmt.Sprintf("  - %s/%s (short: %s)", n.Plural, n.Singular, n.Short)
}
