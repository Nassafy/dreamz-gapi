package db

import (
	"log"
	"time"

	"dreamz.com/api/model"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

func GetDreamDays(store *Store, userId string) []model.DreamDay {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := store.Client.Database("dreamz").Collection("dreamdays").Find(ctx, bson.D{{"userId", userId}})
	if err != nil {
		log.Fatal("Error retrieving dreamDays: ", err)
	}
	var dreamDays []model.DreamDay
	if err = cur.All(ctx, &dreamDays); err != nil {
		log.Fatal(err)
	}
	return dreamDays
}

func NewDreamDay(store *Store, dream *model.DreamDay) primitive.ObjectID {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if dream.Id == "" {
		dream.Id = uuid.NewV4().String()
	}
	dream.HandleDefault()
	insered, err := getCollection(store).InsertOne(ctx, dream)
	if err != nil {
		log.Fatal(err)
	}
	return insered.InsertedID.(primitive.ObjectID)
}

func UpdateDreamDay(store *Store, dream *model.DreamDay, id string) *model.DreamDay {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d := getCollection(store).FindOne(ctx, bson.D{{"id", id}, {"userId", dream.UserId}})

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
	insered := getCollection(store).FindOneAndUpdate(ctx, bson.D{{"id", id}, {"userId", dream.UserId}}, bson.D{{"$set", dream}}, &opt)
	var nDream model.DreamDay
	err = insered.Decode(nDream)
	if err != nil {
		log.Print("update dream: ", err)
	}
	return &nDream

}

func getCollection(store *Store) *mongo.Collection {
	return store.Client.Database("dreamz").Collection("dreamdays")
}
