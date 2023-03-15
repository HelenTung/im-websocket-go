package service

import (
	"log"
	"main/helper"
	"main/module"
	"net/http"

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
	err = helper.SendCode(email, "3241")
	if err != nil {
		log.Panicln("[ERROR]:", err)
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
