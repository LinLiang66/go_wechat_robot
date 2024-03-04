package wechat

import (
	"go_wechat_robot/model"
)

type MessageHandlerInterface interface {
	WechatHandler(appid string, appCache model.WechatAppCache, event *Msg) error
}

// wechatHandler 所有消息类型类型的处理器
var wechatHandler MessageHandlerInterface

func InitHandlers() {
	wechatHandler = NewMessageHandler()
}

func Handler(appid string, appCache model.WechatAppCache, event *Msg) error {
	return wechatHandler.WechatHandler(appid, appCache, event)
}

func NewMessageHandler() MessageHandlerInterface {
	return &MessageHandler{}
}
