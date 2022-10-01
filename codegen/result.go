package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// func buildEnums(req *plugin.CodeGenRequest) []Enum {
// 	var enums []Enum
// 	for _, schema := range req.Catalog.Schemas {
// 		if schema.Name == "pg_catalog" {
// 			continue
// 		}
// 		for _, enum := range schema.Enums {
// 			var enumName string
// 			if schema.Name == req.Catalog.DefaultSchema {
// 				enumName = enum.Name
// 			} else {
// 				enumName = schema.Name + "_" + enum.Name
// 			}
// 			e := Enum{
// 				Name:    StructName(enumName, req.Settings),
// 				Comment: enum.Comment,
// 			}
// 			seen := make(map[string]struct{}, len(enum.Vals))
// 			for i, v := range enum.Vals {
// 				value := EnumReplace(v)
// 				if _, found := seen[value]; found || value == "" {
// 					value = fmt.Sprintf("value_%d", i)
// 				}
// 				e.Constants = append(e.Constants, Constant{
// 					Name:  StructName(enumName+"_"+value, req.Settings),
// 					Value: v,
// 					Type:  e.Name,
// 				})
// 				seen[value] = struct{}{}
// 			}
// 			enums = append(enums, e)
// 		}
// 	}
// 	if len(enums) > 0 {
// 		sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })
// 	}
// 	return enums
// }

func buildStructs(req *CodeGenRequest) []Struct {
	var structs []Struct

	gqlTypes := map[string]*ast.Definition{} // 保存所有需要的变量

	for _, opd := range req.OperationList {
		// 收集输入参数的类型
		for _, v := range opd.VariableDefinitions {
			name := v.Definition.Name
			if v.Definition.Kind != "SCALAR" {
				res = append(res, name)
				diguiFind(name, req.Schema.Types)
			}
		}

		// opd := doc.Operations.ForName("")

		// opd.Operation
		// 收集响应参数的类型
		for _, v := range opd.SelectionSet {
			filed, ok := v.(*ast.Field)
			if ok {
				diguiFind2(filed)
			}
		}

	}

	// todo:res 要复原？
	for _, v := range res {
		d := req.Schema.Types[v]
		fmt.Println(d.Name)
		gqlTypes[v] = req.Schema.Types[v]
	}

	for structName, def := range gqlTypes {
		s := Struct{
			// Table:   plugin.Identifier{Schema: schema.Name, Name: table.Rel.Name},
			Name:    StructName(structName),
			Comment: "input_struct",
		}
		for _, field := range def.Fields {

			tags := map[string]string{}
			// if req.Settings.Go.EmitJsonTags {
			// 	// tags["json"] = JSONTagName(column.Name, req.Settings)
			// }
			// addExtraGoStructTags(tags, req, column)
			s.Fields = append(s.Fields, Field{
				Name:    field.Name, //StructName(column.Name)
				Type:    goType(req, field),
				Tags:    tags,
				Comment: "",
			})
		}
		fmt.Println(s)

		structs = append(structs, s)
	}

	for _, filed := range fieldList {
		fmt.Println(filed.Alias)
		s := Struct{
			// Table:   plugin.Identifier{Schema: schema.Name, Name: table.Rel.Name},
			Name:    filed.Alias, // StructName(filed.Alias)
			Comment: "res_struct",
		}
		for _, selection := range filed.SelectionSet {
			field2, ok := selection.(*ast.Field)
			if ok {
				tags := map[string]string{}
				// addExtraGoStructTags(tags, req, column)
				s.Fields = append(s.Fields, Field{
					Name:    field2.Name, //StructName(column.Name)
					Type:    goType(req, field2.Definition),
					Tags:    tags,
					Comment: "-",
				})
			}
		}
		structs = append(structs, s)
		fmt.Println(s)
	}

	// for _, schema := range req.Catalog.Schemas {
	// 	for _, table := range schema.Tables {
	// 		var tableName string
	// 		if schema.Name == req.Catalog.DefaultSchema {
	// 			tableName = table.Rel.Name
	// 		} else {
	// 			tableName = schema.Name + "_" + table.Rel.Name
	// 		}
	// 		structName := tableName
	// 		// if !req.Settings.Go.EmitExactTableNames {
	// 		// 	// structName = inflection.Singular(structName)
	// 		// }

	// 	}
	// }

	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}
	return structs
}

// type goColumn struct {
// 	id int
// 	*plugin.Column
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

func argName(name string) string {
	out := ""
	for i, p := range strings.Split(name, "_") {
		if i == 0 {
			out += strings.ToLower(p)
		} else if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

func buildQueries(req *CodeGenRequest, structs []Struct) ([]Query, error) {
	qs := make([]Query, 0, len(req.OperationList))
	for _, query := range req.OperationList {
		// if query.Name == "" {
		// 	continue
		// }
		// if query.Cmd == "" {
		// 	continue
		// }

		var constantName string
		// if req.Settings.Go.EmitExportedQueries {
		// 	constantName = Title(query.Name)
		// } else {
		// 	constantName = LowerTitle(query.Name)
		// }

		gq := Query{
			// Cmd:          query.Cmd,
			ConstantName: constantName,
			FieldName:    LowerTitle(query.Name) + "Stmt",
			MethodName:   query.Name,
			// SourceName:   query.Filename,
			// SQL:          query.Text,
			// Comments:     query.Comments,
			// Table:        query.InsertIntoTable,
		}
		// sqlpkg := SQLPackageFromString(req.Settings.Go.SqlPackage)

		if len(query.VariableDefinitions) == 1 {
			p := query.VariableDefinitions[0]
			gq.Arg = QueryValue{
				Name: p.Definition.Name, //paramName(p)
				// Typ:  goType(req, p.Column),
				// SQLPackage: sqlpkg,
			}
		}
		// else if len(query.VariableDefinitions) > 1 {
		// 	var cols []goColumn
		// 	for _, p := range query.VariableDefinitions {
		// 		cols = append(cols, goColumn{
		// 			id:     int(p.Number),
		// 			Column: p.Column,
		// 		})
		// 	}
		// 	s, err := columnsToStruct(req, gq.MethodName+"Params", cols, false)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	gq.Arg = QueryValue{
		// 		Emit:   true,
		// 		Name:   "arg",
		// 		Struct: s,
		// 		// SQLPackage:  sqlpkg,
		// 		// EmitPointer: req.Settings.Go.EmitParamsStructPointers,
		// 	}
		// }

		// if len(query.Columns) == 1 {
		// 	c := query.Columns[0]
		// 	name := columnName(c, 0)
		// 	if c.IsFuncCall {
		// 		name = strings.Replace(name, "$", "_", -1)
		// 	}
		// 	gq.Ret = QueryValue{
		// 		Name:       name,
		// 		Typ:        goType(req, c),
		// 		SQLPackage: sqlpkg,
		// 	}
		// } else if len(query.Columns) > 1 {
		// 	var gs *Struct
		// 	var emit bool

		// 	for _, s := range structs {
		// 		if len(s.Fields) != len(query.Columns) {
		// 			continue
		// 		}
		// 		same := true
		// 		for i, f := range s.Fields {
		// 			c := query.Columns[i]
		// 			sameName := f.Name == StructName(columnName(c, i), req.Settings)
		// 			sameType := f.Type == goType(req, c)
		// 			sameTable := sdk.SameTableName(c.Table, &s.Table, req.Catalog.DefaultSchema)
		// 			if !sameName || !sameType || !sameTable {
		// 				same = false
		// 			}
		// 		}
		// 		if same {
		// 			gs = &s
		// 			break
		// 		}
		// 	}

		// 	if gs == nil {
		// 		var columns []goColumn
		// 		for i, c := range query.Columns {
		// 			columns = append(columns, goColumn{
		// 				id:     i,
		// 				Column: c,
		// 			})
		// 		}
		// 		var err error
		// 		gs, err = columnsToStruct(req, gq.MethodName+"Row", columns, true)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		emit = true
		// 	}
		// 	gq.Ret = QueryValue{
		// 		Emit:        emit,
		// 		Name:        "i",
		// 		Struct:      gs,
		// 		SQLPackage:  sqlpkg,
		// 		EmitPointer: req.Settings.Go.EmitResultStructPointers,
		// 	}
		// }

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

func checkIncompatibleFieldTypes(fields []Field) error {
	fieldTypes := map[string]string{}
	for _, field := range fields {
		if fieldType, found := fieldTypes[field.Name]; !found {
			fieldTypes[field.Name] = field.Type
		} else if field.Type != fieldType {
			return fmt.Errorf("named param %s has incompatible types: %s, %s", field.Name, field.Type, fieldType)
		}
	}
	return nil
}
