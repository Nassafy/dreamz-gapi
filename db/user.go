package db

import (
	"context"
	"log"
	"time"

	"dreamz.com/api/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers(store *Store) []model.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := store.Client.Database("dreamz").Collection("users").Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Error retrieving user: ", err)
	}
	var users []model.User
	if err = cur.All(ctx, &users); err != nil {
		log.Fatal(err)
	}
	return users
}

func GetUser(store *Store, username string) model.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user model.User
	store.Client.Database("dreamz").Collection("users").FindOne(ctx, bson.D{{"username", username}}).Decode(&user)
	return user
}
