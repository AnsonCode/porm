package generate

type PostListRelationFilter struct {
	Every *PostWhereInput `json:"every,omitempty"`

	Some *PostWhereInput `json:"some,omitempty"`

	None *PostWhereInput `json:"none,omitempty"`
}

type PostWhereInput struct {
	AND *PostWhereInput `json:"AND,omitempty"`

	OR *[]PostWhereInput `json:"OR,omitempty"`

	NOT *PostWhereInput `json:"NOT,omitempty"`

	Id *StringFilter `json:"id,omitempty"`

	CreatedAt *DateTimeFilter `json:"createdAt,omitempty"`

	UpdatedAt *DateTimeFilter `json:"updatedAt,omitempty"`

	Title *StringFilter `json:"title,omitempty"`

	Published *BoolFilter `json:"published,omitempty"`

	Desc *StringNullableFilter `json:"desc,omitempty"`

	Author *UserRelationFilter `json:"author,omitempty"`

	UserId *StringFilter `json:"userId,omitempty"`
}

type UserRelationFilter struct {
	Is *UserWhereInput `json:"is,omitempty"`

	IsNot *UserWhereInput `json:"isNot,omitempty"`
}

type UserWhereInput struct {
	AND *UserWhereInput `json:"AND,omitempty"`

	OR *[]UserWhereInput `json:"OR,omitempty"`

	NOT *UserWhereInput `json:"NOT,omitempty"`

	Id *StringFilter `json:"id,omitempty"`

	CreatedAt *DateTimeFilter `json:"createdAt,omitempty"`

	UpdatedAt *DateTimeFilter `json:"updatedAt,omitempty"`

	Name *StringFilter `json:"name,omitempty"`

	Gender *BoolFilter `json:"gender,omitempty"`

	Desc *StringNullableFilter `json:"desc,omitempty"`

	Post *PostListRelationFilter `json:"Post,omitempty"`
}
