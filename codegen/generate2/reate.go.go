// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc porm_V1.20
// source: reate

package generate2

import (
	"context"
	"fmt"

	"github.com/AnsonCode/porm/utils"
)

const reate = `# name: Reate 
	mutation reate ($title: String!, $published: Boolean!, $userId: String!) {
	createOnePost(data: {title:$title,published:$published,author:{connect:{id:$userId}}}) {
		id
		userId
	}
}

	`

//,title string

//,title string,published bool

//,title string,published bool,userId string

func (t Queries) Reate(ctx context.Context, title string, published bool, userId string) (res *ReateResponse, err error) {

	input := map[string]interface{}{

		"title": title, //string

		"published": published, //bool

		"userId": userId, //string

	}
	qry, _ := utils.InlineQuery(reate, input) //这里要优化？
	err = Do(ctx, t.e, qry, &res)
	if err != nil {
		fmt.Println(err)
	}
	return
}