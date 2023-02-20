package schema

type Resource struct {
	Key        string              `json:"key,omitempty" yaml:"key,omitempty"`
	Properties map[string]string   `json:"properties,omitempty" yaml:"properties,omitempty"`
	Conditions []map[string]string `json:"when,omitempty" yaml:"when,omitempty"`
}

type ResourceList []Resource

func (l ResourceList) Where(condition func(t Resource) bool) (fl ResourceList) {
	fl = make(ResourceList, 0)
	for _, t := range l {
		if condition(t) {
			fl = append(fl, t)
		}
	}
	return
}

func (l ResourceList) SelectMany(selector func(i Resource) string) (out []string) {
	for _, i := range l {
		out = append(out, selector(i))
	}
	return
}
