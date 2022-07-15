package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/joycastle/casual-server-lib/config"
)

func SendChatMessageRPC(accountID string, userID int64, guildID int64, chatMsg string) (*MatchingRPCResponse, error) {
	data := url.Values{}

	data.Set("req_id", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	data.Set("from", "robot-service")
	data.Set("account_id", accountID)
	data.Set("user_id", fmt.Sprintf("%d", userID))
	data.Set("lang_type", "1")
	data.Set("device_type", "9")

	data.Set("guild_id", fmt.Sprintf("%d", guildID))
	data.Set("chat_msg", chatMsg)

	rpcHost := config.Grpc["default"]

	resp, err := http.PostForm(rpcHost+"/guild/chat", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &MatchingRPCResponse{}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(err.Error() + string(body))
	}

	if output.Code != 0 {
		return output, errors.New(output.Data + " " + output.Errmsg)
	}

	return output, nil
}

func SendRequestRPC(accountID string, userID int64, guildID int64) (*MatchingRPCResponse, error) {
	data := url.Values{}

	data.Set("req_id", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	data.Set("from", "robot-service")
	data.Set("account_id", accountID)
	data.Set("user_id", fmt.Sprintf("%d", userID))
	data.Set("lang_type", "1")
	data.Set("device_type", "9")

	data.Set("guild_id", fmt.Sprintf("%d", guildID))

	rpcHost := config.Grpc["default"]

	resp, err := http.PostForm(rpcHost+"/guild/request", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &MatchingRPCResponse{}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(err.Error() + string(body))
	}

	if output.Code != 0 {
		return output, errors.New(output.Data + " " + output.Errmsg)
	}

	return output, nil
}

func SendRequestHelpRPC(accountID string, userID int64, guildID int64, helpID int64) (*MatchingRPCResponse, error) {
	data := url.Values{}

	data.Set("req_id", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	data.Set("from", "robot-service")
	data.Set("account_id", accountID)
	data.Set("user_id", fmt.Sprintf("%d", userID))
	data.Set("lang_type", "1")
	data.Set("device_type", "9")

	data.Set("guild_id", fmt.Sprintf("%d", guildID))
	data.Set("req_help_id", fmt.Sprintf("%d", helpID))

	rpcHost := config.Grpc["default"]

	resp, err := http.PostForm(rpcHost+"/guild/response", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &MatchingRPCResponse{}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(err.Error() + string(body))
	}

	if output.Code != 0 {
		return output, errors.New(output.Data + " " + output.Errmsg)
	}

	return output, nil
}

func SendJoinToGuildRPC(accountID string, userID int64, guildID int64) (*MatchingRPCResponse, error) {
	data := url.Values{}

	data.Set("req_id", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	data.Set("from", "robot-service")
	data.Set("account_id", accountID)
	data.Set("user_id", fmt.Sprintf("%d", userID))
	data.Set("lang_type", "1")
	data.Set("device_type", "9")

	data.Set("guild_id", fmt.Sprintf("%d", guildID))

	rpcHost := config.Grpc["default"]

	resp, err := http.PostForm(rpcHost+"/guild/join", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &MatchingRPCResponse{}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(err.Error() + string(body))
	}

	if output.Code != 0 {
		return output, errors.New(output.Data + " " + output.Errmsg)
	}

	return output, nil
}

func SendUpdateScoreRPC(accountID string, userID int64, score int) (*MatchingRPCResponse, error) {
	data := url.Values{}

	data.Set("req_id", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	data.Set("from", "robot-service")
	data.Set("account_id", accountID)
	data.Set("user_id", fmt.Sprintf("%d", userID))
	data.Set("lang_type", "1")
	data.Set("device_type", "9")

	data.Set("level_id", fmt.Sprintf("robot_%d", time.Now().UnixNano()/1000000))
	data.Set("score", fmt.Sprintf("%d", score))

	rpcHost := config.Grpc["default"]

	resp, err := http.PostForm(rpcHost+"/guild/updatescore", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &MatchingRPCResponse{}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(err.Error() + string(body))
	}

	if output.Code != 0 {
		return output, errors.New(output.Data + " " + output.Errmsg)
	}

	return output, nil
}

func SendLeaveGuildRPC(accountID string, userID int64, guildId int64) (*MatchingRPCResponse, error) {
	data := url.Values{}

	if len(accountID) == 0 {
		users, err := GetUserInfosWithField([]int64{userID}, []string{"account_id"})
		if err != nil {
			return nil, err
		}
		if len(users) != 1 {
			return nil, fmt.Errorf("not found user infos user_id:%d", userID)
		}
		accountID = users[0].AccountID
	}

	data.Set("req_id", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	data.Set("from", "robot-service")
	data.Set("account_id", accountID)
	data.Set("user_id", fmt.Sprintf("%d", userID))
	data.Set("lang_type", "1")
	data.Set("device_type", "9")

	data.Set("guild_id", fmt.Sprintf("%d", guildId))

	rpcHost := config.Grpc["default"]

	resp, err := http.PostForm(rpcHost+"/guild/leave", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &MatchingRPCResponse{}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(err.Error() + string(body))
	}

	if output.Code != 0 {
		return output, errors.New(output.Data + " " + output.Errmsg)
	}

	return output, nil
}
