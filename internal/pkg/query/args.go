package query

import "strings"

func ParseArgs(args []string) (kvp map[string]string, err error) {
	kvp = make(map[string]string)
	for _, a := range args {
		parts := strings.Split(a, "=")
		kvp[parts[0]] = parts[1]
	}
	return
}
