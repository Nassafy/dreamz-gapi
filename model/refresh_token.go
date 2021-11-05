package model

type RefreshToken struct {
	Id       string `bson:"id"`
	Token    string `bson:"token"`
	Valid    bool   `bson:"valid"`
	ExpireAt int64  `bson:"exp"`
	UserId   string `bson:"userId"`
}
