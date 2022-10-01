package codegen

import "github.com/vektah/gqlparser/v2/ast"

type File struct {
	Name     string `json:"name,omitempty"`
	Contents []byte `json:"contents,omitempty"`
}

type Override struct {
	// name of the type to use, e.g. `github.com/segmentio/ksuid.KSUID` or `mymodule.Type`
	CodeType string `protobuf:"bytes,1,opt,name=code_type,proto3" json:"code_type,omitempty"`
	// name of the type to use, e.g. `text`
	DbType string `protobuf:"bytes,3,opt,name=db_type,proto3" json:"db_type,omitempty"`
	// True if the override should apply to a nullable database type
	Nullable bool `protobuf:"varint,5,opt,name=nullable,proto3" json:"nullable,omitempty"`
	// fully qualified name of the column, e.g. `accounts.id`
	Column     string        `protobuf:"bytes,6,opt,name=column,proto3" json:"column,omitempty"`
	Table      *Identifier   `protobuf:"bytes,7,opt,name=table,proto3" json:"table,omitempty"`
	ColumnName string        `protobuf:"bytes,8,opt,name=column_name,proto3" json:"column_name,omitempty"`
	GoType     *ParsedGoType `protobuf:"bytes,10,opt,name=go_type,json=goType,proto3" json:"go_type,omitempty"`
}

type ParsedGoType struct {
	ImportPath string            `protobuf:"bytes,1,opt,name=import_path,json=importPath,proto3" json:"import_path,omitempty"`
	Package    string            `protobuf:"bytes,2,opt,name=package,proto3" json:"package,omitempty"`
	TypeName   string            `protobuf:"bytes,3,opt,name=type_name,json=typeName,proto3" json:"type_name,omitempty"`
	BasicType  bool              `protobuf:"varint,4,opt,name=basic_type,json=basicType,proto3" json:"basic_type,omitempty"`
	StructTags map[string]string `protobuf:"bytes,5,rep,name=struct_tags,json=structTags,proto3" json:"struct_tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type Settings struct {
	Version   string            `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Engine    string            `protobuf:"bytes,2,opt,name=engine,proto3" json:"engine,omitempty"`
	Schema    []string          `protobuf:"bytes,3,rep,name=schema,proto3" json:"schema,omitempty"`
	Queries   []string          `protobuf:"bytes,4,rep,name=queries,proto3" json:"queries,omitempty"`
	Rename    map[string]string `protobuf:"bytes,5,rep,name=rename,proto3" json:"rename,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Overrides []*Override       `protobuf:"bytes,6,rep,name=overrides,proto3" json:"overrides,omitempty"`
	Codegen   *Codegen          `protobuf:"bytes,12,opt,name=codegen,proto3" json:"codegen,omitempty"`
	// TODO: Refactor codegen settings
	Go *GoCode `protobuf:"bytes,10,opt,name=go,proto3" json:"go,omitempty"`
}

type Codegen struct {
	Out     string `protobuf:"bytes,1,opt,name=out,proto3" json:"out,omitempty"`
	Plugin  string `protobuf:"bytes,2,opt,name=plugin,proto3" json:"plugin,omitempty"`
	Options []byte `protobuf:"bytes,3,opt,name=options,proto3" json:"options,omitempty"`
}

type GoCode struct {
	EmitInterface             bool   `protobuf:"varint,1,opt,name=emit_interface,json=emitInterface,proto3" json:"emit_interface,omitempty"`
	EmitJsonTags              bool   `protobuf:"varint,2,opt,name=emit_json_tags,json=emitJsonTags,proto3" json:"emit_json_tags,omitempty"`
	EmitDbTags                bool   `protobuf:"varint,3,opt,name=emit_db_tags,json=emitDbTags,proto3" json:"emit_db_tags,omitempty"`
	EmitPreparedQueries       bool   `protobuf:"varint,4,opt,name=emit_prepared_queries,json=emitPreparedQueries,proto3" json:"emit_prepared_queries,omitempty"`
	EmitExactTableNames       bool   `protobuf:"varint,5,opt,name=emit_exact_table_names,json=emitExactTableNames,proto3" json:"emit_exact_table_names,omitempty"`
	EmitEmptySlices           bool   `protobuf:"varint,6,opt,name=emit_empty_slices,json=emitEmptySlices,proto3" json:"emit_empty_slices,omitempty"`
	EmitExportedQueries       bool   `protobuf:"varint,7,opt,name=emit_exported_queries,json=emitExportedQueries,proto3" json:"emit_exported_queries,omitempty"`
	EmitResultStructPointers  bool   `protobuf:"varint,8,opt,name=emit_result_struct_pointers,json=emitResultStructPointers,proto3" json:"emit_result_struct_pointers,omitempty"`
	EmitParamsStructPointers  bool   `protobuf:"varint,9,opt,name=emit_params_struct_pointers,json=emitParamsStructPointers,proto3" json:"emit_params_struct_pointers,omitempty"`
	EmitMethodsWithDbArgument bool   `protobuf:"varint,10,opt,name=emit_methods_with_db_argument,json=emitMethodsWithDbArgument,proto3" json:"emit_methods_with_db_argument,omitempty"`
	JsonTagsCaseStyle         string `protobuf:"bytes,11,opt,name=json_tags_case_style,json=jsonTagsCaseStyle,proto3" json:"json_tags_case_style,omitempty"`
	Package                   string `protobuf:"bytes,12,opt,name=package,proto3" json:"package,omitempty"`
	Out                       string `protobuf:"bytes,13,opt,name=out,proto3" json:"out,omitempty"`
	SqlPackage                string `protobuf:"bytes,14,opt,name=sql_package,json=sqlPackage,proto3" json:"sql_package,omitempty"`
	OutputDbFileName          string `protobuf:"bytes,15,opt,name=output_db_file_name,json=outputDbFileName,proto3" json:"output_db_file_name,omitempty"`
	OutputModelsFileName      string `protobuf:"bytes,16,opt,name=output_models_file_name,json=outputModelsFileName,proto3" json:"output_models_file_name,omitempty"`
	OutputQuerierFileName     string `protobuf:"bytes,17,opt,name=output_querier_file_name,json=outputQuerierFileName,proto3" json:"output_querier_file_name,omitempty"`
	OutputFilesSuffix         string `protobuf:"bytes,18,opt,name=output_files_suffix,json=outputFilesSuffix,proto3" json:"output_files_suffix,omitempty"`
	EmitEnumValidMethod       bool   `protobuf:"varint,19,opt,name=emit_enum_valid_method,json=emitEnumValidMethod,proto3" json:"emit_enum_valid_method,omitempty"`
	EmitAllEnumValues         bool   `protobuf:"varint,20,opt,name=emit_all_enum_values,json=emitAllEnumValues,proto3" json:"emit_all_enum_values,omitempty"`
}

type Catalog struct {
	Comment       string    `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
	DefaultSchema string    `protobuf:"bytes,2,opt,name=default_schema,json=defaultSchema,proto3" json:"default_schema,omitempty"`
	Name          string    `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Schemas       []*Schema `protobuf:"bytes,4,rep,name=schemas,proto3" json:"schemas,omitempty"`
}

type Schema struct {
	Comment        string           `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
	Name           string           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Tables         []*Table         `protobuf:"bytes,3,rep,name=tables,proto3" json:"tables,omitempty"`
	Enums          []*Enum          `protobuf:"bytes,4,rep,name=enums,proto3" json:"enums,omitempty"`
	CompositeTypes []*CompositeType `protobuf:"bytes,5,rep,name=composite_types,json=compositeTypes,proto3" json:"composite_types,omitempty"`
}

type CompositeType struct {
	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Comment string `protobuf:"bytes,2,opt,name=comment,proto3" json:"comment,omitempty"`
}

type Enum struct {
	Name    string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Vals    []string `protobuf:"bytes,2,rep,name=vals,proto3" json:"vals,omitempty"`
	Comment string   `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty"`
}

type Table struct {
	Rel     *Identifier `protobuf:"bytes,1,opt,name=rel,proto3" json:"rel,omitempty"`
	Columns []*Column   `protobuf:"bytes,2,rep,name=columns,proto3" json:"columns,omitempty"`
	Comment string      `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty"`
}

type Identifier struct {
	Catalog string `json:"catalog,omitempty"`
	Schema  string `json:"schema,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Column struct {
	Name         string `json:"name,omitempty"`
	NotNull      bool   `json:"not_null,omitempty"`
	IsArray      bool   `json:"is_array,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Length       int32  `json:"length,omitempty"`
	IsNamedParam bool   `json:"is_named_param,omitempty"`
	IsFuncCall   bool   `json:"is_func_call,omitempty"`
	// XXX: Figure out what PostgreSQL calls `foo.id`
	Scope      string      `json:"scope,omitempty"`
	Table      *Identifier `json:"table,omitempty"`
	TableAlias string      `json:"table_alias,omitempty"`
	Type       *Identifier `json:"type,omitempty"`
}

// type Query struct {
// 	Text            string       `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
// 	Name            string       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
// 	Cmd             string       `protobuf:"bytes,3,opt,name=cmd,proto3" json:"cmd,omitempty"`
// 	Columns         []*Column    `protobuf:"bytes,4,rep,name=columns,proto3" json:"columns,omitempty"`
// 	Params          []*Parameter `protobuf:"bytes,5,rep,name=params,json=parameters,proto3" json:"params,omitempty"`
// 	Comments        []string     `protobuf:"bytes,6,rep,name=comments,proto3" json:"comments,omitempty"`
// 	Filename        string       `protobuf:"bytes,7,opt,name=filename,proto3" json:"filename,omitempty"`
// 	InsertIntoTable *Identifier  `protobuf:"bytes,8,opt,name=insert_into_table,proto3" json:"insert_into_table,omitempty"`
// }

type Parameter struct {
	Number int32   `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Column *Column `protobuf:"bytes,2,opt,name=column,proto3" json:"column,omitempty"`
}

// sch, _ := os.ReadFile("../schema.graphql")
// 	str, _ := os.ReadFile("../operations/query.graphql")
// 	// doc, _ := parser.ParseQuery(&ast.Source{Input: string(str)})
// 	schema := gqlparser.MustLoadSchema(&ast.Source{Input: string(sch)})

// 	doc, _ := gqlparser.LoadQuery(schema, string(str))
type CodeGenRequest struct {
	Settings      *Settings `json:"settings,omitempty"`
	Schema        *ast.Schema
	OperationList []*ast.OperationDefinition

	// Catalog     *Catalog `json:"catalog,omitempty"`
	// Queries     []*Query `json:"queries,omitempty"`
	SqlcVersion string `json:"sqlc_version,omitempty"`
}

type CodeGenResponse struct {
	Files []*File `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *CodeGenResponse) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}
