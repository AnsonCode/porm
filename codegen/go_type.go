package codegen

import "github.com/vektah/gqlparser/v2/ast"

func addExtraGoStructTags(tags map[string]string, req *CodeGenRequest, col *Column) {
	for _, oride := range req.Settings.Overrides {
		if oride.GoType.StructTags == nil {
			continue
		}
		// if !sdk.Matches(oride, col.Table, req.Catalog.DefaultSchema) {
		// 	// Different table.
		// 	continue
		// }
		// if !sdk.MatchString(oride.ColumnName, col.Name) {
		// 	// Different column.
		// 	continue
		// }
		// Add the extra tags.
		for k, v := range oride.GoType.StructTags {
			tags[k] = v
		}
	}
}

func goType(req *CodeGenRequest, field *ast.FieldDefinition) string {
	// Check if the column's type has been overridden
	// for _, oride := range req.Settings.Overrides {
	// 	if oride.GoType.TypeName == "" {
	// 		continue
	// 	}
	// 	sameTable := Matches(oride, col.Table, req.Catalog.DefaultSchema)
	// 	if oride.Column != "" && MatchString(oride.ColumnName, col.Name) && sameTable {
	// 		return oride.GoType.TypeName
	// 	}
	// }
	col := &Column{
		Name:    field.Name,
		IsArray: field.Type.NamedType == "",
		Type: &Identifier{
			Name: field.Type.NamedType,
		},
		NotNull: field.Type.NonNull,
		Comment: "todo",
	}
	if col.IsArray {
		col.Type = &Identifier{
			Name: field.Type.Elem.NamedType,
		}
	}

	// {{range .}}
	//     type {{.Name}} struct {
	// 		{{range .Fields}}
	// 			{{if eq .Type.NamedType "" }}
	// 				{{Title .Name}} *[]{{.Type.Elem.NamedType}}   `json:"{{.Name}},omitempty"`
	// 			{{else}}
	// 				{{Title .Name}} *{{.Type.NamedType}}  `json:"{{.Name}},omitempty"`
	//     		{{end }}

	// 		{{ end }}
	//     }
	// {{ end }}
	// 	type Column struct {
	// 	Name         string `json:"name,omitempty"`
	// 	NotNull      bool   `json:"not_null,omitempty"`
	// 	IsArray      bool   `json:"is_array,omitempty"`
	// 	Comment      string `json:"comment,omitempty"`
	// 	Length       int32  `json:"length,omitempty"`
	// 	IsNamedParam bool   `json:"is_named_param,omitempty"`
	// 	IsFuncCall   bool   `json:"is_func_call,omitempty"`
	// 	// XXX: Figure out what PostgreSQL calls `foo.id`
	// 	Scope      string      `json:"scope,omitempty"`
	// 	Table      *Identifier `json:"table,omitempty"`
	// 	TableAlias string      `json:"table_alias,omitempty"`
	// 	Type       *Identifier `json:"type,omitempty"`
	// }

	typ := goInnerType(req, col)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func goInnerType(req *CodeGenRequest, col *Column) string {
	// columnType := DataType(col.Type)
	// notNull := col.NotNull || col.IsArray

	// package overrides have a higher precedence
	// for _, oride := range req.Settings.Overrides {
	// 	if oride.GoType.TypeName == "" {
	// 		continue
	// 	}
	// 	if oride.DbType != "" && oride.DbType == columnType && oride.Nullable != notNull {
	// 		return oride.GoType.TypeName
	// 	}
	// }

	// TODO: Extend the engine interface to handle types
	return graphqlType(req, col)
}

// https://chenyitian.gitbooks.io/graphql/content/schema.html#scalar-types
// Int：有符号 32 位整数。
// Float：有符号双精度浮点值。
// String：UTF‐8 字符序列。
// Boolean：true 或者 false。
func graphqlType(req *CodeGenRequest, col *Column) string {
	columnType := DataType(col.Type)
	notNull := col.NotNull || col.IsArray

	switch columnType {

	case "String":
		if notNull {
			return "string"
		}
		return "*string"

	case "Int":
		if notNull {
			return "int32"
		}
		return "*int32"

	case "Float":
		if notNull {
			return "float64"
		}
		return "*float64"

	case "enum":
		// TODO: Proper Enum support
		return "string"

	// case "date", "timestamp", "datetime", "time":
	// 	if notNull {
	// 		return "time.Time"
	// 	}
	// 	return "sql.NullTime"

	case "Boolean":
		if notNull {
			return "bool"
		}
		return "*bool"

	// case "json":
	// 	return "json.RawMessage"

	case "any":
		return "interface{}"

	default:
		// for _, schema := range req.Catalog.Schemas {
		// 	for _, enum := range schema.Enums {
		// 		if enum.Name == columnType {
		// 			if notNull {
		// 				if schema.Name == req.Catalog.DefaultSchema {
		// 					return StructName(enum.Name, req.Settings)
		// 				}
		// 				return StructName(schema.Name+"_"+enum.Name, req.Settings)
		// 			} else {
		// 				if schema.Name == req.Catalog.DefaultSchema {
		// 					return "Null" + StructName(enum.Name, req.Settings)
		// 				}
		// 				return "Null" + StructName(schema.Name+"_"+enum.Name, req.Settings)
		// 			}
		// 		}
		// 	}
		// }
		// if debug.Active {
		// 	log.Printf("Unknown MySQL type: %s\n", columnType)
		// }
		return columnType

	}
}
