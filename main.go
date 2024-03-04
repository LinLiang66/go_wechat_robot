package main

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"go_wechat_robot/redisService"
	"go_wechat_robot/routers"
	"go_wechat_robot/utils"
	"log"
	"time"
)

func main() {
	//初始化全局Redis服务，后续服务使用redis时  直接使用 redisService即可
	redisService.InitRedisUtil("你的redis连接地址一般为IP", 6379, "redis连接密码")

	// 注册gin处理器 默认开启日志打印
	//g := gin.Default()

	// 注册gin处理器 默认关闭日志打印
	//设置不打印请求路径日志
	g := gin.New()
	g.Use(utils.CustomMiddleware())
	//测试服务状态接口
	g.GET("/ping", func(c *gin.Context) {
		c.Header("Server", "Go-Gin-Server")
		c.JSON(200, gin.H{
			"message":   "pong",
			"code":      200,
			"success":   true,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
	})
	//添加路由
	routers.RegisterRouter(g)

	// 启动服务
	err := g.Run(":80") //
	if err != nil {
		log.Printf("failed to start server: %v", err)
	}

}
