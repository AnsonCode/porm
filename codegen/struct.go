package codegen

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Struct struct {
	// Table   Identifier
	Name    string
	Kind    string
	IsRes   bool // 标识当前结构是不是响应，如果是响应需要加上前缀
	Fields  []Field
	Comment string
}

func StructName(name string) string {
	// if rename := settings.Rename[name]; rename != "" {
	// 	return rename
	// }
	out := ""
	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}

	// If a name has a digit as its first char, prepand an underscore to make it a valid Go name.
	r, _ := utf8.DecodeRuneInString(out)
	if unicode.IsDigit(r) {
		return "_" + out
	} else {
		return out
	}
}
