package service

import (
	"context"
	"fmt"
	"time"

	"github.com/joycastle/casual-server-lib/config"
	"github.com/joycastle/matching-story-robot-service/api"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/spf13/cast"
	"google.golang.org/grpc"
)

func SendChatMessage(roomID int64, userInfo model.User, content string) error {
	rpcHost := config.Grpc["chat"]
	fmt.Println(rpcHost)
	conn, err := grpc.Dial(rpcHost, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := api.NewPigeonGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := api.RpcSendReq{Param: &api.RpcMessageParam{
		Route:        "receive",
		IsPersist:    true,
		IgnoreMyself: true,
		User:         &api.User{UserId: userInfo.UserID},
		Room:         &api.Room{RoomId: roomID},
	},
		Body: &api.Message{Id: lib.Generate().Int64(),
			From:     userInfo.UserID,
			Name:     userInfo.UserName,
			HeadIcon: userInfo.UserHeadIcon,
			Type:     "txt",
			Content:  cast.ToString(content),
			Time:     time.Now().Unix()},
	}

	if rsp, err := c.Send(ctx, &req); err != nil {
		return err
	} else {
		fmt.Println(rsp)
	}

	return nil
}

func SendJoinRoomMessage(roomID int64, userInfo model.User) error {
	rpcHost := config.Grpc["chat"]
	fmt.Println(rpcHost)
	conn, err := grpc.Dial(rpcHost, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := api.NewPigeonGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := api.RpcSendReq{Param: &api.RpcMessageParam{
		Route:        "receive",
		IsPersist:    true,
		IgnoreMyself: true,
		User:         &api.User{UserId: userInfo.UserID},
		Room:         &api.Room{RoomId: roomID},
	},
		Body: &api.Message{Id: lib.Generate().Int64(),
			From:     userInfo.UserID,
			Name:     userInfo.UserName,
			HeadIcon: userInfo.UserHeadIcon,
			Type:     "join",
			Time:     time.Now().Unix()},
	}

	if rsp, err := c.Send(ctx, &req); err != nil {
		return err
	} else {
		fmt.Println(rsp)
	}

	return nil
}
