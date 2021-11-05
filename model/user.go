package model

type User struct {
	Id       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

func (user *User) HandleDefault() {}
