// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc porm_V1.20

package generate2

import "time"

type DateTime time.Time

// 枚举定义开始

// tes
type Test string

const (
	TestTest Test = "test"
)

// 结构定义开始

// INPUT_OBJECT
type BoolFilter struct {
	// df
	Equals *bool `json:"equals,omitempty"`
	// df
	Not *NestedBoolFilter `json:"not,omitempty"`
}

// INPUT_OBJECT
type DateTimeFilter struct {
	// df
	Equals *DateTime `json:"equals,omitempty"`
	// df
	In []*DateTime `json:"in,omitempty"`
	// df
	NotIn []*DateTime `json:"notIn,omitempty"`
	// df
	Lt *DateTime `json:"lt,omitempty"`
	// df
	Lte *DateTime `json:"lte,omitempty"`
	// df
	Gt *DateTime `json:"gt,omitempty"`
	// df
	Gte *DateTime `json:"gte,omitempty"`
	// df
	Not *NestedDateTimeFilter `json:"not,omitempty"`
}

// INPUT_OBJECT
type EnumSexFilter struct {
	// df
	Equals *Sex `json:"equals,omitempty"`
	// df
	In []*Sex `json:"in,omitempty"`
	// df
	NotIn []*Sex `json:"notIn,omitempty"`
	// df
	Not *Sex `json:"not,omitempty"`
}

// INPUT_OBJECT
type NestedBoolFilter struct {
	// df
	Equals *bool `json:"equals,omitempty"`
	// df
	Not *NestedBoolFilter `json:"not,omitempty"`
}

// INPUT_OBJECT
type NestedDateTimeFilter struct {
	// df
	Equals *DateTime `json:"equals,omitempty"`
	// df
	In []*DateTime `json:"in,omitempty"`
	// df
	NotIn []*DateTime `json:"notIn,omitempty"`
	// df
	Lt *DateTime `json:"lt,omitempty"`
	// df
	Lte *DateTime `json:"lte,omitempty"`
	// df
	Gt *DateTime `json:"gt,omitempty"`
	// df
	Gte *DateTime `json:"gte,omitempty"`
	// df
	Not *NestedDateTimeFilter `json:"not,omitempty"`
}

// INPUT_OBJECT
type NestedStringFilter struct {
	// df
	Equals *string `json:"equals,omitempty"`
	// df
	In []string `json:"in,omitempty"`
	// df
	NotIn []string `json:"notIn,omitempty"`
	// df
	Lt *string `json:"lt,omitempty"`
	// df
	Lte *string `json:"lte,omitempty"`
	// df
	Gt *string `json:"gt,omitempty"`
	// df
	Gte *string `json:"gte,omitempty"`
	// df
	Contains *string `json:"contains,omitempty"`
	// df
	StartsWith *string `json:"startsWith,omitempty"`
	// df
	EndsWith *string `json:"endsWith,omitempty"`
	// df
	Not *NestedStringFilter `json:"not,omitempty"`
}

// INPUT_OBJECT
type NestedStringNullableFilter struct {
	// df
	Equals *string `json:"equals,omitempty"`
	// df
	In []string `json:"in,omitempty"`
	// df
	NotIn []string `json:"notIn,omitempty"`
	// df
	Lt *string `json:"lt,omitempty"`
	// df
	Lte *string `json:"lte,omitempty"`
	// df
	Gt *string `json:"gt,omitempty"`
	// df
	Gte *string `json:"gte,omitempty"`
	// df
	Contains *string `json:"contains,omitempty"`
	// df
	StartsWith *string `json:"startsWith,omitempty"`
	// df
	EndsWith *string `json:"endsWith,omitempty"`
	// df
	Not *NestedStringNullableFilter `json:"not,omitempty"`
}

// INPUT_OBJECT
type PostListRelationFilter struct {
	// df
	Every *PostWhereInput `json:"every,omitempty"`
	// df
	Some *PostWhereInput `json:"some,omitempty"`
	// df
	None *PostWhereInput `json:"none,omitempty"`
}

// INPUT_OBJECT
type PostWhereInput struct {
	// df
	AND *PostWhereInput `json:"AND,omitempty"`
	// df
	OR []*PostWhereInput `json:"OR,omitempty"`
	// df
	NOT *PostWhereInput `json:"NOT,omitempty"`
	// df
	ID *StringFilter `json:"id,omitempty"`
	// df
	CreatedAt *DateTimeFilter `json:"createdAt,omitempty"`
	// df
	UpdatedAt *DateTimeFilter `json:"updatedAt,omitempty"`
	// df
	Title *StringFilter `json:"title,omitempty"`
	// df
	Published *BoolFilter `json:"published,omitempty"`
	// df
	Desc *StringNullableFilter `json:"desc,omitempty"`
	// df
	Author *UserRelationFilter `json:"author,omitempty"`
	// df
	UserId *StringFilter `json:"userId,omitempty"`
}

// OBJECT_Mutation_PART
type ReateCreateOnePost struct {
	// -
	ID string `json:"id"`
	// -
	UserId string `json:"userId"`
}

// res2_struct
type ReateResponse struct {
	// -
	CreateOnePost *ReateCreateOnePost `json:"createOnePost"`
}

// ENUM
type Sex struct {
}

// INPUT_OBJECT
type StringFilter struct {
	// df
	Equals *string `json:"equals,omitempty"`
	// df
	In []string `json:"in,omitempty"`
	// df
	NotIn []string `json:"notIn,omitempty"`
	// df
	Lt *string `json:"lt,omitempty"`
	// df
	Lte *string `json:"lte,omitempty"`
	// df
	Gt *string `json:"gt,omitempty"`
	// df
	Gte *string `json:"gte,omitempty"`
	// df
	Contains *string `json:"contains,omitempty"`
	// df
	StartsWith *string `json:"startsWith,omitempty"`
	// df
	EndsWith *string `json:"endsWith,omitempty"`
	// df
	Not *NestedStringFilter `json:"not,omitempty"`
}

// INPUT_OBJECT
type StringNullableFilter struct {
	// df
	Equals *string `json:"equals,omitempty"`
	// df
	In []string `json:"in,omitempty"`
	// df
	NotIn []string `json:"notIn,omitempty"`
	// df
	Lt *string `json:"lt,omitempty"`
	// df
	Lte *string `json:"lte,omitempty"`
	// df
	Gt *string `json:"gt,omitempty"`
	// df
	Gte *string `json:"gte,omitempty"`
	// df
	Contains *string `json:"contains,omitempty"`
	// df
	StartsWith *string `json:"startsWith,omitempty"`
	// df
	EndsWith *string `json:"endsWith,omitempty"`
	// df
	Not *NestedStringNullableFilter `json:"not,omitempty"`
}

// OBJECT_Post_PART
type TestAuthor struct {
	// -
	Desc *string `json:"desc"`
	// -
	ID string `json:"id"`
}

// OBJECT_Query_PART
type TestFindManyPost struct {
	// -
	ID string `json:"id"`
	// -
	Published bool `json:"published"`
	// -
	Title string `json:"title"`
	// -
	Author *TestAuthor `json:"author"`
}

// res2_struct
type TestResponse struct {
	// -
	FindManyPost []*TestFindManyPost `json:"findManyPost"`
}

// INPUT_OBJECT
type UserRelationFilter struct {
	// df
	Is *UserWhereInput `json:"is,omitempty"`
	// df
	IsNot *UserWhereInput `json:"isNot,omitempty"`
}

// INPUT_OBJECT
type UserWhereInput struct {
	// df
	AND *UserWhereInput `json:"AND,omitempty"`
	// df
	OR []*UserWhereInput `json:"OR,omitempty"`
	// df
	NOT *UserWhereInput `json:"NOT,omitempty"`
	// df
	ID *StringFilter `json:"id,omitempty"`
	// df
	CreatedAt *DateTimeFilter `json:"createdAt,omitempty"`
	// df
	UpdatedAt *DateTimeFilter `json:"updatedAt,omitempty"`
	// df
	Name *StringFilter `json:"name,omitempty"`
	// df
	Gender *BoolFilter `json:"gender,omitempty"`
	// df
	Sex *EnumSexFilter `json:"sex,omitempty"`
	// df
	Desc *StringNullableFilter `json:"desc,omitempty"`
	// df
	Post *PostListRelationFilter `json:"Post,omitempty"`
}
