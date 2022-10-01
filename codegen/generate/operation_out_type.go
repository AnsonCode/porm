// 目前发现可能会重复
package generate

type FindManyPost struct {

	// 这里做布尔类型的判断 String

	Id *string //`structs:"id,omitempty"`

	// 这里做布尔类型的判断 Boolean

	Published *bool //`structs:"published,omitempty"`

	// 这里做布尔类型的判断 String

	Title *string //`structs:"title,omitempty"`

	// 这里做布尔类型的判断 User

	Author *Author //`structs:"author,omitempty"`

}

type Author struct {

	// 这里做布尔类型的判断 String

	Desc *string //`structs:"desc,omitempty"`

	// 这里做布尔类型的判断 String

	Id *string //`structs:"id,omitempty"`

}
