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

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.M{"date": -1})

	cur, err := store.getDreamCollection().Find(ctx, bson.M{"userId": userId}, &queryOptions)
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

func GetTodayDream(store *Store, userId string) *model.DreamDay {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	endTime := startTime.Add(time.Hour * 24)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := store.getDreamCollection().FindOne(ctx, bson.M{"userId": userId, "date": bson.M{"$gte": startTime, "$lte": endTime}})

	var dreamDay model.DreamDay
	if err := result.Decode(&dreamDay); err != nil {
		return nil
	}
	return &dreamDay
}

func UpdateDreamDay(store *Store, dream *model.DreamDay) *model.DreamDay {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var oDream model.DreamDay
	d := store.getDreamCollection().FindOne(ctx, bson.M{"id": dream.Id, "userId": dream.UserId})
	err := d.Decode(&oDream)

	if err == nil {
		oDream = *dream
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
	insered := store.getDreamCollection().FindOneAndUpdate(ctx, bson.M{"id": dream.Id, "userId": dream.UserId}, bson.M{"$set": dream}, &opt)
	var nDream model.DreamDay
	err = insered.Decode(&nDream)
	if err != nil {
		log.Panic("error in update dream: ", err)
	}
	return &nDream

}

func (store *Store) getDreamCollection() *mongo.Collection {
	return store.Client.Database("dreamz").Collection("dreamdays")
}
