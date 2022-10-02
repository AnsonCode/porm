package codegen

import (
	"fmt"
	"sort"
	"strings"
)

type Field struct {
	Name string // CamelCased name for Go
	// DBName   string // Name as used in the DB
	Type     string
	Tags     map[string]string
	Comment  string
	IsObject bool //  标识该字段是否为对象
	// Struct   *Struct // 关联哪个对象
}

func (gf Field) Tag() string {
	tags := make([]string, 0, len(gf.Tags))
	for key, val := range gf.Tags {
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", key, val))
	}
	if len(tags) == 0 {
		return ""
	}
	sort.Strings(tags)
	return strings.Join(tags, " ")
}

func JSONTagName2(name string, omitempty bool) string {
	if omitempty {
		return name + ",omitempty"

	}
	return name
}

func JSONTagName(name string, settings *Settings) string {
	style := settings.Go.JsonTagsCaseStyle
	if style == "" || style == "none" {
		return name
	} else {
		return SetCaseStyle(name, style)
	}
}

func SetCaseStyle(name string, style string) string {
	switch style {
	case "camel":
		return toCamelCase(name)
	case "pascal":
		return toPascalCase(name)
	case "snake":
		return toSnakeCase(name)
	default:
		panic(fmt.Sprintf("unsupported JSON tags case style: '%s'", style))
	}
}

func toSnakeCase(s string) string {
	return strings.ToLower(s)
}

func toCamelCase(s string) string {
	return toCamelInitCase(s, false)
}

func toPascalCase(s string) string {
	return toCamelInitCase(s, true)
}

func toCamelInitCase(name string, initUpper bool) string {
	out := ""
	for i, p := range strings.Split(name, "_") {
		if !initUpper && i == 0 {
			out += p
			continue
		}
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}
