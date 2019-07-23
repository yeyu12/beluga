package main

import (
	"strings"
	"reflect"
	"fmt"
	"encoding/json"
)

// 用来研究key转map然后转json
/*
测试数据
test.test1
asdas.qwa
asdas.qaz
test.test2
test.test22.wer.wer1.dfg.fgh
test.test22.wer2
 */
func main() {

	text_map := map[string]interface{}{
		"test.test1":                   "2",
		"asdas.qwa":                    "342",
		"asdas.qaz":                    "1122",
		"test.test2":                   "1",
		"test.test22.wer.wer1.dfg.fgh": "4",
		"test.test22.wer2":             "5ffff",
	}

	test_map1 := make(map[string]interface{})

	for key, val := range text_map {
		key_arr := strings.Split(key, ".")
		key_arr_len := len(key_arr)
		test_map := make(map[string]interface{})

		for i := 0; i < key_arr_len; i++ {
			if i == (key_arr_len - 1) {
				digui(key_arr, key_arr[i], val,  test_map)
			} else {
				digui(key_arr, key_arr[i], make(map[string]interface{}), test_map)
			}
		}

		test_map1 = JsonMerge(test_map1, test_map)

		//fmt.Println(test)
	}

	json_str, _ := json.Marshal(test_map1)
	fmt.Println(string(json_str))
}

func JsonMerge(dst, src map[string]interface{}) map[string]interface{} {
	return jsMerge(dst, src, 0)
}

var jsonMergeDepth = 32

func jsMerge(dst, src map[string]interface{}, depth int) map[string]interface{} {
	if depth > jsonMergeDepth {
		return dst
	}

	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {

			srcMap, srcMapOk := jsMapify(srcVal)
			dstMap, dstMapOk := jsMapify(dstVal)

			if srcMapOk && dstMapOk {
				srcVal = jsMerge(dstMap, srcMap, depth+1)
			}
		}

		dst[key] = srcVal
	}

	return dst
}

func jsMapify(i interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Map {
		m := map[string]interface{}{}

		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}

		return m, true
	}

	return map[string]interface{}{}, false
}

func digui(key_arr []string, key string, val interface{}, test_map map[string]interface{}) {
	if key_arr[0] == key{
		test_map[key] = val
	}else {
		key_arr_temp := key_arr[1:]
		test_map_temp := test_map[key_arr[0]]
		digui(key_arr_temp, key, val, test_map_temp.(map[string]interface{}))
	}
}
