package define

//var MailPassword = os.Getenv("MailPassword")

var MailPassword = "fdzfhynnasqutkoa"

//定义消息结构

type MessageStruct struct {
	Message      string `json:"message"`
	RoomIdentity string `json:"room_identity"`
}
