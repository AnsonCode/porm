package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnsonCode/porm/engine"
	"github.com/AnsonCode/porm/example/tutorial"
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
	// _, schema := testSdl()

	engine := engine.NewQueryEngine(GLOBAL_SCHEME, "8123", "./query-engine")
	engine.Connect()
	defer engine.Disconnect()

	query := tutorial.NewClient(engine)

	id := "1"
	whe := &tutorial.PostWhereInput{
		ID: &tutorial.StringFilter{
			Equals: &id,
			// Gt:     "ss",
		},
	}
	ctx := context.TODO()
	res, err2 := query.Test3(ctx, whe, 3)

	fmt.Println(res, err2)
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig

	// testDmmf()
	// testSdl()
}

// func testDmmf() {
// 	engine := engine.NewQueryEngine(GLOBAL_SCHEME, false)
// 	defer engine.Disconnect()
// 	if err := engine.Connect(); err != nil {
// 		panic(err)
// 	}
// 	dmmf, err := engine.IntrospectDMMF(context.TODO())
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(dmmf.Datamodel)
// }

// func testSdl() ([]byte, *ast.Schema) {
// 	engine := engine.NewQueryEngine(GLOBAL_SCHEME, false)
// 	defer engine.Disconnect()
// 	if err := engine.Connect(); err != nil {
// 		panic(err)
// 	}
// 	ctx := context.TODO()
// 	sdl, err := engine.IntrospectSDL(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// fmt.Println(string(sdl))

// 	if err := os.WriteFile("schema2.graphql", sdl, 0666); err != nil {
// 		panic(err)
// 	}

// 	schema := gqlparser.MustLoadSchema(&ast.Source{Input: string(sdl)})
// 	return sdl, schema
// }
