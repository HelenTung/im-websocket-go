package test

import (
	"context"
	"fmt"
	"main/module"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestFindOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")
	user := &module.UserBasic{
		Account:   "ArchLinux",
		Password:  "123456",
		Nickname:  "helen",
		CreatAt:   1,
		UpdatedAt: 1,
		Avatar:    "头像",
		Sex:       0,
		Email:     "1635161916@qq.com",
	}
	res, err := db.Collection("user_basic").InsertOne(ctx, user)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}
	ub := new(module.UserBasic)
	filter := bson.D{}
	if err := db.Collection("user_basic").FindOne(context.Background(), filter).Decode(ub); err != nil {
		t.Fatal(err)
	}
	fmt.Println("ub====>", ub)
}

func TestFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")
	user_room := &module.UserRoom{
		UserIdentity: "1",
		RoomIdentity: "2",
		CreatAt:      1,
		UpdatedAt:    1,
	}
	for i := 0; i < 3; i++ {
		res, err := db.Collection("user_room").InsertOne(ctx, user_room)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(res)
		}
	}
	urs := make([]*module.UserRoom, 0)
	cur, err := db.Collection("user_room").Find(context.Background(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	for cur.Next(context.Background()) {
		ur := new(module.UserRoom)
		err := cur.Decode(ur)
		if err != nil {
			t.Fatal(err)
		}
		urs = append(urs, ur)
	}
	for _, v := range urs {
		fmt.Println("user_room======>", v)
	}
}

func TestUserRoomFind(t *testing.T) {
	ur := new(module.UserRoom)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")
	db.Collection("user_room").
		FindOne(context.Background(), bson.M{"user_identity": "1", "room_identity": "helen"}).
		Decode(ur)
	fmt.Println(ur)
}

func TestInsertFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")
	u := &module.UserBasic{
		Account:   "Arch",
		Password:  "123456",
		Nickname:  "helen",
		CreatAt:   1,
		UpdatedAt: 1,
		Avatar:    "头像",
		Sex:       0,
		Email:     "6545698@qq.com",
	}
	res, err := db.Collection("user_basic").InsertOne(ctx, u)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}
	user := new(module.UserBasic)
	for i := 0; i < 3; i++ {
		db.Collection("user_basic").
			FindOne(context.Background(), bson.M{"account": "ArchLinux", "password": "123456"}).
			Decode(user)
	}
	fmt.Println(user)
}
