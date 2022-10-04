package utils

import (
	"bytes"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

func FormatOperateionDocument(operate *ast.OperationDefinition) string {

	query := &ast.QueryDocument{
		Operations: ast.OperationList{operate},
	}
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)

	bufstr := buf.String()

	return bufstr
}
