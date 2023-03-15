package module

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func GetMsgListByRoomIdentity(RoomIdentity string, limit, skip *int64) ([]*MessageBasic, error) {
	data := make([]*MessageBasic, 0)
	cur, err := Mongo.Collection(MessageBasic{}.CollectionName()).
		Find(context.Background(), bson.M{"room_identity": RoomIdentity}, &options.FindOptions{
			Limit: limit,
			Skip:  skip,
			Sort: bson.M{
				"created_at": -1,
			},
		})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		mb := new(MessageBasic)
		err := cur.Decode(mb)
		if err != nil {
			return nil, err
		}
		data = append(data, mb)
	}
	return data, nil
}
