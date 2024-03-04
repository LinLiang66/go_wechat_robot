package wechat

import "go_wechat_robot/model"

// 责任链
func chain(appid string, appCache model.WechatAppCache, event *Msg, actions ...Action) bool {
	for _, v := range actions {
		if !v.Execute(appid, appCache, event) {
			return false
		}
	}
	return true
}

type MessageHandler struct {
}

func (m MessageHandler) WechatHandler(appid string, appCache model.WechatAppCache, event *Msg) error {

	actions := []Action{
		&ProcessedUniqueAction{}, //避免重复处理
		&EmptyAction{},           //空消息处理
		&RobotAction{},           //大模型兜底处理
	}
	chain(appid, appCache, event, actions...)

	return nil
}
