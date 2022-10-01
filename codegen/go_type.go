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
	col := &Column{
		Name:    field.Name,
		IsArray: field.Type.NamedType == "",
		Type: &Identifier{
			Name: field.Type.NamedType,
		},
		NotNull: field.Type.NonNull,
		Comment: "todo1",
	}
	if col.IsArray {
		col.Type = &Identifier{
			Name: field.Type.Elem.NamedType,
		}
	}
	// if field.SelectionSet != nil {
	// 			// 是否数组的处理
	// 			gotypeName = StructName(field2.Name)
	// 		}

	typ := graphqlType(req, col)

	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func goType2(req *CodeGenRequest, variable *ast.VariableDefinition) string {
	col := &Column{
		Name:    variable.Variable,
		IsArray: variable.Type.NamedType == "",
		Type: &Identifier{
			Name: variable.Type.NamedType,
		},
		NotNull: variable.Type.NonNull,
		Comment: "todo2",
	}
	if col.IsArray {
		col.Type = &Identifier{
			Name: variable.Type.Elem.NamedType,
		}
	}

	typ := graphqlType(req, col)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
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
		return "*" + columnType

	}
}
