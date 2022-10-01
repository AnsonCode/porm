原计划采用修改模型的方式做内联，但循环构造 Value 对象比较麻烦~

	sets := []ast.Selection(doc.Operations[0].SelectionSet)
	set := sets[0].(*ast.Field)
	res := set.ArgumentMap(vars)
	fmt.Println(res)
	// argments := []*ast.Argument(set.Arguments)
	// for _, v := range argments {

	// 	if v.Value.Kind == ast.Variable {
	// 		// TODO:应该新建value值

	// 		// val := variable.ForName(v.Value.Raw)
	// 		k, ok := vars[v.Value.Raw]
	// 		if !ok {
	// 			continue
	// 		}
	// 		switch vv := k.(type) {
	// 		case string:
	// 			fmt.Println(k, "is string", vv)
	// 			// v.Value.Raw = k.(string) // 根据raw 去变量里面选值
	// 			v.Value.Kind = ast.StringValue

	// 		case int:
	// 			fmt.Println(k, "is int ", vv)
	// 			// v.Value.Raw = fmt.Sprintf("%d", k) // 根据raw 去变量里面选值
	// 			v.Value.Kind = ast.IntValue
	// 		case float64:
	// 			fmt.Println(k, "is float64 ", vv)
	// 			// v.Value.Raw = fmt.Sprintf("%d", k) // 根据raw 去变量里面选值
	// 			v.Value.Kind = ast.FloatValue
	// 		case map[string]interface{}:
	// 			fmt.Println(k, "is obj ", vv)
	// 			// b, _ := json.Marshal(k)
	// 			// v.Value.Raw = string(b) // 根据raw 去变量里面选值]
	// 			v.Value.Kind = ast.ObjectValue

	// 		case []interface{}:
	// 			fmt.Println(k, "is array ", vv)
	// 			// b, _ := json.Marshal(k)
	// 			// v.Value.Raw = string(b) // 根据raw 去变量里面选值
	// 			v.Value.Kind = ast.ListValue
	// 		}
	// 		v.Value.Value(vars)

	// 	}
	// }