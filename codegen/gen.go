package codegen

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/vektah/gqlparser/v2"
)

func Read() {

	// Define a template.
	sch, _ := os.ReadFile("../schema.graphql")
	str, _ := os.ReadFile("../operations/create.graphql")
	// doc, _ := parser.ParseQuery(&ast.Source{Input: string(str)})
	schema := gqlparser.MustLoadSchema(&ast.Source{Input: string(sch)})

	doc, _ := gqlparser.LoadQuery(schema, string(str))

	// Create a new template and parse the letter into it.
	// opd := doc.Operations.ForName("")

	req := &CodeGenRequest{
		Schema:        schema,
		OperationList: doc.Operations,
		PormVersion:   "porm_V1.20",
		Settings: &Settings{
			Go: &GoCode{
				Package: "generate2",
			},
		},
	}
	valid(req)
	ctx := context.TODO()
	Generate2(ctx, req)
}

func valid(req *CodeGenRequest) {
	chars := req.OperationList

	for i := 0; i < len(chars); i++ {
		if req.OperationList[i].Name == "" {
			fmt.Println("操作必须命名，未命名操作已忽略")
			chars = append(chars[:i], chars[i+1:]...)
			i-- // form the remove item index to start iterate next item
		}
	}
	req.OperationList = chars
	fmt.Printf("%+v num:%v", len(chars), len(req.OperationList))
}

func Generate2(ctx context.Context, req *CodeGenRequest) (*CodeGenResponse, error) {
	// enums := buildEnums(req)
	structs, enums := buildStructs(req)
	fmt.Println(structs)

	queries, err := buildQueries(req, structs)
	if err != nil {
		return nil, err
	}
	fmt.Println(queries)

	return generate(req, enums, structs, queries)
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

func generate(req *CodeGenRequest, enums []Enum, structs []Struct, queries []Query) (*CodeGenResponse, error) {
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
		Enums:       enums,
		Structs:     structs,
		PormVersion: req.PormVersion,
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
	// if golang.OutputModelsFileName != "" {
	// 	modelsFileName = golang.OutputModelsFileName
	// }
	// querierFileName := "querier.go"
	// if golang.OutputQuerierFileName != "" {
	// 	querierFileName = golang.OutputQuerierFileName
	// }

	// if err := execute(".gitignore", "gitignore.tmpl"); err != nil {
	// 	return nil, err
	// }
	if err := execute("engine", "engine.tmpl"); err != nil {
		return nil, err
	}
	// if err := execute(dbFileName, "dbFile"); err != nil {
	// 	return nil, err
	// }
	modelsFileName := "models.go"
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
	}

	return &resp, nil
}
