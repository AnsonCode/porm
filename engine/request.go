package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AnsonCode/porm/engine/ast/dmmf"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

var ErrNotFound = errors.New("ErrNotFound")
var internalUpdateNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to update not found.\"))))"
var internalDeleteNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to delete does not exist.\"))))"

func Do(ctx context.Context, port int, query string, variables map[string]interface{}, v interface{}) error {
	// fmt.Println(vars, qry)
	engine := &QueryEngine{
		port: port,
		// hasBinaryTargets: hasBinaryTargets,
		http: &http.Client{},
	}

	return engine.Do(ctx, query, variables, v)
}

func (g *QueryEngine) Do(ctx context.Context, query string, variables map[string]interface{}, v interface{}) error {
	// TODO:这里要兼容下sql

	// fmt.Println(vars, qry)
	qry, _ := InlineQuery(query, variables) //这里要优化？

	queryObj, _ := parser.ParseQuery(&ast.Source{Input: qry})

	if len(queryObj.Operations) != 1 {
		return fmt.Errorf("一次只能查询一个operation")
	}
	ope := queryObj.Operations[0]

	if len(ope.SelectionSet) == 1 {
		return g.do2(ctx, qry, v)
	}

	requests := make([]GQLRequest, len(ope.SelectionSet))

	selectionset := ope.SelectionSet
	for i, selection := range selectionset {
		ope.SelectionSet = ast.SelectionSet{selection}
		requests[i] = GQLRequest{
			Query:     FormatOperateionDocument(ope),
			Variables: map[string]interface{}{},
		}
	}
	var result GQLBatchResponse
	payload := GQLBatchRequest{
		Batch:       requests,
		Transaction: true,
	}
	if err := g.Batch(ctx, payload, &result); err != nil {
		return fmt.Errorf("could not send raw query: %w", err)
	}
	if len(result.Errors) > 0 {
		first := result.Errors[0]
		return fmt.Errorf("pql error: %s", first.RawMessage())
	}
	// 合并JSON字符串
	var tmpRes string
	for idx, inner := range result.Result {
		if len(inner.Errors) > 0 {
			first := result.Errors[0]
			return fmt.Errorf("pql error: %s", first.RawMessage())
		}
		str := string(inner.Data)
		// 最后一条
		if idx == len(result.Result)-1 {
			tmpRes = tmpRes + str[1:] // 删除开头的{
		} else {
			// 非最后一条
			tmpRes = tmpRes + str[:len(str)-1] + "," // 删除结尾的}
		}
		// r.queries[i].ExtractQuery().TxResult <- inner.Data.Result
	}
	fmt.Println(tmpRes)
	if err := json.Unmarshal([]byte(tmpRes), v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}
	return nil
	// TODO：更改sql值  SELECT * FROM Post WHERE id= {{.id}}
}
func (g *QueryEngine) RawSQL(ctx context.Context, sql string, variables map[string]interface{}, v interface{}) error {
	// TODO:合成sql
	newsql := sql
	for key, vari := range variables {
		old := "${" + key + "}"
		str := fmt.Sprintf("%v", vari)
		newsql = strings.ReplaceAll(newsql, old, str)
	}

	query := fmt.Sprintf(`
		mutation {
			queryRaw(query: "%s", parameters: "[]")
		}
		`, newsql)

	return g.do2(ctx, query, v)
}

func (e *QueryEngine) do2(ctx context.Context, query string, v interface{}) error {

	payload := GQLRequest{
		Query:     query,
		Variables: map[string]interface{}{},
	}
	err := e.do(ctx, payload, v)
	if err != nil {
		return err
	}
	return nil
}

// Do sends the http Request to the query engine and unmarshals the response
func (e *QueryEngine) do(ctx context.Context, payload interface{}, v interface{}) error {
	startReq := time.Now()

	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	fmt.Printf("[timing] query engine request took %s\n", time.Since(startReq))

	startParse := time.Now()

	var response GQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if len(response.Errors) > 0 {
		first := response.Errors[0]
		if first.RawMessage() == internalUpdateNotFoundMessage ||
			first.RawMessage() == internalDeleteNotFoundMessage {
			return ErrNotFound
		}
		return fmt.Errorf("pql error: %s ", first.RawMessage())
	}

	if err := json.Unmarshal(response.Data, v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	fmt.Printf("[timing] request unmarshaling took %s\n", time.Since(startParse))

	return nil
}

// Do sends the http Request to the query engine and unmarshals the response
func (e *QueryEngine) Batch(ctx context.Context, payload interface{}, v interface{}) error {
	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if err := json.Unmarshal(body, &v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}

func (e *QueryEngine) Request(ctx context.Context, method string, path string, payload interface{}) ([]byte, error) {
	if e.disconnected {
		fmt.Printf("A query was executed after Disconnect() was called. Make sure to not send any queries after disconnecting the client.")
		return nil, fmt.Errorf("client is disconnected")
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("payload marshal: %w", err)
	}

	// TODO use specific log level
	// if logger.Enabled {
	// }
	fmt.Printf("prisma engine payload: `%s`", requestBody)
	url := fmt.Sprintf("http://localhost:%d", e.port)

	return request(ctx, e.http, method, url+path, requestBody, func(req *http.Request) {
		req.Header.Set("content-type", "application/json")
	})
}

func (e *QueryEngine) IntrospectDMMF(ctx context.Context) (*dmmf.Document, error) {
	startReq := time.Now()
	body, err := e.Request(ctx, "GET", "/dmmf", nil)
	if err != nil {
		fmt.Printf("dmmf request failed:  %s", err)
		return nil, err
	}

	fmt.Printf("[timing] query engine dmmf request took %s", time.Since(startReq))

	startParse := time.Now()

	var response dmmf.Document
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("json unmarshal: %s", err)

		return nil, err
	}

	fmt.Printf("[timing] request unmarshaling took %s", time.Since(startParse))

	return &response, nil
}

func (e *QueryEngine) IntrospectSDL(ctx context.Context) ([]byte, error) {

	startReq := time.Now()

	body, err := e.Request(ctx, "GET", "/sdl", nil)
	if err != nil {
		fmt.Printf("sdl request failed:  %s", err)
		return nil, err
	}

	fmt.Printf("[timing] query engine sdl request took %s", time.Since(startReq))

	return body, nil
}

var errNotFound = fmt.Errorf("not found; re-upload schema")

func request(ctx context.Context, client *http.Client, method string, url string, payload []byte, apply func(*http.Request)) ([]byte, error) {
	// if logger.Enabled {
	fmt.Printf("prisma engine payload: `%s`", payload)
	// }

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("raw post: %w", err)
	}

	apply(req)

	req = req.WithContext(ctx)

	startReq := time.Now()
	rawResponse, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("raw post: %w", err)
	}
	defer func() {
		if err := rawResponse.Body.Close(); err != nil {
			panic(err)
		}
	}()
	reqDuration := time.Since(startReq)
	fmt.Printf("[timing] query engine raw request took %s", reqDuration)

	responseBody, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("raw read: %w", err)
	}

	if rawResponse.StatusCode == http.StatusNotFound {
		fmt.Printf("status not found with response body %s", responseBody)
		return nil, errNotFound
	}

	if rawResponse.StatusCode != http.StatusOK && rawResponse.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("http status code %d with response %s", rawResponse.StatusCode, responseBody)
	}

	// if logger.Enabled {
	if elapsedRaw := rawResponse.Header["X-Elapsed"]; len(elapsedRaw) > 0 {
		elapsed, _ := strconv.Atoi(elapsedRaw[0])
		duration := time.Duration(elapsed) * time.Microsecond
		fmt.Printf("[timing] elapsed: %s", duration)

		diff := reqDuration - duration
		fmt.Printf("[timing] just http: %s", diff)
		fmt.Printf("[timing] http percentage: %.2f%%", float64(diff)/float64(reqDuration)*100)
	}
	// }

	return responseBody, nil
}

func (e *QueryEngine) startPlayground() {
	handler := NewHandler(e.port)
	pgPort := e.playground_port
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", pgPort),
		Handler: handler,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			// log.Fatalln("listen and serve", err)
			log.Panic("listen and serve", err)
		}
	}()
	fmt.Printf("playground，访问：http://localhost:%d\n", pgPort)
}
