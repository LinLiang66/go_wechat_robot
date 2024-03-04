package wechat

import (
	"encoding/json"
	"fmt"
	"go_wechat_robot/model"
	"io"
	"net/http"
	"strings"
)

var wxUrl = "https://api.weixin.qq.com"

// SendKFMsg  发送客服消息 https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html
func SendKFMsg(accessToken string, msgEntity model.MsgEntity) bool {
	msgbyte, err := json.Marshal(msgEntity)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// 构造请求对象
	req, err := http.NewRequest("POST", wxUrl+"/cgi-bin/message/custom/send?access_token="+accessToken, strings.NewReader(string(msgbyte)))
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	// 发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var res model.RequestEntity
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Printf("请求结果转换为实体: %v\n", resp)
		return false
	}
	return res.Errcode == 0 && res.Errmsg == "ok"
}
