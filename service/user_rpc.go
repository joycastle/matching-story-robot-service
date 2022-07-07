package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/joycastle/casual-server-lib/config"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
)

//创建工会机器人
func CreateGuildRobotUserRPC(name, headIcon string, likeCnt, level int) (model.User, error) {
	var u model.User

	data := url.Values{}
	data.Set("req_id", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	data.Set("user_id", "1001")
	data.Set("lang_type", "0")
	data.Set("from", "robot-service")
	data.Set("account_id", lib.Md5(name))
	data.Set("name", name)
	data.Set("head_icon", headIcon)
	data.Set("country", "CN")
	data.Set("user_type", USERTYPE_CLUB_ROBOT_SERVICE)
	data.Set("device", lib.Md5(name))
	data.Set("device_type", "9")
	data.Set("language", "0")
	data.Set("like_cnt", fmt.Sprintf("%d", likeCnt))
	data.Set("level", fmt.Sprintf("%d", level))
	data.Set("channel", "0")

	rpcHost := config.Grpc["default"]

	resp, err := http.PostForm(rpcHost+"/user/create", data)
	if err != nil {
		return u, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return u, err
	}

	output := &MatchingRPCResponse{}

	if err := json.Unmarshal(body, output); err != nil {
		return u, fmt.Errorf(err.Error() + string(body))
	}

	if output.Code != 0 {
		return u, errors.New(output.Data + " " + output.Errmsg)
	}

	userID, err := strconv.ParseInt(output.Data, 10, 64)
	if err != nil {
		return u, err
	}

	users, err := GetUserInfosWithField([]int64{userID}, []string{"*"})
	if err != nil {
		return u, err
	}
	if len(users) == 0 {
		return u, fmt.Errorf("find error user_id:%d", userID)
	}

	return users[0], nil
}
