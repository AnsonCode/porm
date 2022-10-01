package codegen

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/vektah/gqlparser/v2"
)

var defaultGqlDefinitionNamedType = []string{"StringFilter", "StringNullableFilter", "DateTimeFilter", "BoolFilter"}

func Read() {

	// Define a template.
	sch, _ := os.ReadFile("../schema.graphql")
	str, _ := os.ReadFile("../operations/query.graphql")
	// doc, _ := parser.ParseQuery(&ast.Source{Input: string(str)})
	schema := gqlparser.MustLoadSchema(&ast.Source{Input: string(sch)})

	doc, _ := gqlparser.LoadQuery(schema, string(str))

	// Create a new template and parse the letter into it.
	// opd := doc.Operations.ForName("")

	req := &CodeGenRequest{
		Schema:        schema,
		OperationList: doc.Operations,
		SqlcVersion:   "test",
		Settings: &Settings{
			Go: &GoCode{
				Package: "generate2",
			},
		},
	}
	ctx := context.TODO()
	Generate2(ctx, req)
}
func Generate2(ctx context.Context, req *CodeGenRequest) (*CodeGenResponse, error) {
	// enums := buildEnums(req)
	structs := buildStructs(req)
	fmt.Println(structs)

	queries, err := buildQueries(req, structs)
	if err != nil {
		return nil, err
	}
	fmt.Println(queries)

	return generate(req, structs, queries)
}

type tmplCtx struct {
	Q             string
	SourceName    string
	Package       string
	GqlSchemaPath string
	OperationDir  string

	Operations []*ast.Operation

	Enums       []Enum
	Structs     []Struct
	GoQueries   []Query
	PormVersion string
}

func (t *tmplCtx) OutputQuery(sourceName string) bool {
	return t.SourceName == sourceName
}

// enums []Enum,
func generate(req *CodeGenRequest, structs []Struct, queries []Query) (*CodeGenResponse, error) {
	// i := &importer{
	// 	Settings: req.Settings,
	// 	Queries:  queries,
	// 	Enums:    enums,
	// 	Structs:  structs,
	// }

	funcMap := template.FuncMap{
		"lowerTitle": LowerTitle,
		"comment":    DoubleSlashComment,
		"escape":     EscapeBacktick,
		// "imports":    i.Imports,
		"hasPrefix": strings.HasPrefix,
	}

	tmpl := template.Must(
		template.New("table").
			Funcs(funcMap).
			ParseFS(
				templates,
				"templates/*.tmpl",
				// "templates/*/*.tmpl",
			),
	)

	golang := req.Settings.Go
	tctx := tmplCtx{
		Q:           "`",
		Package:     golang.Package,
		GoQueries:   queries,
		Enums:       []Enum{},
		Structs:     structs,
		PormVersion: req.SqlcVersion,
	}

	output := map[string]string{}

	execute := func(name, templateName string) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := tmpl.ExecuteTemplate(w, templateName, &tctx)
		w.Flush()
		if err != nil {
			return err
		}
		code, err := format.Source(b.Bytes())
		if err != nil {
			fmt.Println(b.String())
			return fmt.Errorf("source error: %w", err)
		}

		if templateName == "queryFile" && golang.OutputFilesSuffix != "" {
			name += golang.OutputFilesSuffix
		}

		if !strings.HasSuffix(name, ".go") {
			name += ".go"
		}
		output[name] = string(code)
		return nil
	}

	// dbFileName := "db.go"
	// if golang.OutputDbFileName != "" {
	// 	dbFileName = golang.OutputDbFileName
	// }
	modelsFileName := "models.go"
	// if golang.OutputModelsFileName != "" {
	// 	modelsFileName = golang.OutputModelsFileName
	// }
	// querierFileName := "querier.go"
	// if golang.OutputQuerierFileName != "" {
	// 	querierFileName = golang.OutputQuerierFileName
	// }

	// if err := execute(dbFileName, "dbFile"); err != nil {
	// 	return nil, err
	// }
	if err := execute(modelsFileName, "modelsFile"); err != nil {
		return nil, err
	}
	// if golang.EmitInterface {
	// 	if err := execute(querierFileName, "interfaceFile"); err != nil {
	// 		return nil, err
	// 	}
	// }

	files := map[string]struct{}{}
	for _, gq := range queries {
		files[gq.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, "queryFile"); err != nil {
			return nil, err
		}
	}
	resp := CodeGenResponse{}

	for filename, code := range output {
		resp.Files = append(resp.Files, &File{
			Name:     filename,
			Contents: []byte(code),
		})
		if err := ioutil.WriteFile("./generate2/"+filename+".go", []byte(code), 0666); err != nil {
			panic(err)
		}
	}

	return &resp, nil
}

func Generate() {
	// Define a template.
	sch, _ := os.ReadFile("../schema.graphql")
	str, _ := os.ReadFile("../operations/query.graphql")
	// doc, _ := parser.ParseQuery(&ast.Source{Input: string(str)})
	schema := gqlparser.MustLoadSchema(&ast.Source{Input: string(sch)})

	doc, _ := gqlparser.LoadQuery(schema, string(str))

	// Create a new template and parse the letter into it.
	opd := doc.Operations.ForName("")

	// opd.Operation
	for _, v := range opd.SelectionSet {
		filed, ok := v.(*ast.Field)
		if ok {
			diguiFind2("T", filed)
		}
	}
	for _, filed := range fieldList {
		fmt.Println(filed.Alias)
		for _, v := range filed.SelectionSet {
			field2, ok := v.(*ast.Field)
			if ok {
				tpy := field2.Definition.Type.NamedType
				if tpy == "" {
					fmt.Println(field2.Alias, "[]", field2.Definition.Type.Elem.NamedType)
				} else {
					fmt.Println(field2.Alias, "*", tpy)
				}
			}
		}
	}

	for _, v := range opd.VariableDefinitions {
		name := v.Definition.Name
		if v.Definition.Kind != "SCALAR" {
			res = append(res, name)
			diguiFind(name, schema.Types)
		}
	}
	fmt.Println(res)
	gqlTypes := map[string]*ast.Definition{} // 保存所有需要的变量

	// todo:合成
	for _, v := range res {
		d := schema.Types[v]
		fmt.Println(d.Name)
		gqlTypes[v] = schema.Types[v]
	}

	fmt.Println(opd)
	funcMap := template.FuncMap{
		"lowerTitle": LowerTitle,
		"Title":      Title,
		"comment":    DoubleSlashComment,
		"escape":     EscapeBacktick,
	}

	t := template.Must(template.New("letter").Funcs(funcMap).ParseFS(templates, "templates/*.tmpl"))

	Write(t, "operation", opd)
	Write(t, "operation_out_type", fieldList)

	Write(t, "custom_type", gqlTypes)
}

func Write(t *template.Template, name string, data any) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := t.ExecuteTemplate(w, name+".tmpl", data); err != nil {
		panic(err)
	}
	w.Flush()
	code, _ := format.Source(b.Bytes())

	if err := ioutil.WriteFile("./generate/"+name+".go", code, 0666); err != nil {
		panic(err)
	}
}

// 死循环了~
var res = []string{}

func diguiFind(defname string, all map[string]*ast.Definition) {
	def, ok := all[defname]
	if !ok {
		fmt.Println(defname)
		return
	}
	for _, v2 := range def.Fields {
		if v2.Type.NamedType == "" {
			continue
		}
		// 默认的忽略
		if in(v2.Type.NamedType, defaultGqlDefinitionNamedType) {
			continue
		}
		// 已经加进去的忽略
		if in(v2.Type.NamedType, res) {
			continue
		}
		res = append(res, v2.Type.NamedType)
		// 还要继续遍历（递归下）
		diguiFind(v2.Type.NamedType, all)
		// res = append(res, childres...)
	}
}

var fieldList = []*ast.Field{}

func diguiFind2(optName string, root *ast.Field) {
	if root.SelectionSet == nil {
		return
	}
	root.Name = optName + root.Name
	// root.Alias = optName + root.Alias
	fieldList = append(fieldList, root)

	for _, v := range root.SelectionSet {
		filed, ok := v.(*ast.Field)
		if ok {
			if filed.SelectionSet != nil {
				// fieldList = append(fieldList, filed)
				// 这里递归下
				diguiFind2(optName, filed)
			}
		}
		// var _ = (ast.Field)(&container{})
	}
}
func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}
