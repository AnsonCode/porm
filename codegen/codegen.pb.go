package codegen

import "github.com/vektah/gqlparser/v2/ast"

type File struct {
	Name     string `json:"name,omitempty"`
	Contents []byte `json:"contents,omitempty"`
}

type ParsedGoType struct {
	ImportPath string            `json:"import_path,omitempty"`
	Package    string            `json:"package,omitempty"`
	TypeName   string            `json:"type_name,omitempty"`
	BasicType  bool              `json:"basic_type,omitempty"`
	StructTags map[string]string `json:"struct_tags,omitempty"`
}

type Settings struct {
	Version string            `json:"version,omitempty"`
	Engine  string            `json:"engine,omitempty"`
	Schema  []string          `json:"schema,omitempty"`
	Queries []string          `json:"queries,omitempty"`
	Rename  map[string]string `json:"rename,omitempty"`
	// Overrides []*Override       `json:"overrides,omitempty"`
	Codegen *Codegen `json:"codegen,omitempty"`
	// TODO: Refactor codegen settings
	Go *GoCode `json:"go,omitempty"`
}

type Codegen struct {
	Out     string `json:"out,omitempty"`
	Plugin  string `json:"plugin,omitempty"`
	Options []byte `json:"options,omitempty"`
}

type GoCode struct {
	EmitInterface             bool   `json:"emit_interface,omitempty"`
	EmitJsonTags              bool   `json:"emit_json_tags,omitempty"`
	EmitDbTags                bool   `json:"emit_db_tags,omitempty"`
	EmitPreparedQueries       bool   `json:"emit_prepared_queries,omitempty"`
	EmitExactTableNames       bool   `json:"emit_exact_table_names,omitempty"`
	EmitEmptySlices           bool   `json:"emit_empty_slices,omitempty"`
	EmitExportedQueries       bool   `json:"emit_exported_queries,omitempty"`
	EmitResultStructPointers  bool   `json:"emit_result_struct_pointers,omitempty"`
	EmitParamsStructPointers  bool   `json:"emit_params_struct_pointers,omitempty"`
	EmitMethodsWithDbArgument bool   `json:"emit_methods_with_db_argument,omitempty"`
	JsonTagsCaseStyle         string `json:"json_tags_case_style,omitempty"`
	Package                   string `json:"package,omitempty"`
	Out                       string `json:"out,omitempty"`
	SqlPackage                string `json:"sql_package,omitempty"`
	OutputDbFileName          string `json:"output_db_file_name,omitempty"`
	OutputModelsFileName      string `json:"output_models_file_name,omitempty"`
	OutputQuerierFileName     string `json:"output_querier_file_name,omitempty"`
	OutputFilesSuffix         string `json:"output_files_suffix,omitempty"`
	EmitEnumValidMethod       bool   `json:"emit_enum_valid_method,omitempty"`
	EmitAllEnumValues         bool   `json:"emit_all_enum_values,omitempty"`
}

// type Catalog struct {
// 	Comment       string    `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
// 	DefaultSchema string    `protobuf:"bytes,2,opt,name=default_schema,json=defaultSchema,proto3" json:"default_schema,omitempty"`
// 	Name          string    `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
// 	Schemas       []*Schema `protobuf:"bytes,4,rep,name=schemas,proto3" json:"schemas,omitempty"`
// }

// type Schema struct {
// 	Comment string `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
// 	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
// 	// Tables         []*Table         `protobuf:"bytes,3,rep,name=tables,proto3" json:"tables,omitempty"`
// 	Enums          []*Enum          `protobuf:"bytes,4,rep,name=enums,proto3" json:"enums,omitempty"`
// 	CompositeTypes []*CompositeType `protobuf:"bytes,5,rep,name=composite_types,json=compositeTypes,proto3" json:"composite_types,omitempty"`
// }

// type CompositeType struct {
// 	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
// 	Comment string `protobuf:"bytes,2,opt,name=comment,proto3" json:"comment,omitempty"`
// }

// type Enum struct {
// 	Name    string   `json:"name,omitempty"`
// 	Vals    []string `json:"vals,omitempty"`
// 	Comment string   `json:"comment,omitempty"`
// }

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

type CodeGenRequest struct {
	Settings      *Settings `json:"settings,omitempty"`
	Schema        *ast.Schema
	OperationList []*ast.OperationDefinition

	// Catalog     *Catalog `json:"catalog,omitempty"`
	// Queries     []*Query `json:"queries,omitempty"`
	PormVersion string `json:"porm_version,omitempty"`
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
