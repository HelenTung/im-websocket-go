package service

import (
	"main/helper"
	"main/module"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ChatList(c *gin.Context) {
	RoomIdentity := c.Query("room_identity")
	if RoomIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间号不能为空",
		})
		return
	}
	//判断用户是否属于该房间
	uc := c.MustGet("user_claims").(*helper.UserClaims)
	_, err := module.GetUserRoomByIdentity(uc.Identity, RoomIdentity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "此用户不在此房间、非法访问",
		})
		return
	}
	pageIndex, _ := strconv.ParseInt(c.Query("page_index"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_index"), 10, 32)
	skin := (pageIndex - 1) * pageSize
	//chat 查询
	data, err := module.GetMsgListByRoomIdentity(RoomIdentity, &pageSize, &skin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库故障、无法查询" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据加载成功",
		"data": gin.H{
			"list": data,
		},
	})

}
