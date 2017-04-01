package consullauncher

import (
	"html/template"
	"github.com/Masterminds/sprig"
	"fmt"
	"bytes"
)

func templatePostIteration(configFiles []Entry, dest string) {

}

var templatePlugin = Plugin{
	PostIteration: func([] Entry, string) {},
	CheckActivation:func(flag uint64) bool {
		return flag | 2 > 0
	},
	ProcessContent:func(content []byte) []byte {
		template, err := template.New("template").Funcs(sprig.FuncMap()).Parse(string(content))
		if (err != nil) {
			fmt.Println("Error on parsing template")
			return content
		}
		var result bytes.Buffer
		template.Execute(&result, make(map[string]string))
		return result.Bytes()
	},
}