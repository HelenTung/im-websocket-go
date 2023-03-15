package router

import (
	"main/middleware"
	"main/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	//login
	r.POST("/login", service.Login)
	//send code
	r.POST("/send/code", service.SendCode)
	auth := r.Group("/u", middleware.AuthCheck())
	//用户获取状态详情
	auth.GET("/user/detail", service.UserDetail)

	//发送接收消息
	auth.GET("/websocket/message", service.WebsocketMessage)
	return r
}
