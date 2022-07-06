package qa

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/joycastle/casual-server-lib/log"
)

type QaKv struct {
	Value       string
	Type        string
	Description string
}

type QaDebug struct {
	list map[string]QaKv
	sort []string
	mu   *sync.Mutex
}

func NewQaDebug() *QaDebug {
	return &QaDebug{list: make(map[string]QaKv), sort: []string{}, mu: new(sync.Mutex)}
}

var IsOpenQaDebug bool = false
var qaDebug *QaDebug = NewQaDebug()

func (qd *QaDebug) AddInt(key string, desc string) {
	qd.mu.Lock()
	defer qd.mu.Unlock()
	qakv := QaKv{
		Value:       "\"not set\"",
		Type:        "int",
		Description: desc,
	}
	qd.list[key] = qakv
	qd.sort = append(qd.sort, key)
}

func (qd *QaDebug) AddString(key string, desc string) {
	qd.mu.Lock()
	defer qd.mu.Unlock()
	qakv := QaKv{
		Value:       "\"not set\"",
		Type:        "string",
		Description: desc,
	}
	qd.list[key] = qakv
	qd.sort = append(qd.sort, key)
}

func (qd *QaDebug) Update(key string, value string) error {
	qd.mu.Lock()
	defer qd.mu.Unlock()

	if v, ok := qd.list[key]; !ok {
		return fmt.Errorf("服务端没有配置该变量的使用，设置无效")
	} else {
		if v.Type == "int" {
			vals := strings.Split(value, ",")
			for _, val := range vals {
				if !IsAllNumber(val) {
					return fmt.Errorf("不能包含非数字字符，设置无效")
				}
			}
		}
		v.Value = value
		qd.list[key] = v
	}
	return nil
}

func (qd *QaDebug) Get(key string) string {
	qd.mu.Lock()
	defer qd.mu.Unlock()
	return qd.list[key].Value
}

func (qd *QaDebug) GetOptions() string {
	qd.mu.Lock()
	defer qd.mu.Unlock()

	result := ""
	if IsOpenQaDebug {
		result += "当前调试功能已打开......\n"
	} else {
		result += "当前调试功能已关闭......\n"
	}
	result += "\n选项:\n"
	for index, key := range qd.sort {
		result = result + fmt.Sprintf("%-2d Key: %-10s Value: %-20s Description: %-20s Usage: /set?key=%s&value={VALUE}\n", index, key, qd.list[key].Value, qd.list[key].Description, key)
	}

	result = result + "\n\n\n\n\n其他：\n"
	result = result + "1.多个值可用英文逗号【,】分割， 如：/set?key=guild_id&value=9064735890735104,9130434411626496,9194785835319296\n"
	result = result + "2./qaclose 关闭调试功能\n"
	result = result + "3./qaopen 打开调试功能\n"
	return result
}

func StartQaDebugMode(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, qaDebug.GetOptions())
		log.Get("run").Info("QaDebug-SET", "/index")
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if err := qaDebug.Update(values.Get("key"), values.Get("value")); err != nil {
			fmt.Fprintf(w, fmt.Sprintln("错误:", err.Error()))
		} else {
			fmt.Fprintf(w, qaDebug.GetOptions())
		}
		log.Get("run").Info("QaDebug-SET", "key:", values.Get("key"), "value:", values.Get("value"))
	})

	http.HandleFunc("/qaopen", func(w http.ResponseWriter, r *http.Request) {
		IsOpenQaDebug = true
		fmt.Fprintf(w, "已打开调试功能")
		log.Get("run").Info("QaDebug-SET", "qa-open")
	})

	http.HandleFunc("/qaclose", func(w http.ResponseWriter, r *http.Request) {
		IsOpenQaDebug = false
		fmt.Fprintf(w, "已关闭调试功能")
		log.Get("run").Info("QaDebug-SET", "qa-close")
	})

	IsOpenQaDebug = true
	http.ListenAndServe(addr, nil)
}

func OpenQaDebug() bool {
	return IsOpenQaDebug
}

func IsAllNumber(str string) bool {
	filter := make(map[rune]struct{})
	filter['.'] = struct{}{}
	filter['1'] = struct{}{}
	filter['2'] = struct{}{}
	filter['3'] = struct{}{}
	filter['4'] = struct{}{}
	filter['5'] = struct{}{}
	filter['6'] = struct{}{}
	filter['7'] = struct{}{}
	filter['8'] = struct{}{}
	filter['9'] = struct{}{}
	filter['0'] = struct{}{}

	for index, b := range str {
		if index == 0 {
			if b == '-' {
				continue
			}
		}
		if _, ok := filter[b]; !ok {
			return false
		}
	}

	return true
}
