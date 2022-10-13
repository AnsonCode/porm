package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AnsonCode/porm/codegen"

	"github.com/AnsonCode/porm/config"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 at the top of porm.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The only supported version is "1".
`

const errMessageNoPackages = `No packages are configured`

// func printFileErr(stderr io.Writer, dir string, fileErr *multierr.FileError) {
// 	filename := strings.TrimPrefix(fileErr.Filename, dir+"/")
// 	fmt.Fprintf(stderr, "%s:%d:%d: %s\n", filename, fileErr.Line, fileErr.Column, fileErr.Err)
// }

type outPair struct {
	Gen    config.SQLGen
	Plugin *config.Codegen

	config.SQL
}

func readConfig(stderr io.Writer, dir, filename string) (string, *config.Config, error) {
	configPath := ""
	if filename != "" {
		configPath = filepath.Join(dir, filename)
	} else {
		var yamlMissing, jsonMissing bool
		yamlPath := filepath.Join(dir, "porm.yaml")
		jsonPath := filepath.Join(dir, "porm.json")

		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			yamlMissing = true
		}
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			jsonMissing = true
		}

		if yamlMissing && jsonMissing {
			fmt.Fprintln(stderr, "error parsing configuration files. porm.yaml or porm.json: file does not exist")
			return "", nil, errors.New("config file missing")
		}

		if !yamlMissing && !jsonMissing {
			fmt.Fprintln(stderr, "error: both porm.json and porm.yaml files present")
			return "", nil, errors.New("porm.json and porm.yaml present")
		}

		configPath = yamlPath
		if yamlMissing {
			configPath = jsonPath
		}
	}

	base := filepath.Base(configPath)
	blob, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %s: file does not exist\n", base)
		return "", nil, err
	}

	conf, err := config.ParseConfig(bytes.NewReader(blob))
	if err != nil {
		switch err {
		case config.ErrMissingVersion:
			fmt.Fprintf(stderr, errMessageNoVersion)
		case config.ErrUnknownVersion:
			fmt.Fprintf(stderr, errMessageUnknownVersion)
		case config.ErrNoPackages:
			fmt.Fprintf(stderr, errMessageNoPackages)
		}
		fmt.Fprintf(stderr, "error parsing %s: %s\n", base, err)
		return "", nil, err
	}

	return configPath, &conf, nil
}

func Generate(ctx context.Context, e Env, dir, filename string, stderr io.Writer) (map[string]string, error) {
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return nil, err
	}

	base := filepath.Base(configPath)
	if err := config.Validate(conf); err != nil {
		fmt.Fprintf(stderr, "error validating %s: %s\n", base, err)
		return nil, err
	}

	output := map[string]string{}
	errored := false
	// var pairs []outPair
	for _, sql := range conf.SQL {
		// if sql.Gen.Go != nil {
		// 	pairs = append(pairs, outPair{
		// 		SQL: sql,
		// 		Gen: config.SQLGen{Go: sql.Gen.Go},
		// 	})
		// }

		// for i, _ := range sql.Codegen {
		// 	pairs = append(pairs, outPair{
		// 		SQL:    sql,
		// 		Plugin: &sql.Codegen[i],
		// 	})
		// }
		// combo := config.Combine(*conf, sql.SQL)
		// if sql.Plugin != nil {
		// 	combo.Codegen = *sql.Plugin
		// }
		schema, _ := parseFSSchema(os.DirFS(dir), sql.Schema)
		operations, _ := parseFS(os.DirFS(dir), sql.Queries)

		// result, failed := parse(ctx, e, name, dir, sql.SQL, combo, parseOpts, stderr)
		// if failed {
		// 	if packageRegion != nil {
		// 		packageRegion.End()
		// 	}
		// 	errored = true
		// 	break
		// }

		resp, err := codegen2(ctx, sql, schema, operations)
		// fmt.Println(operations)

		// out, resp, err := codegen(ctx, combo, sql, result)
		if err != nil {
			name := sql.Gen.Go.Package
			fmt.Fprintf(stderr, "# package %s\n", name)
			fmt.Fprintf(stderr, "error generating code: %s\n", err)
			errored = true

			continue
		}

		files := map[string]string{}
		for _, file := range resp.Files {
			files[file.Name] = string(file.Contents)
		}
		for n, source := range files {
			out := sql.Gen.Go.Out
			filename := filepath.Join(dir, out, n)
			output[filename] = source
		}
	}
	// for _, sql := range pairs {

	// }
	if errored {
		return nil, fmt.Errorf("errored")
	}
	return output, nil
}

func codegen2(ctx context.Context, sql config.SQL, schemaSource []*ast.Source, operationStr []byte) (*codegen.CodeGenResponse, error) {

	// sch, _ := os.ReadFile("../schema.graphql")
	// str, _ := os.ReadFile("../operations/create.graphql")
	// doc, _ := parser.ParseQuery(&ast.Source{Input: string(str)})
	// os.WriteFile("final.graphql", schemaStr, 0666)

	// schema := gqlparser.MustLoadSchema(&ast.Source{Input: string(schemaStr)})
	schema := gqlparser.MustLoadSchema(schemaSource...)

	// TODO:逐个加载operation文件，失败的忽略，合法校验，操作重命名

	doc, err := gqlparser.LoadQuery(schema, string(operationStr))
	if err != nil {
		panic(err)
	}

	// Create a new template and parse the letter into it.
	// opd := doc.Operations.ForName("")

	req := &codegen.CodeGenRequest{
		Schema:        schema,
		OperationList: doc.Operations,
		PormVersion:   "porm_V1.20",
		Settings: &codegen.Settings{
			Go: &codegen.GoCode{
				Package: sql.Gen.Go.Package,
			},
		},
	}
	// valid(req)
	return codegen.Generate2(ctx, req)
}

func parseFSSchema(fsys fs.FS, patterns []string) ([]*ast.Source, error) {
	// var filenames []string
	var result []*ast.Source
	for _, pattern := range patterns {
		list, err := fs.Glob(fsys, pattern)
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
		}
		// filenames = append(filenames, list...)
		for _, path := range list {
			tmpcontent, _ := ioutil.ReadFile(path)
			source := ast.Source{Input: string(tmpcontent)}
			result = append(result, &source)
		}

	}

	return result, nil
}

func parseFS(fsys fs.FS, patterns []string) ([]byte, error) {
	// var filenames []string
	var content []byte
	for _, pattern := range patterns {
		list, err := fs.Glob(fsys, pattern)
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
		}
		// filenames = append(filenames, list...)
		for _, path := range list {
			tmpcontent, _ := ioutil.ReadFile(path)
			content = append(content, tmpcontent...)
		}

	}

	return content, nil
}
