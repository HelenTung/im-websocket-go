package module

import "context"

type RoomBasic struct {
	Identity     string `bson:"identity,omitempty"`
	Number       string `bson:"number,omitempty"`
	Name         string `bson:"name,omitempty"`
	Info         string `bson:"info,omitempty"`
	UserIdentity string `bson:"user_identity,omitempty"`
	CreatAt      int64  `bson:"created_at,omitempty"`
	UpdatedAt    int64  `bson:"updated_at,omitempty"`
}

func (RoomBasic) CollectionName() string {
	return "room_basic"
}

func InsertOneRoomBasic(r *RoomBasic) error {
	_, err := Mongo.Collection(RoomBasic{}.CollectionName()).InsertOne(context.Background(), r)
	if err != nil {
		return err
	}
	return nil
}
