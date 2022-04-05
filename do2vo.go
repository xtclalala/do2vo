package do2vo

import (
	"fmt"
	"reflect"
	"strings"
)

func GetResponse(d any, target string) (any, error) {
	typeOf := reflect.TypeOf(d)
	if typeOf.Kind() != reflect.Struct {
		if typeOf.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("不是结构体")
		}
	}
	valueOf := reflect.ValueOf(d)
	response, err := rec(valueOf, target)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func rec(v reflect.Value, target string) (map[any]any, error) {

	response := map[any]any{}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		tField := t.Field(i)
		tag, err := tField.Tag.Lookup("y1t")
		if err == false {
			continue
		}
		if ok := check(tag, target); ok == false {
			continue
		}
		key := getJsonTag(tField)
		switch kind := tField.Type.Kind(); kind {
		case reflect.String:
			response[key] = vField.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			response[key] = vField.Int()
		case reflect.Bool:
			response[key] = vField.Bool()
		case reflect.Float64, reflect.Float32:
			response[key] = vField.Float()
		case reflect.Struct:
			res, err := rec(vField, target)
			if err != nil {
				return nil, err
			}
			response[key] = res
		case reflect.Slice, reflect.Array:
			var slices []map[any]any
			for item := 0; item < vField.Len(); item++ {
				res, err := rec(vField.Index(item), target)
				if err != nil {
					return nil, err
				}
				slices = append(slices, res)
			}
			response[key] = slices
		}
	}
	return response, nil
}

func check(value string, target string) bool {
	values := strings.Split(value, ",")
	for _, i := range values {
		if target == i || i == "-" {
			return true
		}
	}
	return false
}

func getJsonTag(t reflect.StructField) string {
	if j, ok := t.Tag.Lookup("json"); ok == false {
		return t.Name
	} else {
		return j
	}
}
