package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Song struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Artist   string             `bson:"artist" json:"artist"`
	Duration string             `bson:"duration" json:"duration"`
	Album    string             `bson:"album" json:"album"`
	Artwork  string             `bson:"artwork" json:"artwork"`
	Price    string             `bson:"price" json:"price"`
	Origin   string             `bson:"origin" json:"origin"`
}
