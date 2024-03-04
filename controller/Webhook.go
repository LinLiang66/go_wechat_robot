package controller

import (
	"github.com/gin-gonic/gin"
	"go_wechat_robot/model"
	"go_wechat_robot/redisService"
	"go_wechat_robot/utils"
	"go_wechat_robot/wechat"
	"io"
	"log"
	"strconv"
	"time"
)

type WebhookController struct {
}

// WechatSave 用于公众号凭证入库 redis 如果需要长期持久化的话可以自行配置 redis持久化
// 也可以搭配 Mysql 进行持久化处理 每次启动服务之前 去Mysql中初始化数据到 redis 即可
func (controller *WebhookController) WechatSave(c *gin.Context) {
	appid := c.Param("appid")
	c.Header("Server", "Go-Gin-Server")
	if len(appid) == 0 {
		c.XML(400, utils.H{
			"message":   "appid Cannot be empty!!",
			"code":      400,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		return
	}
	Jsonstr, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(200, gin.H{
			"message":   "Failed to read the requested content.",
			"code":      -1,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		return
	}

	//微信公众号对应的app凭证入库redis
	//对应 JSON 结构如下：
	//{"appid":appid,"app_secret":app_secret,"token":token,"aes_key",aes_key}
	err = redisService.RedisClient.SetStr("robot:wechat_app_key:"+appid, string(Jsonstr))

	if err != nil {
		c.JSON(200, gin.H{
			"message":   "Failed to put data into redis storage.",
			"code":      -1,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		return
	}
	c.JSON(200, gin.H{
		"message":   "success",
		"code":      200,
		"success":   true,
		"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
	})
	return
}

// WechatCheck 用于公众号自动验证 https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
// 详情请参考消息体签名及加解密部分的文档 https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Before_Develop/Message_encryption_and_decryption.html
func (controller *WebhookController) WechatCheck(c *gin.Context) {
	appid := c.Param("appid")
	c.Header("Server", "Go-Gin-Server")
	if len(appid) == 0 {
		c.XML(400, utils.H{
			"message":   "appid Cannot be empty!!",
			"code":      400,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		return
	}
	appCache, exist := redisService.GetWechatCache(appid)
	if !exist {
		c.XML(400, utils.H{
			"message":   "appid is invalid",
			"code":      400,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		return
	} else {
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echostr := c.Query("echostr")
		// 校验
		if wechat.CheckSignature(signature, timestamp, nonce, appCache.Token) {
			c.String(200, echostr)
			return
		}
		c.XML(400, utils.H{
			"message":   "此接口为公众号验证，不应该被手动调用，公众号接入校验失败",
			"code":      400,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		log.Println("此接口为公众号验证，不应该被手动调用，公众号接入校验失败")
	}
}

// WechatEvent https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Passive_user_reply_message.html
// 微信服务器在五秒内收不到响应会断掉连接，并且重新发起请求，总共重试三次
func (controller *WebhookController) WechatEvent(c *gin.Context) {
	appid := c.Param("appid")
	c.Header("Server", "Go-Gin-Server")
	if len(appid) == 0 {
		c.XML(400, utils.H{
			"message":   "appid Cannot be empty!!",
			"code":      400,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		return
	}
	appCache, exist := redisService.GetWechatCache(appid)
	if !exist {
		c.XML(400, utils.H{
			"message":   "appid is invalid",
			"code":      400,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
	} else {
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		encType := c.Query("encrypt_type")
		// 校验
		if !wechat.CheckSignature(signature, timestamp, nonce, appCache.Token) {
			c.XML(400, utils.H{
				"message":   "WeChat official account message interface in xml format, please do not call it manually! !",
				"code":      400,
				"success":   false,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		bs, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.XML(400, utils.H{
				"message":   "Error reading request body!!",
				"code":      400,
				"success":   false,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		msg := wechat.NewMsg(bs)
		if msg == nil {
			c.XML(400, utils.H{
				"message":   "WeChat official account message interface in xml format, please do not call it manually! !",
				"code":      400,
				"success":   false,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		if encType == "aes" {
			_, rawMsgXMLByte, err := utils.DecryptMsg(appid, msg.Encrypt, appCache.AesKey)
			if err != nil {
				c.XML(400, utils.H{
					"message":   err.Error(),
					"code":      400,
					"success":   false,
					"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
				})
				log.Printf(err.Error())
				return
			}
			msg = wechat.NewMsg(rawMsgXMLByte)
		}
		//直接先返回 success  避免超时未响应导致微信服务器重推事件
		c.String(200, "success")
		//继续处理其他耗时操作
		if msg.MsgType != "event" {
			if redisService.RedisClient.KEYEXISTS("robot:wechat_event:" + strconv.FormatInt(msg.MsgId, 10)) {
				return
			}
		}
		switch msg.MsgType {
		case "event":
			switch msg.Event {
			case "subscribe": //关注事件
				tokenCache, tokenExist := redisService.GetWechatToken(appid, appCache)
				if tokenExist {
					go wechat.SendKFMsg(tokenCache.AccessToken, model.MsgEntity{Touser: msg.FromUserName, Msgtype: "text", Text: model.Text{Content: "哎呀，终于等到你了，欢迎关注！"}})
				}
				return
			case "unsubscribe": //取消关注事件
				return
			case "CLICK": //点击菜单事件
				tokenCache, tokenExist := redisService.GetWechatToken(appid, appCache)
				if tokenExist {
					go wechat.SendKFMsg(tokenCache.AccessToken, model.MsgEntity{Touser: msg.FromUserName, Msgtype: "text", Text: model.Text{Content: "哎呀，实在抱歉，当前菜单已失效"}})
				}
				return
			default:
				log.Printf(string(bs))
				log.Printf("未实现的事件%s", msg.Event)
				return
			}

		case "text": //文本消息处理
			go func() {
				err := wechat.Handler(appid, appCache, msg)
				if err != nil {
					log.Printf("消息处理责任链报错%s\n", err)
				}
			}()
			return
		// 未写的类型
		default:
			go redisService.RedisClient.SetStrWithExpire("robot:wechat_event:"+strconv.FormatInt(msg.MsgId, 10), "Message has been handle", 25200)
			tokenCache, tokenExist := redisService.GetWechatToken(appid, appCache)
			if tokenExist {
				go wechat.SendKFMsg(tokenCache.AccessToken, model.MsgEntity{Touser: msg.FromUserName, Msgtype: "text", Text: model.Text{Content: "🤖🤖️：亲，暂不支持当前消息类型，请发送纯文本消息哈~"}})
			}
			log.Printf("未实现的消息类型%s\n", msg.MsgType)
			return
		}
	}
}
