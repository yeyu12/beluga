package drive

import (
	"sort"
	"strings"
	"reflect"
)

type Configuration struct {
	Key string
	Val string
}

type SortConfig []Configuration

func SetKeyToJson(config []Configuration) map[string]interface{} {
	// 排序配置,按照分割大小排序
	sort.Sort(SortConfig(config))

	// 配置分组,使用结构体来实现
	data := setConfigMap(config)

	return data
}

func setConfigMap(config []Configuration) map[string]interface{} {
	config_map := make(map[string]interface{})

	for _, val := range config {
		key_arr := strings.Split(val.Key, ".")
		key_arr_len := len(key_arr)
		temp_map := make(map[string]interface{})

		for i := 0; i < key_arr_len; i++ {
			if i == (key_arr_len - 1) {
				strToMap(key_arr, key_arr[i],0, i, val.Val, temp_map)
			} else {
				strToMap(key_arr, key_arr[i], 0, i, make(map[string]interface{}), temp_map)
			}
		}

		config_map = jsonMerge(config_map, temp_map)
	}

	return config_map
}

func jsonMerge(dst, src map[string]interface{}) map[string]interface{} {
	return jsMerge(dst, src, 0)
}

func jsMerge(dst, src map[string]interface{}, depth int) map[string]interface{} {
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

// key=》val转map格式
// params key_arr []string 所有的key
// params key string 要添加的key
// params base_index int 要开始找的下标
// params index int key下标位置所在
// params val interface 值
// params temp_map map[string]interface 要保存的数据
func strToMap(key_arr []string, key string, base_index, index int, val interface{}, temp_map map[string]interface{}) {
	if (base_index == index) && (key_arr[index] == key) {
		temp_map[key] = val
	} else {
		test_map_temp := temp_map[key_arr[base_index]]

		if test_map_temp == nil {
			test_map_temp = make(map[string]interface{})
		}

		base_index++
		strToMap(key_arr, key, base_index, index, val, test_map_temp.(map[string]interface{}))
	}

	/*if key_arr[0] == key {
		temp_map[key] = val
	} else {
		// 剩余的键
		key_arr_temp := key_arr[1:]
		// 当前键
		test_map_temp := temp_map[key_arr[0]]

		if test_map_temp == nil {
			test_map_temp = make(map[string]interface{})
		}

		index++

		strToMap(key_arr_temp, key, index, val, test_map_temp.(map[string]interface{}))
	}*/
}

// 重写 Len() 方法
func (a SortConfig) Len() int {
	return len(a)
}

// 重写 Swap() 方法
func (a SortConfig) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// 重写 Less() 方法， 指定排序规则
func (a SortConfig) Less(i, j int) bool {
	if a[i].Key[len(a[i].Key)-1] == '.' {
		a[i].Key = string([]byte(a[i].Key)[0 : len(a[i].Key)-1])
	}

	return len(strings.Split(a[j].Key, ".")) < len(strings.Split(a[i].Key, "."))
}
