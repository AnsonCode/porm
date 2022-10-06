package main

import (
	"encoding/json"
	"fmt"
)

type TestStruct struct {
	QueryRaw *Json `json:"queryRaw"`
}
type Json struct {
	json.RawMessage
}

func main() {
	input := `
       {
        "Type": 1,
        "Body":[
			{ 
            "Name":"ff",
            "Age" : 19
         }
		]
    }`
	input = `{"queryRaw":[{"id":"1","createdAt":"2022-10-02T21:51:44+00:00","updatedAt":"2022-10-02T21:51:47+00:00","title":"SS","published":1,"desc":"SSS","userId":"1"}]}`
	// ""
	// var res tutorial.RawSqlResponse
	ts := TestStruct{}

	if err := json.Unmarshal([]byte(input), &ts); err != nil {
		panic(err)
	}

	fmt.Print(string(ts.QueryRaw.RawMessage))

	// sql := "SELECT * FROM Post WHERE id= {{.id}}"

	// switch ts.Type {
	// case 1:
	// 	var p Person
	// 	if err := json.Unmarshal(ts.Body, &p); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(p)
	// case 2:
	// 	var w Worker
	// 	if err := json.Unmarshal(ts.Body, &w); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(w)
	// }

}
