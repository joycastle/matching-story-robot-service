package lib

import (
	"fmt"
	"net/url"
	"testing"
)

type Request struct {
	Code int
	Msg  string
	Data string
}

func Test_HttpPostFormJson(t *testing.T) {

	data := url.Values{}
	data.Set("account_id", "AJSKLDJFJLSJDLFJSLKJFLKJASLJFDKSJF")
	data.Set("user_id", "12121212")
	data.Set("lang_type", "1")
	data.Set("device_type", "9")

	var ret Request

	if err := PostFormToJson("http://127.0.0.1:8081/test", data, &ret); err != nil {
		t.Fatal("http-error", err)
	}

	fmt.Println(ret)
}

func Test_Http1(t *testing.T) {
	data := url.Values{}
	data.Set("req_id", "11")
	data.Set("from", "robot-service")
	data.Set("account_id", "44531B0BCB34A58BBB1D9F92CA5330B7")
	data.Set("user_id", "130290000000")
	data.Set("lang_type", "6")
	data.Set("device_type", "0")
	data.Set("req_help_id", "124317902189887488")
	data.Set("guild_id", "124317585276665856")

	var ret Request

	if err := PostFormToJson("http://127.0.0.1:8081/guild/response", data, &ret); err != nil {
		t.Fatal("http-error", err)
	}

	fmt.Println(ret)
}
