# im-websocket-go
项目运用了redis、mongodb、websocket

项目的主要功能有：
1、用户通过邮箱注册
2、添加好友
3、删除好友
4、查看聊天记录
5、好友聊天
6、群聊
7、查询用户个人信息
8、查询是否为好友

定义主要结构如下、放在module模块下、其结构的主要方法也放置在对应文件下、恰好对应数据库的中四个collection
```go
type MessageBasic struct {
	UserIdentity string `bson:"user_identity,omitempty"`
	RoomIdentity string `bson:"room_identity,omitempty"`
	Data         string `bson:"data,omitempty"`
	CreatAt      int64  `bson:"created_at,omitempty"`
	UpdatedAt    int64  `bson:"updated_at,omitempty"`
}
type RoomBasic struct {
	Identity     string `bson:"identity,omitempty"`
	Number       string `bson:"number,omitempty"`
	Name         string `bson:"name,omitempty"`
	Info         string `bson:"info,omitempty"`
	UserIdentity string `bson:"user_identity,omitempty"`
	CreatAt      int64  `bson:"created_at,omitempty"`
	UpdatedAt    int64  `bson:"updated_at,omitempty"`
}
type UserBasic struct {
	Identity  string `bson:"identity,omitempty"`
	Account   string `bson:"account,omitempty"`
	Password  string `bson:"password,omitempty"`
	Nickname  string `bson:"nickname,omitempty"`
	CreatAt   int64  `bson:"created_at,omitempty"`
	UpdatedAt int64  `bson:"updated_at,omitempty"`
	Avatar    string `bson:"avatar,omitempty"`
	Sex       int    `bson:"sex,omitempty"`
	Email     string `bson:"email,omitempty"`
}
type UserRoom struct {
	UserIdentity string `bson:"user_identity,omitempty"`
	RoomIdentity string `bson:"room_identity,omitempty"`
	RoomType     int64  `bson:"room_type,omitempty"` // 房间类型，1为单独房间，其他为群聊
	CreatAt      int64  `bson:"created_at,omitempty"`
	UpdatedAt    int64  `bson:"updated_at,omitempty"`
}
```

其中运用了middleware、用来做身份检测
```go
func AuthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		user_claims, err := helper.AnalyseToken(token)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户认证不通过",
			})
			return
		}
		ctx.Set("user_claims", user_claims)
		ctx.Next()
	}
}
```


第三方包为
```go
"github.com/gin-gonic/gin"
"github.com/gorilla/websocket"
"github.com/dgrijalva/jwt-go"
"github.com/gofrs/uuid"
"github.com/jordan-wright/email"
```
