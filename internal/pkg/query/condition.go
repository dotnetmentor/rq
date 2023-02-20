package query

import (
	"fmt"
	"regexp"

	"github.com/rs/zerolog/log"
)

type Condition struct {
	Name         string
	MatchMissing bool
	MatchRegExp  bool
	Left         map[string]string
	Right        map[string]string
}

func (c Condition) Match() bool {
	matches := 0

	i := 0
	for lk, lv := range c.Left {
		i++
		rv, ok := c.Right[lk]

		isMatch := false
		if c.MatchRegExp {
			matched, _ := regexp.MatchString(fmt.Sprintf("^%s$", rv), lv)
			isMatch = (ok && matched) || (!ok && c.MatchMissing)
			log.Trace().Str(c.Name, fmt.Sprintf("%s: %s =~ %s", lk, lv, rv)).Str("match", fmt.Sprintf("%v", isMatch)).Str("missing", fmt.Sprintf("%v", !ok)).Msgf("matching %s (key=%d/%d)", c.Name, i, len(c.Left))
		} else {
			isMatch = (ok && rv == lv) || (!ok && c.MatchMissing)
			log.Trace().Str(c.Name, fmt.Sprintf("%s: %s == %s", lk, lv, rv)).Str("match", fmt.Sprintf("%v", isMatch)).Str("missing", fmt.Sprintf("%v", !ok)).Msgf("matching %s (key=%d/%d)", c.Name, i, len(c.Left))
		}

		if !isMatch {
			break
		}
		matches++
	}

	return matches == len(c.Left)
}
