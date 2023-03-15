package module

import "context"

type MessageBasic struct {
	UserIdentity string `bson:"user_identity,omitempty"`
	RoomIdentity string `bson:"room_identity,omitempty"`
	Data         string `bson:"data,omitempty"`
	CreatAt      int64  `bson:"created_at,omitempty"`
	UpdatedAt    int64  `bson:"updated_at,omitempty"`
}

func (MessageBasic) CollectionName() string {
	return "message_basic"
}

func InsertOneMessage(ms *MessageBasic) error {

	_, err := Mongo.Collection(MessageBasic{}.CollectionName()).InsertOne(context.Background(), ms)
	if err != nil {
		return err
	}
	return nil
}
