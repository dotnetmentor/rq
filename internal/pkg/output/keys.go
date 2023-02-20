package output

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	Newline string = "newline"
	Xargs   string = "xargs"
	Json    string = "json"
)

func WriteKeys(keys []string, o string) {
	switch o {
	case Newline:
		for _, k := range keys {
			fmt.Println(k)
		}
	case Xargs:
		fmt.Println(strings.Join(keys, " "))
	case Json:
		b, _ := json.Marshal(keys)
		fmt.Println(string(b))
	}
}
