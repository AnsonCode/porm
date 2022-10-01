package generate2

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/runtime/types"
)

func NewClient(path string) *Queries {
	connectionStr := mustReadFile(path)
	e := engine.NewQueryEngine(connectionStr, false)
	return &Queries{
		e,
	}
}
func NewClient2(connectionStr string) *Queries {
	e := engine.NewQueryEngine(connectionStr, false)
	return &Queries{
		e,
	}
}

type Queries struct {
	e *engine.QueryEngine
}

func (q *Queries) Connect() {
	if err := q.e.Connect(); err != nil {
		panic(err)
	}
}

func (q *Queries) Disconnect() {
	if q.e != nil {
		q.e.Disconnect()
	}
}

// Do sends the http Request to the query engine and unmarshals the response
func Do(ctx context.Context, e *engine.QueryEngine, qry string, v interface{}) error {

	// fmt.Println(vars, qry)
	payload := engine.GQLRequest{
		Query:     qry,
		Variables: map[string]interface{}{},
	}
	startReq := time.Now()

	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	fmt.Printf("[timing] query engine request took %s", time.Since(startReq))

	startParse := time.Now()
	type GQLResponse struct {
		Data       json.RawMessage        `json:"data"`
		Errors     []engine.GQLError      `json:"errors"`
		Extensions map[string]interface{} `json:"extensions"`
	}
	var response GQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if len(response.Errors) > 0 {
		first := response.Errors[0]
		var internalUpdateNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to update not found.\"))))"
		var internalDeleteNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to delete does not exist.\"))))"

		if first.RawMessage() == internalUpdateNotFoundMessage ||
			first.RawMessage() == internalDeleteNotFoundMessage {
			return types.ErrNotFound
		}
		return fmt.Errorf("pql error: %s", first.RawMessage())
	}

	if err := json.Unmarshal(response.Data, v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	fmt.Printf("[timing] request unmarshaling took %s", time.Since(startParse))

	return nil
}

func mustReadFile(name string) string {
	src, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return string(src)
}
