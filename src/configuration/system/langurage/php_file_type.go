package langurage

import (
	"reflect"
	"strings"
	"log"
)

type PhpFileType struct {
}

func NewPhpFileType() *PhpFileType {
	return &PhpFileType{}
}

func StrToPhpFileFormat(frame_type, conf_str, namespace string) string {
	defer func() {
		if r := recover(); r != nil {
			log.Println("不能解析该类型的PHP框架", r)
		}
	}()

	frame_type = strings.Title(strings.ToLower(frame_type))

	php_file_type_obj := NewPhpFileType()
	ref_obj := reflect.ValueOf(&php_file_type_obj).Elem()

	param := make([]reflect.Value, 1)
	param[0] = reflect.ValueOf(map[string]interface{}{
		"conf_str":  conf_str,
		"namespace": namespace,
	})

	res := ref_obj.MethodByName(frame_type).Call(param)[0].String()

	return res
}

// sd框架
func (c *PhpFileType) Sd(data map[string]interface{}) string {
	conf_str := "<?php\n\n"
	conf_arr := strings.Split(data["conf_str"].(string), "\n")

	for _, v := range conf_arr {
		if v == "" {
			continue
		}
		conf_str += "$config" + v + "\n"
	}

	conf_str += "\nreturn $config;\n"

	return conf_str
}

// ci框架
func (c *PhpFileType) Codeigniter(data map[string]interface{}) string {
	conf_str := "<?php  if ( ! defined('BASEPATH')) exit('No direct script access allowed');\n\n"
	conf_arr := strings.Split(data["conf_str"].(string), "\n")

	for _, v := range conf_arr {
		if v == "" {
			continue
		}
		conf_str += "$config" + v + "\n"
	}

	return conf_str
}

// tp框架
func (c *PhpFileType) Thinkphp(data map[string]interface{}) string {
	conf_str := "<?php\n\n"
	conf_arr := strings.Split(data["conf_str"].(string), "\n")

	for _, v := range conf_arr {
		if v == "" {
			continue
		}
		conf_str += "$config" + v + "\n"
	}

	conf_str += "\nreturn $config;\n"

	return conf_str
}
