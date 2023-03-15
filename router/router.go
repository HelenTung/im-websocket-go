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

	//用户注册
	r.POST("/register", service.Register)
	auth := r.Group("/u", middleware.AuthCheck())
	//用户获取状态详情
	auth.GET("/user/detail", service.UserDetail)
	//查询用户个人信息
	auth.POST("/user/query", service.UserQuery)
	//发送接收消息
	auth.GET("/websocket/message", service.WebsocketMessage)
	//聊天记录列表
	auth.GET("/chat/list", service.ChatList)
	return r
}
