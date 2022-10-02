package main

import (
	"context"
	"fmt"
	"os"

	"github.com/AnsonCode/porm/codegen/generate2"
	"github.com/prisma/prisma-client-go/engine"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

var GLOBAL_SCHEME string

func init() {
	GLOBAL_SCHEME = mustReadFile("./schema.prisma")
}
func mustReadFile(name string) string {
	src, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return string(src)
}

func main() {
	_, schema := testSdl()
	operationStr := mustReadFile("./operations/query3.graphql")

	doc, err := gqlparser.LoadQuery(schema, operationStr)
	fmt.Println(doc, err)

	if len(doc.Operations) != 1 {
		panic("必须只有一个查询")
	}

	query := generate2.NewClient2(GLOBAL_SCHEME)
	query.Connect()
	defer query.Disconnect()

	id := "1"
	whe := &generate2.PostWhereInput{
		ID: &generate2.StringFilter{
			Equals: &id,
			// Gt:     "ss",
		},
	}
	ctx := context.TODO()
	res, err2 := query.Test(ctx, whe, 3)

	fmt.Println(res, err2)

	// testDmmf()
	// testSdl()
}

func testDmmf() {
	engine := engine.NewQueryEngine(GLOBAL_SCHEME, false)
	defer engine.Disconnect()
	if err := engine.Connect(); err != nil {
		panic(err)
	}
	dmmf, err := engine.IntrospectDMMF(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(dmmf.Datamodel)
}

func testSdl() ([]byte, *ast.Schema) {
	engine := engine.NewQueryEngine(GLOBAL_SCHEME, false)
	defer engine.Disconnect()
	if err := engine.Connect(); err != nil {
		panic(err)
	}
	ctx := context.TODO()
	sdl, err := engine.IntrospectSDL(ctx)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(sdl))

	if err := os.WriteFile("schema2.graphql", sdl, 0666); err != nil {
		panic(err)
	}

	schema := gqlparser.MustLoadSchema(&ast.Source{Input: string(sdl)})
	return sdl, schema
}
