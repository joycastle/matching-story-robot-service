package create

import "encoding/json"

var errTexts = map[int]string{
	0: "success",
	//db
	100: "db read error",
	101: "db data empty",
	102: "db update error",
	103: "db filed value is empty",
	104: "db delete error",

	//config
	200: "config read error",
	201: "config is empty",

	//data parse
	300: "parameter is invalid",

	//rpc
	400: "rpc:rpc send error",

	//create
	500: "conditions are not met",
	501: "robot create error",
	502: "join club error",
}

type Result struct {
	Code    int    `json:"code"`
	Errmsg  string `json:"errmsg,omitempty"`
	Details []any  `json:"detail,omitempty"`
}

func (r *Result) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *Result) Detail(args ...any) *Result {
	for _, v := range args {
		r.Details = append(r.Details, v)
	}
	return r
}

func (r *Result) SetCode(code int) *Result {
	r.Code = code
	r.Errmsg = errTexts[code]
	return r
}

func ResultError(code int) *Result {
	return &Result{
		Code:   code,
		Errmsg: errTexts[code],
	}
}

func ResultSuccess() *Result {
	code := 0
	return &Result{
		Code:   code,
		Errmsg: errTexts[code],
	}
}

type CrontabLog struct {
	Total int `json:"total"`
	New   int `json:"new"`
}

func (cl *CrontabLog) String() string {
	b, _ := json.Marshal(cl)
	return string(b)
}

func NewCrontabLog() *CrontabLog {
	return &CrontabLog{}
}

func (cl *CrontabLog) SetNew(v int) *CrontabLog {
	cl.New = v
	return cl
}

func (cl *CrontabLog) SetTotal(v int) *CrontabLog {
	cl.Total = v
	return cl
}

type UpdateLog struct {
	New       int    `json:"new"`
	Delete    int    `json:"delete"`
	DeleteAct int    `json:"delete_act"`
	NewAct    int    `json:"new_act"`
	Errmsg    string `json:"errmsg,omitempty"`
}

func NewUpdateLog() *UpdateLog {
	return &UpdateLog{}
}

func (ul *UpdateLog) String() string {
	b, _ := json.Marshal(ul)
	return string(b)
}

func (ul *UpdateLog) SetNew(v int) *UpdateLog {
	ul.New = v
	return ul
}

func (ul *UpdateLog) SetNewAct(v int) *UpdateLog {
	ul.NewAct = v
	return ul
}

func (ul *UpdateLog) SetDelete(v int) *UpdateLog {
	ul.Delete = v
	return ul
}

func (ul *UpdateLog) SetDeleteAct(v int) *UpdateLog {
	ul.DeleteAct = v
	return ul
}

func (ul *UpdateLog) SetErrmsg(v string) *UpdateLog {
	ul.Errmsg = v
	return ul
}
