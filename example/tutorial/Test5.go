// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc porm_V1.20
// source: Test5

package tutorial

import (
	"context"
	"fmt"
)

const test5 = `# name: Test5 
	query Test5 ($tak: Int!) {
	findManyPost(take: $tak, skip: 0) {
		id
		title
		author {
			desc
			id
		}
	}
}

	`

//,tak int32

func (t Queries) Test5(ctx context.Context, tak int32) (res *Test5Response, err error) {

	input := map[string]interface{}{

		"tak": tak, //int32

	}
	qry, _ := InlineQuery(test5, input) //这里要优化？
	err = Do(ctx, t.e, qry, &res)
	if err != nil {
		fmt.Println(err)
	}
	return
}
