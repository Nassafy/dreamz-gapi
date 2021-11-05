package db

import (
	"log"

	"dreamz.com/api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRefreshToken(store *Store, refreshToken model.RefreshToken) string {
	ctx, cancel := GetContext()
	defer cancel()
	result, err := getRefreshCollection(store).InsertOne(ctx, refreshToken)
	if err != nil {
		log.Panic("Error saving refresh token:", err)
	}
	return result.InsertedID.(primitive.ObjectID).Hex()
}

func InvalidateRefreshToken(store *Store, id string) {
	ctx, cancel := GetContext()
	defer cancel()
	result := getRefreshCollection(store).FindOne(ctx, bson.M{"id": id})
	var refreshToken model.RefreshToken
	err := result.Decode(&refreshToken)
	if err != nil {
		log.Panic("Error decoding old refresh token", err)
	}
	refreshToken.Valid = false
	_, err = getRefreshCollection(store).UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": refreshToken})
	if err != nil {
		log.Panic("Error updating refresh token:", err)
	}
}

func InvalidateUserRefreshToken(store *Store, userId string) int64 {
	ctx, cancel := GetContext()
	defer cancel()
	result, err := getRefreshCollection(store).UpdateMany(
		ctx, bson.M{"userId": userId, "valid": false},
		bson.M{"$set": bson.M{"valid": false}},
	)
	if err != nil {
		log.Panic("Error invalidating token: ", err)
	}
	return result.MatchedCount
}

func GetRefreshToken(store *Store, tokenId string, userId string) *model.RefreshToken {
	ctx, cancel := GetContext()
	defer cancel()
	result := getRefreshCollection(store).FindOne(ctx, bson.M{"userId": userId, "id": tokenId})
	var token *model.RefreshToken
	err := result.Decode(&token)
	if err != nil {
		log.Panic("Error decoding token: ", err)
	}
	return token
}

func getRefreshCollection(store *Store) *mongo.Collection {
	return store.Client.Database("dreamz").Collection("refresh")
}
