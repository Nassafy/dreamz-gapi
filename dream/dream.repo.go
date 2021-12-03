package dream

import (
	"log"
	"time"

	"dreamz.com/api/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

func dbGetDreamDays(store *common.Store, userId string) []DreamDay {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.M{"date": -1})

	cur, err := getDreamCollection(store).Find(ctx, bson.M{"userId": userId}, &queryOptions)
	if err != nil {
		log.Fatal("Error retrieving dreamDays: ", err)
	}
	var dreamDays []DreamDay
	if err = cur.All(ctx, &dreamDays); err != nil {
		log.Fatal(err)
	}
	if dreamDays == nil {
		dreamDays = []DreamDay{}
	}
	return dreamDays
}

func dbGetTodayDream(store *common.Store, userId string) *DreamDay {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	endTime := startTime.Add(time.Hour * 24)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := getDreamCollection(store).FindOne(ctx, bson.M{"userId": userId, "date": bson.M{"$gte": startTime, "$lte": endTime}})

	var dreamDay DreamDay
	if err := result.Decode(&dreamDay); err != nil {
		return nil
	}
	return &dreamDay
}

func dbUpdateDreamDay(store *common.Store, dream *DreamDay) *DreamDay {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var oDream DreamDay
	d := getDreamCollection(store).FindOne(ctx, bson.M{"id": dream.Id, "userId": dream.UserId})
	err := d.Decode(&oDream)

	if err != nil {
		oDream = *dream
	}

	dream.Id = oDream.Id
	dream.Date = oDream.Date
	dream.handleDefault()

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	insered := getDreamCollection(store).FindOneAndUpdate(ctx, bson.M{"id": dream.Id, "userId": dream.UserId}, bson.M{"$set": dream}, &opt)
	var nDream DreamDay
	err = insered.Decode(&nDream)
	if err != nil {
		log.Panic("error in update dream: ", err)
	}
	return &nDream
}

func DbDeleteDreamDay(store *common.Store, id string, userId string) int64 {
	ctx, cancel := common.GetContext()
	defer cancel()
	res, err := getDreamCollection(store).DeleteOne(ctx, bson.M{"id": id, "userId": userId})
	if err != nil {
		log.Panic(err)
	}
	return res.DeletedCount
}

func getDreamCollection(store *common.Store) *mongo.Collection {
	return store.Client.Database("dreamz").Collection("dreamdays")
}
