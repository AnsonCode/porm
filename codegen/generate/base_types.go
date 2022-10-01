package generate

import (
	"time"
)

// 基本定义

type StringFilter struct {
	Equals    string   `json:"equals,omitempty"`
	In        []string `json:"in,omitempty"`
	NotIn     []string `json:"notIn,omitempty"`
	Lt        string   `json:"lt,omitempty"`
	Lte       string   `json:"lte,omitempty"`
	Gt        string   `json:"gt,omitempty"`
	Gte       string   `json:"gte,omitempty"`
	Contains  string   `json:"contains,omitempty"`
	StartWith string   `json:"startWith,omitempty"`
	EndsWith  string   `json:"endsWith,omitempty"`
	// Not NestedStringFilter
}

type StringNullableFilter struct {
	Equals    string   `json:"equals,omitempty"`
	In        []string `json:"in,omitempty"`
	NotIn     []string `json:"notIn,omitempty"`
	Lt        string   `json:"lt,omitempty"`
	Lte       string   `json:"lte,omitempty"`
	Gt        string   `json:"gt,omitempty"`
	Gte       string   `json:"gte,omitempty"`
	Contains  string   `json:"contains,omitempty"`
	StartWith string   `json:"startWith,omitempty"`
	EndsWith  string   `json:"endsWith,omitempty"`
	// not NestedStringNullableFilter
}

type BoolFilter struct {
	Equals bool `json:"equals,omitempty"`
	// Not NestedBoolFilter
}

// 类型待确认
type DateTimeFilter struct {
	Equals    *time.Time `json:"equals,omitempty"` // 道理是什么呢？
	In        []string   `json:"in,omitempty"`
	NotIn     []string   `json:"notIn,omitempty"`
	Lt        string     `json:"lt,omitempty"`
	Lte       string     `json:"lte,omitempty"`
	Gt        string     `json:"gt,omitempty"`
	Gte       string     `json:"gte,omitempty"`
	Contains  string     `json:"contains,omitempty"`
	StartWith string     `json:"startWith,omitempty"`
	EndsWith  string     `json:"endsWith,omitempty"`
	// Not NestedDateTimeFilter
}

type DateTime time.Time
