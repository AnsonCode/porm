package codegen

import "github.com/vektah/gqlparser/v2/ast"

// func addExtraGoStructTags(tags map[string]string, req *CodeGenRequest, col *Column) {
// 	for _, oride := range req.Settings.Overrides {
// 		if oride.GoType.StructTags == nil {
// 			continue
// 		}
// 		// if !sdk.Matches(oride, col.Table, req.Catalog.DefaultSchema) {
// 		// 	// Different table.
// 		// 	continue
// 		// }
// 		// if !sdk.MatchString(oride.ColumnName, col.Name) {
// 		// 	// Different column.
// 		// 	continue
// 		// }
// 		// Add the extra tags.
// 		for k, v := range oride.GoType.StructTags {
// 			tags[k] = v
// 		}
// 	}
// }

func goType3(req *CodeGenRequest, typName string, field *ast.FieldDefinition) string {
	col := &Column{
		Name:    field.Name,
		IsArray: field.Type.NamedType == "",
		Type: &Identifier{
			Name: typName,
		},
		NotNull: field.Type.NonNull,
		Comment: "todo2",
	}
	typ := graphqlType(req, col)

	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func goType(req *CodeGenRequest, field *ast.FieldDefinition) string {
	typName := field.Type.NamedType
	isArray := typName == ""
	// 说明是数组
	if isArray {
		typName = field.Type.Elem.NamedType
	}

	return goType3(req, typName, field)

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
		return "enum"

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
