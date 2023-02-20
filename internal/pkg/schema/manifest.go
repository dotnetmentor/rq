package schema

import (
	"fmt"
	"os"
	"path/filepath"

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

	// TODO: validate manifest
	// - Resource keys must be unique for each resource type

	return m, nil
}

type Manifest struct {
	Config Config                  `yaml:"config,omitempty"`
	Root   map[string]ResourceList `yaml:"resources,omitempty"`
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
