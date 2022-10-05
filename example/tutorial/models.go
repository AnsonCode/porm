// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc porm_V1.20

package tutorial

import "time"

type DateTime time.Time

const (
	timeFormart = "2006-01-02T15:04:05.000Z"
)

func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = DateTime(now)
	return
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t DateTime) String() string {
	return time.Time(t).Format(timeFormart)
}

// 枚举定义开始

// enum
type Sex string

const (
	SexMALE   Sex = "MALE"
	SexFEMAL  Sex = "FEMAL"
	SexUNKOWN Sex = "UNKOWN"
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
type Test2Author struct {
	// -
	Desc *string `json:"desc"`
	// -
	ID string `json:"id"`
	// -
	Sex *Sex `json:"sex"`
}

// OBJECT_Query_PART
type Test2FindManyPost struct {
	// -
	ID string `json:"id"`
	// -
	Published bool `json:"published"`
	// -
	Title string `json:"title"`
	// -
	Author *Test2Author `json:"author"`
}

// res2_struct
type Test2Response struct {
	// -
	FindManyPost []*Test2FindManyPost `json:"findManyPost"`
}

// OBJECT_Post_PART
type Test3Author struct {
	// -
	Desc *string `json:"desc"`
	// -
	ID string `json:"id"`
}

// OBJECT_Query_PART
type Test3FindManyPost struct {
	// -
	ID string `json:"id"`
	// -
	Published bool `json:"published"`
	// -
	Title string `json:"title"`
	// -
	Author *Test3Author `json:"author"`
}

// res2_struct
type Test3Response struct {
	// -
	FindManyPost []*Test3FindManyPost `json:"findManyPost"`
}

// OBJECT_Post_PART
type Test4Author struct {
	// -
	Desc *string `json:"desc"`
	// -
	ID string `json:"id"`
}

// OBJECT_Query_PART
type Test4FindManyPost struct {
	// -
	ID string `json:"id"`
	// -
	UpdatedAt *DateTime `json:"updatedAt"`
	// -
	Title string `json:"title"`
	// -
	Author *Test4Author `json:"author"`
}

// res2_struct
type Test4Response struct {
	// -
	FindManyPost []*Test4FindManyPost `json:"findManyPost"`
}

// OBJECT_Query_PART
type TestFindUniquePost struct {
	// -
	ID string `json:"id"`
	// -
	CreatedAt *DateTime `json:"createdAt"`
	// -
	UpdatedAt *DateTime `json:"updatedAt"`
	// -
	Title string `json:"title"`
	// -
	Published bool `json:"published"`
	// -
	Desc *string `json:"desc"`
}

// res2_struct
type TestResponse struct {
	// -
	FindUniquePost *TestFindUniquePost `json:"result"`
}

// OBJECT_Post_PART
type TestTimeAuthor struct {
	// -
	ID string `json:"id"`
}

// OBJECT_Query_PART
type TestTimeFindManyPost struct {
	// -
	ID string `json:"id"`
	// -
	Title string `json:"title"`
	// -
	CreatedAt *DateTime `json:"createdAt"`
	// -
	Author *TestTimeAuthor `json:"author"`
}

// res2_struct
type TestTimeResponse struct {
	// -
	FindManyPost []*TestTimeFindManyPost `json:"findManyPost"`
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
