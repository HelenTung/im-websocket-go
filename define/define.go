package define

//var MailPassword = os.Getenv("MailPassword")

// gamil
var MailPassword = "fdzfhynnasqutkoa"

// var MailPassword = "gqvityxmdrgiuoip"

// qq
// var MailPassword = "oyimwoehnurdcjhh"
var RegisterPrefix = "TOKEN_"
var ExpireTime = 300

// 定义消息结构
type MessageStruct struct {
	Message      string `json:"message"`
	RoomIdentity string `json:"room_identity"`
}

type UserQueryResult struct {
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Sex      int    `json:"sex,omitempty"`
	Email    string `json:"email,omitempty"`
	IsFriend bool   `json:"is_friend"`
}
