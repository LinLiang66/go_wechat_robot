package routers

import (
	"github.com/gin-gonic/gin"
	"go_wechat_robot/controller"
)

// RegisterRouter 路由设置
func RegisterRouter(router *gin.Engine) {
	routerUser(router)
}

// 用户路由
func routerUser(engine *gin.Engine) {
	con := &controller.WebhookController{}
	// 添加新的路由
	wxApi := engine.Group("/wxapi")
	{
		wxApi.GET("/msg/:appid", con.WechatCheck)  //微信公众号初始化验证
		wxApi.POST("/msg/:appid", con.WechatEvent) //接收微信公众号消息
		wxApi.POST("/save/:appid", con.WechatSave) //微信公众号凭证入库redis
	}
}
