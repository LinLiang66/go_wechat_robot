package wechat

import (
	"go_wechat_robot/model"
	"go_wechat_robot/redisService"
	"strconv"
)

type Action interface {
	Execute(appid string, appCache model.WechatAppCache, event *Msg) bool
}

// ProcessedUniqueAction 消息去重处理
type ProcessedUniqueAction struct { //幂等判断消息唯一性
}

func (*ProcessedUniqueAction) Execute(appid string, appCache model.WechatAppCache, event *Msg) bool {
	if redisService.RedisClient.KEYEXISTS("robot:wechat_event:" + strconv.FormatInt(event.MsgId, 10)) {
		return false
	}
	redisService.RedisClient.SetStrWithExpire("robot:wechat_event:"+strconv.FormatInt(event.MsgId, 10), "Message has been handle", 25200)
	return true
}

// EmptyAction 空内容处理
type EmptyAction struct { /*空消息*/
}

func (*EmptyAction) Execute(appid string, appCache model.WechatAppCache, event *Msg) bool {
	if len(event.Content) == 0 {
		tokenCache, tokenExist := redisService.GetWechatToken(appid, appCache)
		if tokenExist {
			go SendKFMsg(tokenCache.AccessToken, model.MsgEntity{Touser: event.FromUserName, Msgtype: "text", Text: model.Text{Content: "🤖️：您好，请问有什么可以帮到您~"}})
		}
		return false
	}
	return true
}

// RobotAction   机器人兜底处理
type RobotAction struct { /*大模型兜底处理*/
}

func (*RobotAction) Execute(appid string, appCache model.WechatAppCache, event *Msg) bool {
	go SendChatGptMessage(appid, appCache, event)
	return false
}
