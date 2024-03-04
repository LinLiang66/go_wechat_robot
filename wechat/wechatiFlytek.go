package wechat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go_wechat_robot/model"
	"go_wechat_robot/redisService"
	"go_wechat_robot/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func iFlytekSendmessage(a *Msg) string {
	userId := a.FromUserName
	message := a.Content
	appCache, _ := redisService.GetUserCache(userId)
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(assembleAuthUrl1(appCache.RobotSparkUrl, appCache.RobotApiKey, appCache.RobotApiSecret), nil)
	if err != nil {
		panic(readResp(resp) + err.Error())
		return "哎呀，不好意思，这个问题太难了,一会儿再试试"
	} else if resp.StatusCode != 101 {
		panic(readResp(resp) + err.Error())
		return "哎呀，不好意思，这个问题太难了,一会儿再试试"
	}
	go func() {
		data := genParams1(userId, message, appCache)
		conn.WriteJSON(data)
	}()

	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}
		var data map[string]interface{}
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			fmt.Println("Error parsing JSON:", err)

			conn.Close()
			return "处理出错了，请稍后再试"
		}
		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			fmt.Println("讯飞调用报错", data)
			return "处理出错了，请稍后再试"
		}

		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			answer += content

		} else {
			answer += content
			conn.Close()
			break
		}

	}
	return answer
}

// 生成参数
func genParams1(userId string, question string, appCache model.AppCache) map[string]interface{} { // 根据实际情况修改返回的数据结构和字段名
	messages := redisService.GetMessageContext(userId)
	messages = append(messages, model.IFlytekMessage{Role: "user", Content: question})
	newMessage := utils.Checklen(messages)
	redisService.SetMessageContext(userId, newMessage)
	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": appCache.RobotAppID, // 根据实际情况修改返回的数据结构和字段名
			"uid":    userId,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      utils.GetDomainType(appCache.RobotDomain),
				"temperature": appCache.RobotTemperature,
				"top_k":       int64(6),
				"max_tokens":  int64(2048),
				"auditing":    "default",
				"chat_id":     userId,
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": newMessage, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}
