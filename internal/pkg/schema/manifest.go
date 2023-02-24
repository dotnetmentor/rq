package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml2 "gopkg.in/yaml.v2"
)

func NewManifest(path string) (Manifest, error) {
	// read manifest file
	var file []byte

	mp, _ := filepath.Abs(path)

	if _, err := os.Stat(mp); err != nil {
		return Manifest{}, err
	}

	bs, err := os.ReadFile(mp)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to read manifest file (path=%s). %v", mp, err)
	}
	file = bs

	if file == nil {
		return Manifest{}, fmt.Errorf("failed to find manifest file path=%v", path)
	}

	// parse
	m := Manifest{}
	err = yaml2.Unmarshal(file, &m)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to parse manifest yaml. %v", err)
	}
	m.parsed = true

	errors := make([]string, 0)
	for i, cr := range m.Config.Resources {
		if cr.Kind == "" {
			errors = append(errors, fmt.Sprintf("config.resources[%d]: missing property \"kind\"", i))
		}
		if cr.Names.Plural == "" {
			errors = append(errors, fmt.Sprintf("config.resources[%d]: missing property \"plural\"", i))
		}
	}

	for _, rt := range m.Config.Resources {
		keys := make([]string, 0)
		for _, key := range m.ResourcesOfType(rt).SelectMany(func(i Resource) string {
			return i.Key
		}) {
			if contains(keys, key) {
				errors = append(errors, fmt.Sprintf("resources[%s].key: key \"%s\" must be unique", rt.Kind, key))
			}
			keys = append(keys, key)
		}
	}

	if len(errors) > 0 {
		return m, fmt.Errorf("invalid resource manifest, error(s):\n\n%s", strings.Join(errors, "\n"))
	}

	return m, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type Manifest struct {
	parsed bool
	Config Config                  `yaml:"config,omitempty"`
	Root   map[string]ResourceList `yaml:"resources,omitempty"`
}

func (m Manifest) Parsed() bool {
	return m.parsed
}

func (m Manifest) ValidateResourceType(t string) (CustomResource, error) {
	rt, ok := m.ResourceType(t)
	if !ok {
		return CustomResource{}, fmt.Errorf("unknown resource type %s, valid resource types: \n%s", t, m.ResourceTypeNames())
	}
	return rt, nil
}

func (m Manifest) ResourcesOfType(cr CustomResource) ResourceList {
	for rk, rl := range m.Root {
		if cr.Names.Match(rk) {
			return rl
		}
	}
	return ResourceList{}
}

func (m Manifest) ResourceTypeNames() string {
	rt := fmt.Sprintln()
	for _, cr := range m.Config.Resources {
		rt += fmt.Sprintln(cr.Names.String())
	}
	return rt
}

func (m Manifest) HasResourceType(t string) bool {
	_, ok := m.ResourceType(t)
	return ok
}

func (m Manifest) ResourceType(t string) (CustomResource, bool) {
	for _, cr := range m.Config.Resources {
		if cr.Names.Match(t) {
			return cr, true
		}
	}
	return CustomResource{}, false
}
