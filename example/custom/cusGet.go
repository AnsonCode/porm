// versions:
//   sqlc porm_V1.20
// source: cusGet

package custom

import (
	"context"
	"fmt"
)

const cusGet = `# name: CusGet 
  mutation cusGet{
	cusGetPostById:queryRaw(query: "SELECT * FROM Post WHERE id= 1 ", parameters:"[]" )
}
	`

func (t Queries) CusGet(ctx context.Context) (res *CusGetResponse, err error) {

	input := map[string]interface{}{}
	err = t.client.Do(ctx, cusGet, input, &res)

	if err != nil {
		fmt.Println(err)
	}
	return
}
