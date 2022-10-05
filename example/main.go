package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnsonCode/porm/engine"
	"github.com/AnsonCode/porm/example/tutorial"
)

func mustReadFile(name string) string {
	src, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return string(src)
}

func main() {
	ctx := context.TODO()

	// _, schema := testSdl()
	GLOBAL_SCHEME := mustReadFile("./schema.prisma")
	engine := engine.NewQueryEngine(GLOBAL_SCHEME, 8123, "./query-engine")
	engine.Connect()
	defer engine.Disconnect()
	engine.StartPlayground()

	// 内省获取 graphql schema
	sdl, err := engine.IntrospectSDL(ctx)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("schema2.graphql", sdl, 0666); err != nil {
		panic(err)
	}

	// TODO:生成schema.graphql的方式，以便用户好读取
	query := tutorial.NewClient(engine)

	id := "1"
	whe := &tutorial.PostWhereInput{
		ID: &tutorial.StringFilter{
			Equals: &id,
			// Gt:     "ss",
		},
	}
	res, err2 := query.Test3(ctx, whe, 3)
	data, _ := json.MarshalIndent(res, "", "\t")
	fmt.Println(string(data), err2)
	sex := tutorial.SexMALE
	restime, err3 := query.TestTime(ctx, 2, []string{"dd"}, &sex)
	data, _ = json.MarshalIndent(restime, "", "\t")
	fmt.Println(string(data), err3)

	resraw, err3 := query.RawSql(ctx)
	data, _ = json.MarshalIndent(resraw, "", "\t")
	fmt.Println(string(data), err3)

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
