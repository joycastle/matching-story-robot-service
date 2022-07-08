package action

import (
	"encoding/json"
)

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

	//ownai
	1000: "active day is -1 not take effect",
	1001: "not a active day",
	1002: "rule1 match",
	1003: "rule2 match",

	//help
	3000: "help: the result of not completed is empty",
	3001: "help: no new reasonable request",
	3002: "help: normal user request is empty",
	3003: "help: no request to help",

	//request
	4000: "request: not exceeded freeze time",
}

type Result struct {
	Code    int    `json:"code"`
	Errmsg  string `json:"errmsg,omitempty"`
	Details []any  `json:"detail,omitempty"`
}

func (ar *Result) String() string {
	b, _ := json.Marshal(ar)
	return string(b)
}

//Must be an even number
func (ar *Result) Detail(args ...any) *Result {
	for _, v := range args {
		ar.Details = append(ar.Details, v)
	}
	return ar
}

func ErrorText(code int) *Result {
	ar := &Result{
		Code:   code,
		Errmsg: errTexts[code],
	}
	return ar
}

func ActionSuccess() *Result {
	return &Result{
		Code:   0,
		Errmsg: errTexts[0],
	}
}

type ScheduleLog struct {
	Action string  `json:"action"`
	Job    *Job    `json:"job"`
	Msg    string  `json:"msg,omitempty"`
	Result *Result `json:"result,omitempty"`
}

func (sl *ScheduleLog) String() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *ScheduleLog) SetResult(result *Result) *ScheduleLog {
	sl.Result = result
	return sl
}

func (sl *ScheduleLog) SetError(msg string) *ScheduleLog {
	sl.Msg = msg
	return sl
}

func NewScheduleLog(job *Job, action string) *ScheduleLog {
	return &ScheduleLog{Action: action, Job: job}
}

type Job struct {
	GuildID    int64 `json:"guild_id"`
	UserID     int64 `json:"user_id"`
	ActionTime int64 `json:"action_time"`
	RobotNum   int   `json:"robot_num,omitempty"`
}

type CrontabLog struct {
	Action string   `json:"action"`
	Total  int      `json:"total"`
	New    int      `json:"new"`
	Reset  int      `json:"reset"`
	Extra  []string `json:"extra,omitempty"`
}

func (cl *CrontabLog) String() string {
	b, _ := json.Marshal(cl)
	return string(b)
}

func (cl *CrontabLog) SetTotal(v int) *CrontabLog {
	cl.Total = v
	return cl
}

func (cl *CrontabLog) SetNew(v int) *CrontabLog {
	cl.New = v
	return cl
}

func (cl *CrontabLog) SetReset(v int) *CrontabLog {
	cl.Reset = v
	return cl
}

func (cl *CrontabLog) AddExtra(r *Result) *CrontabLog {
	cl.Extra = append(cl.Extra, r.String())
	return cl
}

func NewCrontabLog(action string) *CrontabLog {
	return &CrontabLog{Action: action}
}

type UpdateJobLog struct {
	State  string         `json:"state"`
	Delete map[string]int `json:"delete,omitempty"`
	Create map[string]int `json:"create,omitempty"`
}

func NewUpdateJobLog() *UpdateJobLog {
	return &UpdateJobLog{
		State:  "success",
		Delete: make(map[string]int),
		Create: make(map[string]int),
	}
}

func (ujl *UpdateJobLog) AddDeleteInfo(t string, v int) *UpdateJobLog {
	ujl.Delete[t] = v
	return ujl
}

func (ujl *UpdateJobLog) AddCreateInfo(t string, v int) *UpdateJobLog {
	ujl.Create[t] = v
	return ujl
}

func (ujl *UpdateJobLog) SetState(stat string) *UpdateJobLog {
	ujl.State = stat
	return ujl
}

func (ujl *UpdateJobLog) String() string {
	b, _ := json.Marshal(ujl)
	return string(b)
}
