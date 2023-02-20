package query

func AllOf(conditions ...Condition) bool {
	for _, c := range conditions {
		if !c.Match() {
			return false
		}
	}
	return true
}

func OneOf(conditions ...Condition) bool {
	if len(conditions) == 0 {
		return true
	}
	for _, c := range conditions {
		if c.Match() {
			return true
		}
	}
	return false
}
