package engine

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_convert(t *testing.T) {
	s := `{"id":{"equals":"ssss"}} ["sss","ddd"]`
	reg := regexp.MustCompile("\"(\\w+)\"(\\s*:\\s*)")
	res := reg.ReplaceAllString(s, "$1$2")
	fmt.Println(res) // Abraham Lincoln
	if res != `{id:{equals:"ssss"}} ["sss","ddd"]` {
		t.Fail()
	}
}
