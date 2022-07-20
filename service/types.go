package service

import (
	"fmt"
	"reflect"

	lib_reflect "github.com/joycastle/casual-server-lib/reflect"
)

type MatchingRPCResponse struct {
	Code   int    `json:"code"`
	Errmsg string `json:"errmsg"`
	Data   string `json:"data"`
}

var TypesInt map[reflect.Kind]struct{}

func init() {
	TypesInt = map[reflect.Kind]struct{}{
		reflect.Int:    struct{}{},
		reflect.Int8:   struct{}{},
		reflect.Int16:  struct{}{},
		reflect.Int32:  struct{}{},
		reflect.Int64:  struct{}{},
		reflect.Uint:   struct{}{},
		reflect.Uint8:  struct{}{},
		reflect.Uint16: struct{}{},
		reflect.Uint32: struct{}{},
		reflect.Uint64: struct{}{},
	}
}

func MergeFileds(fileds []string, defaultFileds ...string) []string {
	defautAll := []string{"*"}
	filedMap := make(map[string]struct{})
	for _, v := range defaultFileds {
		if v == "*" {
			return defautAll
		}
		filedMap[v] = struct{}{}
	}
	for _, v := range fileds {
		if v == "*" {
			return defautAll
		}
		filedMap[v] = struct{}{}
	}

	newFileds := []string{}
	for k, _ := range filedMap {
		newFileds = append(newFileds, "`"+k+"`")
	}

	return newFileds
}

func MergeFiledsKV(args ...any) ([]string, error) {
	length := len(args)
	if length == 0 || length%2 != 0 {
		return nil, fmt.Errorf("kv num is fatal")
	}
	out := []string{}
	for i := 1; i < length; i += 2 {
		if reflect.TypeOf(args[i-1]).Kind() != reflect.String {
			return nil, fmt.Errorf("the key in index %d must be string", i)
		}

		if reflect.TypeOf(args[i]).Kind() == reflect.String {
			out = append(out, fmt.Sprintf("`%s`='%s'", args[i-1], args[i]))
		} else if lib_reflect.IsIntType(args[i]) {
			out = append(out, fmt.Sprintf("`%s`=%d", args[i-1], args[i]))
		}
	}
	return out, nil
}
