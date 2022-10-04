// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc porm_V1.20
// source: Test4

package tutorial

import (
	"context"
	"fmt"
)

const test4 = `# name: Test4 
	query Test4 ($tak: Int!) {
	findManyPost(where: {id:{equals:"1"}}, take: $tak, skip: 0) {
		id
		updatedAt
		title
		author {
			desc
			id
		}
	}
}

	`

//,tak int32

func (t Queries) Test4(ctx context.Context, tak int32) (res *Test4Response, err error) {

	input := map[string]interface{}{

		"tak": tak, //int32

	}
	qry, _ := InlineQuery(test4, input) //这里要优化？
	err = Do(ctx, t.e, qry, &res)
	if err != nil {
		fmt.Println(err)
	}
	return
}
