package define

//var MailPassword = os.Getenv("MailPassword")

var MailPassword = "fdzfhynnasqutkoa"

// var MailPassword = "gqvityxmdrgiuoip"
var RegisterPrefix = "TOKEN_"
var ExpireTime = 300

// 定义消息结构
type MessageStruct struct {
	Message      string `json:"message"`
	RoomIdentity string `json:"room_identity"`
}
