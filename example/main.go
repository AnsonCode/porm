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

	GLOBAL_SCHEME := mustReadFile("./schema.prisma")
	client := engine.NewQueryEngine(GLOBAL_SCHEME, 8123, 8124, "./prisma/efdf9b1183dddfd4258cd181a72125755215ab7b/prisma-query-engine-darwin")
	client.Connect()
	defer client.Disconnect()

	// TODO:生成schema.graphql的方式，以便用户好读取
	query := tutorial.NewClient(client)

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

	// []byte
	resraw, err3 := query.QueryRaw2(ctx, &tutorial.Json{RawMessage: []byte(`"[1,2]"`)})
	data, _ = json.MarshalIndent(resraw, "", "\t")
	fmt.Println(string(data), err3)

	resraw2, err3 := query.CusGet(ctx, "1")
	data, _ = json.MarshalIndent(resraw2, "", "\t")
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
