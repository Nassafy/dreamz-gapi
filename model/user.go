package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

func (user *User) HandleDefault() {
	if user.Id == "" {
		user.Id = uuid.NewV4().String()
	}
}
