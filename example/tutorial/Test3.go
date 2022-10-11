// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc porm_V1.20
// source: Test3

package tutorial

import (
	"context"
	"fmt"
)

//# name: Test3
const test3 = `query Test3 ($whe: PostWhereInput!, $tak: Int!) {
	findManyPost(where: $whe, take: $tak, skip: 0) {
		id
		published
		title
		author {
			desc
			id
		}
	}
	findFirstPost {
		desc
	}
}
`

//,whe *PostWhereInput

//,whe *PostWhereInput,tak int32

func (t Queries) Test3(ctx context.Context, whe *PostWhereInput, tak int32) (res *Test3Response, err error) {
	input := map[string]interface{}{

		"whe": whe, //*PostWhereInput

		"tak": tak, //int32

	}
	err = t.client.Do(ctx, test3, input, &res)
	if err != nil {
		fmt.Println(err)
	}
	return
}
