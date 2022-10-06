package utils

import (
	"bytes"
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

// 位置挪走
func FormatOperateionDocument(operate *ast.OperationDefinition) string {

	// 支持SQL指令
	for _, sel := range operate.SelectionSet {
		filed, ok := sel.(*ast.Field)
		if ok {
			for _, directive := range filed.Directives {
				if directive.Name == "sql" {
					sql := directive.Arguments.ForName("raw").Value.Raw
					// ".*\\$\\{([a-z]+)\\}.*"

					fmt.Println(sql)
					// directive.Arguments
					// TODO: 更改当前的结构

					res := fmt.Sprintf(`
					 mutation %s{
						  %s:queryRaw(query: "%s ", parameters:"[]" )
						}
					`, sel.GetPosition().Src.Name, filed.Alias, sql)
					return res
				}
			}
		}
	}

	query := &ast.QueryDocument{
		Operations: ast.OperationList{operate},
	}
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)

	bufstr := buf.String()

	return bufstr
}
