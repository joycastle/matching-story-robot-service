package lib

import (
	"encoding/json"
	"reflect"
)

type LogStructuredJson struct {
	State   string         `json:"state"`
	ErrStep int            `json:"step,omitempty"`
	ErrMsg  string         `json:"errmsg,omitempty"`
	Detail  map[string]any `json:"detail,omitempty"`
}

func NewLogStructed() *LogStructuredJson {
	return &LogStructuredJson{
		State:  "SUCCESS",
		Detail: make(map[string]any),
	}
}

func (lsj *LogStructuredJson) Set(args ...any) *LogStructuredJson {
	length := len(args)
	if length == 0 || length%2 != 0 {
		return lsj
	}
	for i := 1; i < length; i += 2 {
		if reflect.TypeOf(args[i-1]).Kind() == reflect.String && reflect.TypeOf(args[i]).Kind() != reflect.Interface {
			lsj.Detail[args[i-1].(string)] = args[i]
		}
	}
	return lsj
}

func (lsj *LogStructuredJson) Step(step int) *LogStructuredJson {
	lsj.ErrStep = step
	return lsj
}

func (lsj *LogStructuredJson) Failed() *LogStructuredJson {
	lsj.State = "FAILED"
	return lsj
}

func (lsj *LogStructuredJson) IsFailed() bool {
	return lsj.State == "FAILED"
}

func (lsj *LogStructuredJson) Success() *LogStructuredJson {
	lsj.State = "SUCCESS"
	return lsj
}

func (lsj *LogStructuredJson) Err(err error) *LogStructuredJson {
	if err != nil {
		lsj.ErrMsg = err.Error()
	}
	return lsj
}

func (lsj *LogStructuredJson) ErrString(s string) *LogStructuredJson {
	lsj.ErrMsg = s
	return lsj
}

func (lsj *LogStructuredJson) String() string {
	b, _ := json.Marshal(lsj)
	return string(b)
}

func (lsj *LogStructuredJson) Merge(t *LogStructuredJson) *LogStructuredJson {
	lsj.State = t.State
	lsj.ErrStep = t.ErrStep
	lsj.ErrMsg = t.ErrMsg
	for k, v := range t.Detail {
		lsj.Set(k, v)
	}
	return lsj
}
