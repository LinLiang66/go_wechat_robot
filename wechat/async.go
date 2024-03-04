package wechat

import (
	"go_wechat_robot/model"
	"go_wechat_robot/redisService"
)

func SendChatGptMessage(appid string, appCache model.WechatAppCache, event *Msg) {
	Content := iFlytekSendmessage(event)
	tokenCache, tokenExist := redisService.GetWechatToken(appid, appCache)
	if tokenExist {
		SendKFMsg(tokenCache.AccessToken, model.MsgEntity{Touser: event.FromUserName, Msgtype: "text", Text: model.Text{Content: Content}})
	}

}
