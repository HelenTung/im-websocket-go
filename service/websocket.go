package service

import (
	"fmt"
	"log"
	"main/define"
	"main/helper"
	"main/module"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var wc = make(map[string]*websocket.Conn)

func WebsocketMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
		})
		return
	}
	defer conn.Close()
	uc := c.MustGet("user_claims").(*helper.UserClaims)
	wc[uc.Identity] = conn
	for {
		ms := new(define.MessageStruct)
		err := conn.ReadJSON(ms)
		if err != nil {
			log.Printf("Read error:%v\n", err)
			return
		}
		fmt.Println(uc.Identity, ms.RoomIdentity)
		//TODO:判断用户是否属于消息体的房间
		_, err = module.GetUserRoomByIdentity(uc.Identity, ms.RoomIdentity)
		if err != nil {
			log.Printf("UserIdentity:%v RoomIdentity:%v Not Exits\n", uc.Identity, ms.RoomIdentity)
			return
		}
		//TODO:保存消息、消息持久化
		m := &module.MessageBasic{
			UserIdentity: uc.Identity,
			RoomIdentity: ms.RoomIdentity,
			Data:         ms.Message,
			CreatAt:      time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		err = module.InsertOneMessage(m)
		if err != nil {
			log.Println("[DB] Insert Message ERROR", err)
			return
		}
		//TODO:获取在特定房间的在线用户
		urs, err := module.GetUserRoomByRoomIdentity(ms.RoomIdentity)
		if err != nil {
			log.Printf("RoomIdentity:%v Not Exits\n", ms.RoomIdentity)
			return
		}
		for _, room := range urs {
			fmt.Println(room.RoomIdentity, room.UserIdentity)
			if cc, ok := wc[room.UserIdentity]; ok {
				err := cc.WriteMessage(websocket.TextMessage, []byte(ms.Message))
				if err != nil {
					log.Printf("Write error:%v\n", err)
					return
				}
			}
		}
	}
}
