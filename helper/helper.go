package helper

import (
	"crypto/md5"
	"fmt"
	"main/define"
	"math/rand"
	"net/smtp"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jordan-wright/email"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GetMd5
// 生成 md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("gin-gorm-oj-key")

// GenerateToken
// 生成 token
func GenerateToken(identity, email string) (string, error) {

	UserClaim := &UserClaims{
		Identity:       identity,
		Email:          email,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

// SendCode
// 发生验证码
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "<denghailun1635161916@gmail.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已经发送，请注意查收"
	e.HTML = []byte("您的验证码:<b>" + code + "</b>")
	// return e.SendWithStartTLS("smtp.gmail.com:587",
	// 	smtp.PlainAuth("", "denghailun1635161916@gmail.com", define.MailPassword, "smtp.gmail.com"),
	// 	&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.gmail.com"})
	return e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "denghailun1635161916@gmail.com", define.MailPassword, "smtp.gmail.com"))
}

// GetCode
// 生存验证码
func GetCode() string {
	res := ""
	for i := 0; i < 6; i++ {
		res += strconv.Itoa(rand.Intn(10))
	}
	//fmt.Println(res)
	return res
}

// GetUUID
// 生成唯一标识码
func GetUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	fmt.Println(u)
	return fmt.Sprintf("%x", u)
}
