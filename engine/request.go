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
	"time"

	"github.com/AnsonCode/porm/engine/ast/dmmf"
)

var ErrNotFound = errors.New("ErrNotFound")
var internalUpdateNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to update not found.\"))))"
var internalDeleteNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to delete does not exist.\"))))"

func Do(ctx context.Context, port, query string, variables map[string]interface{}, v interface{}) error {
	// fmt.Println(vars, qry)
	engine := &QueryEngine{
		port: port,
		// hasBinaryTargets: hasBinaryTargets,
		http: &http.Client{},
	}

	return engine.Do(ctx, query, variables, v)
}

func (e *QueryEngine) Do(ctx context.Context, query string, variables map[string]interface{}, v interface{}) error {
	// fmt.Println(vars, qry)
	qry, _ := InlineQuery(query, variables) //这里要优化？
	payload := GQLRequest{
		Query:     qry,
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
	url := "http://localhost:" + e.port

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

func (e *QueryEngine) StartPlayground() {
	handler := NewHandler(e.port)
	pgPort := 8124
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
