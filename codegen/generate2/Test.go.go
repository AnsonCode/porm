// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc porm_V1.20
// source: Test

package generate2

import (
	"context"
	"fmt"

	"github.com/AnsonCode/porm/utils"
)

const test = `# name: Test 
	query Test ($whe: PostWhereInput!, $tak: Int!) {
	findManyPost(where: $whe, take: $tak, skip: 0) {
		id
		published
		title
		author {
			desc
			id
			sex
		}
	}
}

	`

//,whe *PostWhereInput

//,whe *PostWhereInput,tak int32

func (t Queries) Test(ctx context.Context, whe *PostWhereInput, tak int32) (res *TestResponse, err error) {

	input := map[string]interface{}{

		"whe": whe, //*PostWhereInput

		"tak": tak, //int32

	}
	qry, _ := utils.InlineQuery(test, input) //这里要优化？
	err = Do(ctx, t.e, qry, &res)
	if err != nil {
		fmt.Println(err)
	}
	return
}
