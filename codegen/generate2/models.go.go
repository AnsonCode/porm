// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc test

package generate2

// input_struct
type PostListRelationFilter struct {
	// df
	Every *PostWhereInput `json:"every,omitempty"`
	// df
	Some *PostWhereInput `json:"some,omitempty"`
	// df
	None *PostWhereInput `json:"none,omitempty"`
}

// input_struct
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

// res2_struct
type TestResponse struct {
	// -
	FindManyPost []*TestfindManyPost `json:"findManyPost"`
}

// res_struct
type Testauthor struct {
	// -
	Desc *string `json:"desc"`
	// -
	ID string `json:"id"`
}

// res_struct
type TestfindManyPost struct {
	// -
	ID string `json:"id"`
	// -
	Published bool `json:"published"`
	// -
	Title string `json:"title"`
	// -
	Author *Testauthor `json:"author"`
}

// input_struct
type UserRelationFilter struct {
	// df
	Is *UserWhereInput `json:"is,omitempty"`
	// df
	IsNot *UserWhereInput `json:"isNot,omitempty"`
}

// input_struct
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
	Desc *StringNullableFilter `json:"desc,omitempty"`
	// df
	Post *PostListRelationFilter `json:"Post,omitempty"`
}
