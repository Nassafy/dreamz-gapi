package db

import (
	"log"
	"time"

	"dreamz.com/api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

func GetDreamDays(store *Store, userId string) []model.DreamDay {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := store.Client.Database("dreamz").Collection("dreamdays").Find(ctx, bson.M{"userId": userId})
	if err != nil {
		log.Fatal("Error retrieving dreamDays: ", err)
	}
	var dreamDays []model.DreamDay
	if err = cur.All(ctx, &dreamDays); err != nil {
		log.Fatal(err)
	}
	if dreamDays == nil {
		dreamDays = []model.DreamDay{}
	}
	return dreamDays
}

func UpdateDreamDay(store *Store, dream *model.DreamDay) *model.DreamDay {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d := getCollection(store).FindOne(ctx, bson.M{"id": dream.Id, "userId": dream.UserId})

	var oDream model.DreamDay
	err := d.Decode(&oDream)
	if err != nil {
		log.Fatal("Error decoding old dream in update: ", err)
		return nil
	}

	dream.Id = oDream.Id
	dream.Date = oDream.Date
	dream.HandleDefault()

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	insered := getCollection(store).FindOneAndUpdate(ctx, bson.M{"id": dream.Id, "userId": dream.UserId}, bson.M{"$set": dream}, &opt)
	var nDream model.DreamDay
	err = insered.Decode(&nDream)
	if err != nil {
		log.Print("update dream: ", err)
	}
	return &nDream

}

func getCollection(store *Store) *mongo.Collection {
	return store.Client.Database("dreamz").Collection("dreamdays")
}
