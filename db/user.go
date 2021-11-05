package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"dreamz.com/api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsers(store *Store) []model.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := store.Client.Database("dreamz").Collection("users").Find(ctx, bson.D{})
	if err != nil {
		log.Panic("Error retrieving user: ", err)
	}
	var users []model.User
	if err = cur.All(ctx, &users); err != nil {
		log.Panic(err)
	}
	return users
}

func GetUserByUsername(store *Store, username string) model.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user model.User
	store.getUserCollection().FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return user
}

func GetUserById(store *Store, id string) model.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user model.User
	store.getUserCollection().FindOne(ctx, bson.M{"id": id}).Decode(&user)
	return user
}

func UpdateUser(store *Store, user *model.User) *model.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d := store.getUserCollection().FindOne(ctx, bson.M{"id": user.Username})

	fmt.Print(d)
	var oUser model.User
	err := d.Decode(&oUser)
	if err != nil {
		log.Println("Error decoding old user in update: ", err)
		oUser = *user
	}

	user.Id = oUser.Id

	user.HandleDefault()

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	insered := store.getUserCollection().FindOneAndUpdate(ctx, bson.M{"username": user.Username}, bson.M{"$set": user}, &opt)
	var nUser model.User
	err = insered.Decode(&nUser)
	if err != nil {
		log.Print("update user: ", err)
	}
	return &nUser
}

func (store *Store) getUserCollection() *mongo.Collection {
	return store.Client.Database("dreamz").Collection("users")
}
