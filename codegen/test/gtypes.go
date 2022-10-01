package test

import (
	"time"
)

// 基本定义

type StringFilter struct {
	Equals    string   `structs:"equals"`
	In        []string `structs:"in"`
	NotIn     []string `structs:"notIn"`
	Lt        string   `structs:"lt"`
	Lte       string   `structs:"lte"`
	Gt        string   `structs:"gt"`
	Gte       string   `structs:"gte"`
	Contains  string   `structs:"contains"`
	StartWith string   `structs:"startWith"`
	EndsWith  string   `structs:"endsWith"`
	// Not NestedStringFilter
}

type StringNullableFilter struct {
	Equals     string
	In         []string
	NotIn      []string
	Lt         string
	Lte        string
	Gt         string
	Gte        string
	Contains   string
	StartsWith string
	EndsWith   string
	// not NestedStringNullableFilter
}

type BoolFilter struct {
	Equals bool `structs:"equals"`
	// Not NestedBoolFilter
}

// 类型待确认
type DateTimeFilter struct {
	Equals    time.Time `structs:"equals"` // 道理是什么呢？
	In        []string  `structs:"in"`
	NotIn     []string  `structs:"notIn"`
	Lt        string    `structs:"lt"`
	Lte       string    `structs:"lte"`
	Gt        string    `structs:"gt"`
	Gte       string    `structs:"gte"`
	Contains  string    `structs:"contains"`
	StartWith string    `structs:"startWith"`
	EndsWith  string    `structs:"endsWith"`
	// Not NestedDateTimeFilter
}

// 动态定义
// "INPUT_OBJECT"
type UserRelationFilter struct {
	Is    *UserWhereInput
	IsNot *UserWhereInput
}

// "INPUT_OBJECT"
type UserWhereInput struct {
	AND       *UserWhereInput
	OR        []UserWhereInput
	NOT       *UserWhereInput
	Id        StringFilter
	CreatedAt *DateTimeFilter
	UpdatedAt *DateTimeFilter
	Name      *StringFilter
	Gender    *BoolFilter
	Desc      *StringNullableFilter
	Post      *PostListRelationFilter
}

// "INPUT_OBJECT"
type PostListRelationFilter struct {
	Every *PostWhereInput
	Some  *PostWhereInput
	None  *PostWhereInput
}

// "INPUT_OBJECT"
type PostWhereInput struct {
	AND       *PostWhereInput
	OR        []PostWhereInput
	NOT       *PostWhereInput
	Id        *StringFilter `structs:"id"`
	CreatedAt *DateTimeFilter
	UpdatedAt *DateTimeFilter
	Title     *StringFilter
	Published *BoolFilter
	Desc      *StringNullableFilter
	Author    *UserRelationFilter // "UserRelationFilter"
	UserId    *StringFilter
}

type PostWhereUniqueInput struct {
	id string
}

// enum SortOrder {
// asc
// desc
// }

type SortOrder string // 只有两种值 desc  asc

type PostOrderByWithRelationInput struct {
	Id        *SortOrder
	CreatedAt *SortOrder
	UpdatedAt *SortOrder
	Title     *SortOrder
	Published *SortOrder
	Desc      *SortOrder
	// author    *UserOrderByWithRelationInput
	UserId *SortOrder
}

func NewDistinctInput() *PostScalarFieldEnum {

	return (&PostScalarFieldEnum{})
}

type PostScalarFieldEnum []string // 这里是枚举
func (s *PostScalarFieldEnum) Id() *PostScalarFieldEnum {
	// 不优雅
	tmp := []string(*s)
	tmp2 := append(tmp, "id")
	tmp3 := PostScalarFieldEnum(tmp2)
	return &tmp3
}
func (s *PostScalarFieldEnum) Title() *PostScalarFieldEnum {
	// 不优雅
	tmp := []string(*s)
	tmp2 := append(tmp, "title")
	tmp3 := PostScalarFieldEnum(tmp2)
	return &tmp3
}

// enum PostScalarFieldEnum {
// id
// createdAt
// updatedAt
// title
// published
// desc
// userId
// }

type TestInput struct {
	Whe *PostWhereInput `structs:"whe"`
	// orderBy: [PostOrderByWithRelationInput]
	// Cursor   *PostWhereUniqueInput
	// Skip     int                  `structs:"tak"`
	Tak      int                  `structs:"tak"`
	Distinct *PostScalarFieldEnum //PostScalarFieldEnum
}

// ): [Post]!

// type Post {
// id: String!
// createdAt: DateTime!
// updatedAt: DateTime!
// title: String!
// published: Boolean!
// desc: String
// author: User!
// userId: String!
// }
