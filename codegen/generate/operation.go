package generate

import (
	"context"
	"fmt"

	"github.com/AnsonCode/porm/utils"
)

//Query

// 合成响应结构体
type QueryTestResponse struct {
	FindManyPost []FindManyPost `json:"findManyPost"`
}

// 这个值需要合成进来~
const testStr = `
# Write your query or mutation here
query Test($whe: PostWhereInput!, $tak: Int!) {
	findManyPost(where: $whe, take: $tak, skip: 0) {
		id
		published
		title
		author {
			desc
			id
		}
	}
}
`

//,whe *PostWhereInput

//,whe *PostWhereInput,tak int

func (t Queries) Test(ctx context.Context, whe *PostWhereInput, tak int) (res *QueryTestResponse, err error) {

	input := map[string]interface{}{

		"whe": whe,

		"tak": tak,
	}
	qry, _ := utils.InlineQuery(testStr, input) //这里要优化？
	err = Do(ctx, t.e, qry, &res)
	if err != nil {
		fmt.Println(err)
	}
	return
}
