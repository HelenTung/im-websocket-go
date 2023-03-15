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
	fmt.Println(Useridentity, RoomIdentity, ur)
	return ur, err
}

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
