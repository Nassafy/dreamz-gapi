package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"dreamz.com/api/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func dbGetUsers(store *common.Store) []User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := store.Client.Database("dreamz").Collection("users").Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Error retrieving user: ", err)
	}
	var users []User
	if err = cur.All(ctx, &users); err != nil {
		log.Fatal(err)
	}
	return users
}

func DbGetUser(store *common.Store, username string) User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user User
	getUserCollection(store).FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return user
}

func dbUpdateUser(store *common.Store, user *User) *User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d := getUserCollection(store).FindOne(ctx, bson.M{"id": user.Username})

	fmt.Print(d)
	var oUser User
	err := d.Decode(&oUser)
	if err != nil {
		log.Println("Error decoding old user in update: ", err)
		oUser = *user
	}

	user.ID = oUser.ID

	user.HandleDefault()

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	insered := getUserCollection(store).FindOneAndUpdate(ctx, bson.M{"username": user.Username}, bson.M{"$set": user}, &opt)
	var nUser User
	err = insered.Decode(&nUser)
	if err != nil {
		log.Print("update user: ", err)
	}
	return &nUser
}

func getUserCollection(store *common.Store) *mongo.Collection {
	return store.Client.Database("dreamz").Collection("users")
}
