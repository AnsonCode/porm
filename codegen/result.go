package codegen

import (
	"fmt"
	"sort"

	"github.com/vektah/gqlparser/v2/ast"
)

func buildStructsAndEnums(req *CodeGenRequest) ([]Struct, []Enum) {
	var structs []Struct

	var inputStructNameList []string

	for _, opd := range req.OperationList {
		// 收集输入参数的类型
		for _, v := range opd.VariableDefinitions {
			name := v.Definition.Name
			if v.Definition.Kind != "SCALAR" {
				// 这个地方追加时，也要判断下
				if !in(name, inputStructNameList) {
					inputStructNameList = append(inputStructNameList, name)
				}
				// TODO:优化递归方法
				recursiveFindInputStructName(name, req.Schema.Types, &inputStructNameList)
			}
		}

		// opd := doc.Operations.ForName("")

		// 收集响应参数的类型
		for _, v := range opd.SelectionSet {
			filed, ok := v.(*ast.Field)
			if ok {
				if opd.Name == "reate" {
					fmt.Println("debug")
				}
				// TODO: 递归的方法尽量不用全局变量
				recursiveResStruct(req, opd.Name, filed, &structs)
			}
		}

		// 这里补充一个 整体的响应
		s := Struct{
			IsRes: true,
			// Table:   Identifier{Schema: "", Name: opd.Name},
			Name:    StructName(opd.Name) + "Response", //
			Comment: "res2_struct",
			Fields:  selectionSet2Fields(req, opd.Name, opd.SelectionSet),
		}
		structs = append(structs, s)
		// 处理响应结果的结构体
		// structs = append(structs, structs...) // 把递归的加进来
		// structs = []Struct{}                  //清空下
	}

	var enums []Enum

	// 处理输入参数的结构体
	for _, structName := range inputStructNameList {
		d := req.Schema.Types[structName]
		fmt.Println(d.Name)
		def := req.Schema.Types[structName]
		if def.Kind == "ENUM" {
			fmt.Print("enum:", structName)
			e := Enum{
				Name:    StructName(structName),
				Comment: "enum",
			}
			for _, enumValue := range def.EnumValues {
				e.Constants = append(e.Constants, Constant{
					Name:  StructName(e.Name + "_" + enumValue.Name),
					Value: enumValue.Name,
					Type:  e.Name,
				})
			}
			enums = append(enums, e)

		} else {
			s := buildStruct(req, def)
			fmt.Println(s)

			structs = append(structs, *s)
		}
	}

	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}
	if len(enums) > 0 {
		sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })
	}
	return structs, enums
}

func buildStruct(req *CodeGenRequest, def *ast.Definition) *Struct {
	s := Struct{
		IsRes: false,
		Kind:  string(def.Kind),
		// Table:   plugin.Identifier{Schema: schema.Name, Name: table.Rel.Name},
		Name:    StructName(def.Name),
		Comment: string(def.Kind),
	}
	for _, field := range def.Fields {

		tags := map[string]string{}
		// if req.Settings.Go.EmitJsonTags {
		// 	// tags["json"] = JSONTagName(column.Name, req.Settings)
		// }
		tags["json"] = JSONTagName2(field.Name, true)
		// addExtraGoStructTags(tags, req, column)
		s.Fields = append(s.Fields, Field{
			Name:    StructName(field.Name), //StructName(column.Name)
			Type:    goType(req, field),
			Tags:    tags,
			Comment: "df",
		})
	}
	return &s
}

//  optName 为响应的结构体加上前缀，避免重复
func recursiveResStruct(req *CodeGenRequest, optName string, root *ast.Field, ret *[]Struct) {
	if root.SelectionSet == nil {
		return
	}
	for _, v := range root.SelectionSet {
		filed, ok := v.(*ast.Field)
		if ok {
			if filed.SelectionSet != nil {
				// 这里递归下
				recursiveResStruct(req, optName, filed, ret)
			}
		}
	}

	kind := string(root.ObjectDefinition.Kind) + "_" + root.ObjectDefinition.Name + "_" + "PART"
	if kind == "OBJECT_Mutation_PART" {
		fmt.Println("xx")
	}
	s := Struct{
		IsRes: true,
		Kind:  kind,
		// Table:   Identifier{Schema: "", Name: optName},
		Name:    StructName(optName) + StructName(root.Name), // 这里修改了对应查询响应的名字（但是应该怎么修改它的字段值的）
		Comment: kind,
	}
	s.Fields = selectionSet2Fields(req, optName, root.SelectionSet)

	*ret = append(*ret, s)
}

func selectionSet2Fields(req *CodeGenRequest, optName string, set ast.SelectionSet) []Field {
	fields := []Field{}
	for _, selection := range set {
		field2, ok := selection.(*ast.Field)
		if ok {
			tags := map[string]string{}
			tags["json"] = JSONTagName(field2.Alias, req.Settings)
			// addExtraGoStructTags(tags, req, column)
			isObject := field2.SelectionSet != nil
			gotypeName := goType(req, field2.Definition)
			fieldName := StructName(field2.Name)
			// 如果是对象，名字要处理下
			if isObject {
				typName := StructName(optName) + StructName(field2.Name)
				gotypeName = goType3(req, typName, field2.Definition)
				// gotypeName = "*" + StructName(optName) + StructName(field2.Name)
			}

			fields = append(fields, Field{
				Name: fieldName, //StructName(column.Name)
				// DBName:   optName,
				Type:     gotypeName,
				Tags:     tags,
				IsObject: isObject,
				// Struct: , // 应该查找出它的对象
				Comment: "-",
			})
		}
	}
	return fields
}

// graphql的标量
var defaultGqlDefinitionNamedType = []string{"String", "DateTime", "Boolean", "Float", "Int", "enum"}

// 怕死循环了~
func recursiveFindInputStructName(defname string, all map[string]*ast.Definition, ret *[]string) {
	def, ok := all[defname]
	if !ok {
		fmt.Println(defname)
		// return []string{}
		return
	}

	for _, v2 := range def.Fields {
		namedType := v2.Type.NamedType
		if namedType == "" {
			namedType = v2.Type.Elem.NamedType
		}
		// 默认的忽略
		if in(namedType, defaultGqlDefinitionNamedType) {
			continue
		}
		// // 已经加进去的忽略
		if in(namedType, *ret) {
			continue
		}
		// inputStructNameList = append(inputStructNameList, namedType)
		*ret = append(*ret, namedType)

		// 还要继续遍历（递归下）
		recursiveFindInputStructName(namedType, all, ret)
		// result = append(result, lastResult...)
		// res = append(res, childres...)
	}
	// return result
}

func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

// func diguiFind2(optName string, root *ast.Field) {
// 	if root.SelectionSet == nil {
// 		return
// 	}
// 	root.Name = optName + root.Name // 这里修改了对应查询响应的名字
// 	// root.Alias = optName + root.Alias
// 	fieldList = append(fieldList, root)

// 	for _, v := range root.SelectionSet {
// 		filed, ok := v.(*ast.Field)
// 		if ok {
// 			if filed.SelectionSet != nil {
// 				// fieldList = append(fieldList, filed)
// 				// 这里递归下
// 				diguiFind2(optName, filed)
// 			}
// 		}
// 		// var _ = (ast.Field)(&container{})
// 	}
// }

// type goColumn struct {
// 	id int
// 	// *plugin.Column
// }

// func columnName(c *plugin.Column, pos int) string {
// 	if c.Name != "" {
// 		return c.Name
// 	}
// 	return fmt.Sprintf("column_%d", pos+1)
// }

// func paramName(p *plugin.Parameter) string {
// 	if p.Column.Name != "" {
// 		return argName(p.Column.Name)
// 	}
// 	return fmt.Sprintf("dollar_%d", p.Number)
// }

// func argName(name string) string {
// 	out := ""
// 	for i, p := range strings.Split(name, "_") {
// 		if i == 0 {
// 			out += strings.ToLower(p)
// 		} else if p == "id" {
// 			out += "ID"
// 		} else {
// 			out += strings.Title(p)
// 		}
// 	}
// 	return out
// }

func parseDirective(operation *ast.OperationDefinition) (string, string) {
	query := FormatOperateionDocument(operation)

	for _, sel := range operation.SelectionSet {
		filed, ok := sel.(*ast.Field)
		if ok {
			for _, directive := range filed.Directives {
				if directive.Name == "sql" {
					query = directive.Arguments.ForName("raw").Value.Raw
					// ".*\\$\\{([a-z]+)\\}.*"
					return "RawSQL", query
				}
			}
		}
	}

	return "Do", query

}

func buildQueries(req *CodeGenRequest, structs []Struct) ([]Query, error) {
	qs := make([]Query, 0, len(req.OperationList))
	for _, operation := range req.OperationList {
		if operation.Name == "" {
			fmt.Println("query name is null ,skip")
			continue
		}
		clientMethod, query := parseDirective(operation)
		gq := Query{
			// Cmd:          query.Cmd,
			ConstantName: LowerTitle(operation.Name),
			// FieldName:    LowerTitle(query.Name) + "Stmt",
			MethodName:   Title(operation.Name),
			SourceName:   operation.Name, // TODO：这里要完善
			Comments:     []string{string(operation.Operation)},
			SQL:          query, // TODO:这里要重新读取
			ClientMethod: clientMethod,
			// Table:        query.InsertIntoTable,
		}
		// sqlpkg := SQLPackageFromString(req.Settings.Go.SqlPackage)

		allArg := []QueryValue{}
		for _, variable := range operation.VariableDefinitions {
			arg := QueryValue{
				Name: variable.Variable, //paramName(p)
				Typ:  goType2(req, variable),
				// SQLPackage: sqlpkg,
			}
			allArg = append(allArg, arg)
		}
		gq.Arg = allArg

		// 这里开始构造方法的返回结果
		gs := &Struct{
			Name: StructName(operation.Name) + "Response",
		}
		for _, selection := range operation.SelectionSet {
			field3, ok := selection.(*ast.Field)
			if ok {
				newfield := Field{
					Name: field3.Name,
					Type: goType(req, field3.Definition),
				}

				gs.Fields = append(gs.Fields, newfield)
			}
		}
		gq.Ret = QueryValue{
			Emit:   true,
			Name:   gs.Name,
			Struct: gs,
			// SQLPackage:  sqlpkg,
			EmitPointer: req.Settings.Go.EmitResultStructPointers,
		}

		qs = append(qs, gq)
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].MethodName < qs[j].MethodName })
	return qs, nil
}

// It's possible that this method will generate duplicate JSON tag values
//
//   Columns: count, count,   count_2
//    Fields: Count, Count_2, Count2
// JSON tags: count, count_2, count_2
//
// This is unlikely to happen, so don't fix it yet
// func columnsToStruct(req *plugin.CodeGenRequest, name string, columns []goColumn, useID bool) (*Struct, error) {
// 	gs := Struct{
// 		Name: name,
// 	}
// 	seen := map[string][]int{}
// 	suffixes := map[int]int{}
// 	for i, c := range columns {
// 		colName := columnName(c.Column, i)
// 		tagName := colName
// 		fieldName := StructName(colName, req.Settings)
// 		baseFieldName := fieldName
// 		// Track suffixes by the ID of the column, so that columns referring to the same numbered parameter can be
// 		// reused.
// 		suffix := 0
// 		if o, ok := suffixes[c.id]; ok && useID {
// 			suffix = o
// 		} else if v := len(seen[fieldName]); v > 0 && !c.IsNamedParam {
// 			suffix = v + 1
// 		}
// 		suffixes[c.id] = suffix
// 		if suffix > 0 {
// 			tagName = fmt.Sprintf("%s_%d", tagName, suffix)
// 			fieldName = fmt.Sprintf("%s_%d", fieldName, suffix)
// 		}
// 		tags := map[string]string{}
// 		if req.Settings.Go.EmitDbTags {
// 			tags["db"] = tagName
// 		}
// 		if req.Settings.Go.EmitJsonTags {
// 			tags["json"] = JSONTagName(tagName, req.Settings)
// 		}
// 		gs.Fields = append(gs.Fields, Field{
// 			Name:   fieldName,
// 			DBName: colName,
// 			Type:   goType(req, c.Column),
// 			Tags:   tags,
// 		})
// 		if _, found := seen[baseFieldName]; !found {
// 			seen[baseFieldName] = []int{i}
// 		} else {
// 			seen[baseFieldName] = append(seen[baseFieldName], i)
// 		}
// 	}

// 	// If a field does not have a known type, but another
// 	// field with the same name has a known type, assign
// 	// the known type to the field without a known type
// 	for i, field := range gs.Fields {
// 		if len(seen[field.Name]) > 1 && field.Type == "interface{}" {
// 			for _, j := range seen[field.Name] {
// 				if i == j {
// 					continue
// 				}
// 				otherField := gs.Fields[j]
// 				if otherField.Type != field.Type {
// 					field.Type = otherField.Type
// 				}
// 				gs.Fields[i] = field
// 			}
// 		}
// 	}

// 	err := checkIncompatibleFieldTypes(gs.Fields)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &gs, nil
// }

// func checkIncompatibleFieldTypes(fields []Field) error {
// 	fieldTypes := map[string]string{}
// 	for _, field := range fields {
// 		if fieldType, found := fieldTypes[field.Name]; !found {
// 			fieldTypes[field.Name] = field.Type
// 		} else if field.Type != fieldType {
// 			return fmt.Errorf("named param %s has incompatible types: %s, %s", field.Name, field.Type, fieldType)
// 		}
// 	}
// 	return nil
// }
