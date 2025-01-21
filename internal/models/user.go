package models

type user struct {
	ID       string `json:"id" bson:"_id"`
	UserName string `json:"username" bson:"_username"`
	Password string `json:"password" bson:"_password"`
}
