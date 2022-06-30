package service

import "reflect"

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
