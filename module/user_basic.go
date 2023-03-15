package module

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

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

func (ub UserBasic) CollectionName() string {
	return "user_basic"
}

func GetUserBasicAccountPassword(account, password string) (*UserBasic, error) {
	ub := new(UserBasic)
	filter := bson.M{"account": account, "password": password}
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), filter).
		Decode(ub)
	fmt.Println(account, password, ub.Account, ub.Password)
	return ub, err
}

func GetUserBasicIdentity(id string) (*UserBasic, error) {
	ub := new(UserBasic)
	filter := bson.M{"identity": id}
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), filter).
		Decode(ub)
	return ub, err
}

func GetUserBasicEmail(email string) (int64, error) {
	return Mongo.Collection(UserBasic{}.CollectionName()).
		CountDocuments(context.Background(), bson.M{"email": email})
}
