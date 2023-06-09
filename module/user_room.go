package module

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRoom struct {
	UserIdentity string `bson:"user_identity,omitempty"`
	RoomIdentity string `bson:"room_identity,omitempty"`
	RoomType     int64  `bson:"room_type,omitempty"` // 房间类型，1为单独房间，其他为群聊
	CreatAt      int64  `bson:"created_at,omitempty"`
	UpdatedAt    int64  `bson:"updated_at,omitempty"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

func GetUserRoomByIdentity(Useridentity, RoomIdentity string) (*UserRoom, error) {
	ur := new(UserRoom)
	err := Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.M{"user_identity": Useridentity, "room_identity": RoomIdentity}).
		Decode(ur)
	return ur, err
}

// 通过房间id获取查找对应关联关系
func GetUserRoomByRoomIdentity(RoomIdentity string) ([]*UserRoom, error) {
	cur, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.M{"room_identity": RoomIdentity})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	urs := make([]*UserRoom, 0)
	for cur.Next(context.Background()) {
		ur := new(UserRoom)
		err := cur.Decode(ur)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		urs = append(urs, ur)
	}
	return urs, nil
}

// 判断是否为好友
func JudgeUserIsFriend(user1, user2 string) bool {
	//查询user1单独房间列表
	cur, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.M{"user_identity": user1, "room_type": 1})
	curs := make([]string, 0)
	if err != nil {
		log.Println("[DB ERROR]", err)
		return false
	}
	for cur.Next(context.Background()) {
		ur := new(UserRoom)
		err := cur.Decode(ur)
		if err != nil {
			return false
		}
		curs = append(curs, ur.RoomIdentity)
	}
	//查询user2有多少关联的单独房间
	cur2, err := Mongo.Collection(UserRoom{}.CollectionName()).
		CountDocuments(context.Background(), bson.M{"user_identity": user2, "room_identity": bson.M{"$in": curs}, "room_type": 1})
	if err != nil {
		log.Println("[DB ERROR]", err)
		return false
	}
	if cur2 > 0 {
		return true
	}
	return false
}

// 根据user查找roomid
func GetUserRoomIdentity(user1, user2 string) []string {
	//查询user1单独房间列表
	cur, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.M{"user_identity": user1, "room_type": 1})
	curs := make([]string, 0)
	if err != nil {
		log.Println("[DB ERROR]", err)
		return nil
	}
	for cur.Next(context.Background()) {
		ur := new(UserRoom)
		err := cur.Decode(ur)
		if err != nil {
			return nil
		}
		curs = append(curs, ur.RoomIdentity)
	}
	slice := make([]string, 0)
	//查询user2有多少关联的单独房间
	cur2, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.M{"user_identity": user2, "room_identity": bson.M{"$in": curs}, "room_type": 1})
	if err != nil {
		log.Println("[DB ERROR]", err)
		return nil
	}
	if cur2.Next(context.Background()) {
		ur := new(UserRoom)
		err := cur2.Decode(ur)
		if err != nil {
			return nil
		}
		slice = append(slice, ur.RoomIdentity)
	}
	if len(slice) == 0 {
		return nil
	}
	fmt.Println(slice)
	return slice
}

// 删除数据
func DeleUserRoom(str []string) error {
	for _, v := range str {
		_, err := Mongo.Collection(UserRoom{}.CollectionName()).
			DeleteMany(context.Background(), bson.M{"room_identity": v})
		if err != nil {
			return err
		}
	}
	return nil
}

// 插入用户与房间对应关系
func InsertOneUserRoom(ur *UserRoom) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).InsertOne(context.Background(), ur)
	if err != nil {
		return err
	}
	return nil
}
