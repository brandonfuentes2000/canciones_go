package storage

import (
	"canciones/internal/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clienteOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(ctx, clienteOptions)
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	log.Println("conexion a mongodb")
	Client = client
	return nil
}

func GetCollection(database, collection string) *mongo.Collection {
	if Client == nil {
		log.Fatal("MongoDB no esta inicializado")
	}
	return Client.Database(database).Collection(collection)
}

func SaveSong(song models.Song) error {
	collection := GetCollection("songs_db", "songs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": song.Name, "artist": song.Artist, "Album": song.Album}
	update := bson.M{"$set": song}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error al guardar la canción '%s' de '%s': %s", song.Name, song.Artist, err.Error())
	}
	return err
}
