package service

import (
	"context"
	"fmt"
	"log"
	"main/define"
	"main/helper"
	"main/module"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	accounts := c.PostForm("account")
	password := c.PostForm("password")
	if accounts == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	ub, err := module.GetUserBasicAccountPassword(accounts, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
		})
		return
	}
	token, err := helper.GenerateToken(ub.Account, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "succes",
		"data": gin.H{
			"token": token,
		},
	})

}

func UserDetail(c *gin.Context) {
	u, _ := c.Get("user_claims")
	uc := u.(*helper.UserClaims)
	userBasic, err := module.GetUserBasicIdentity(uc.Identity)
	if err != nil {
		log.Println("[DB ERROR]:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  "数据查询成功",
		"data": userBasic,
	})
}

func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱不能为空",
		})
	}
	cnt, err := module.GetUserBasicEmail(email)
	if err != nil {
		log.Println("[DB ERROR]:", err)
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "当前邮箱已经注册",
		})
		return
	}
	code := helper.GetCode()
	err = helper.SendCode(email, code)
	if err != nil {
		log.Panicln("[ERROR]:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		return
	}
	err = module.Rdb.Set(context.Background(), define.RegisterPrefix+email, code, time.Second*time.Duration(define.ExpireTime)).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "邮箱发送成功",
	})
}

func Register(c *gin.Context) {

	code := c.Query("code")
	email := c.Query("email")
	account := c.Query("account")
	password := c.Query("password")

	fmt.Println("postform", code, email, account, password)
	if code == "" || email == "" || account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "输入为空、请重新输入",
		})
		return
	}
	//判断账号是否唯一
	acc, err := module.GuiceUserBasicAccount(account)
	if err != nil {
		log.Println("[DB ERROR]", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		return
	}
	if acc > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号已被注册、请重新输入",
		})
		return
	}
	//校验验证码是否正确
	r, err := module.Rdb.Get(context.Background(), define.RegisterPrefix+email).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码已经过期、请重新获取",
		})
		return
	}
	if r != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确、请重新输入",
		})
		return
	}

	//注册流程、申请新对象、写入数据库、获取token
	ub := &module.UserBasic{
		Identity:  helper.GetUUID(),
		Account:   account,
		Password:  password,
		Email:     email,
		CreatAt:   time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err = module.InsertOneUserBasic(ub)
	if err != nil {
		log.Println("[DB ERROR]", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		return
	}
	token, err := helper.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"token": token,
		},
	})
}

func UserQuery(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号输入为空请重新输入",
		})
		return
	}
	userBasic, err := module.GetUserBasicAccount(account)
	if err != nil {
		log.Println("[DB ERROR]:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询失败",
		})
		return
	}
	u, _ := c.MustGet("user_claims").(*helper.UserClaims)
	data := &define.UserQueryResult{
		Nickname: userBasic.Nickname,
		Sex:      userBasic.Sex,
		Avatar:   userBasic.Avatar,
		Email:    userBasic.Email,
		IsFriend: false,
	}
	if module.JudgeUserIsFriend(u.Identity, userBasic.Account) {
		data.IsFriend = true
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户信息",
		"data": data,
	})
}

func UserAdd(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号输入为空请重新输入",
		})
		return
	}
	ub, err := module.GetUserBasicAccount(account)
	fmt.Println(ub)
	if err != nil {
		log.Println("[DB ERROR]:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询失败",
		})
		return
	}
	u, _ := c.MustGet("user_claims").(*helper.UserClaims)
	fmt.Println(u)
	if module.JudgeUserIsFriend(u.Identity, ub.Account) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "互为好友、不可重复添加",
		})
		return
	}
	//添加好友、插入对应关系与房间、插入对应房间
	rb := &module.RoomBasic{
		Identity:     helper.GetUUID(),
		UserIdentity: u.Identity,
		CreatAt:      time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		Info:         ub.Account + " and " + u.Identity + " room",
	}
	//保存房间
	if err := module.InsertOneRoomBasic(rb); err != nil {
		log.Println("[DB ERROR]:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间插入失败",
		})
		return
	}
	r := &module.UserRoom{
		UserIdentity: u.Identity,
		RoomIdentity: rb.Identity,
		RoomType:     1,
		CreatAt:      time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	//保存对应关系
	if err := module.InsertOneUserRoom(r); err != nil {
		log.Println("[DB ERROR]:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "对应关系插入失败",
		})
		return
	}
	r1 := &module.UserRoom{
		UserIdentity: ub.Account,
		RoomIdentity: rb.Identity,
		RoomType:     1,
		CreatAt:      time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	//保存对应关系
	if err := module.InsertOneUserRoom(r1); err != nil {
		log.Println("[DB ERROR]:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "对应关系插入失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "好友添加成功",
	})
}

func UserDelete(c *gin.Context) {
	account := c.Query("account")
	fmt.Println(account)
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除空账号好友",
		})
		return
	}
	u := c.MustGet("user_claims").(*helper.UserClaims)
	//获取房间id
	mrs := module.GetUserRoomIdentity(account, u.Identity)
	if mrs == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该账号并非好友、请确认输入是否正确",
		})
		return
	}
	if err := module.DeleUserRoom(mrs); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除对应关系错误、请重试",
		})
		return
	}
	if err := module.DeleRoomBasic(mrs); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除房间错误错误、请重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功删除好友",
	})

}
