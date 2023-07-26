package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Veri tabanı bağalantı linki ujkbqpruk3Q2cnXo

// Değişkenler
var urlDB mongo.Database
var urlCollection mongo.Collection
var userCollection mongo.Collection

func init() {

	//.env dosyasını yükle
	err := godotenv.Load("enviroments/.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	// Veri tabanına bağlanmak için bir istemci oluştur
	uri := os.Getenv("DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	// Birincil geçikme
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	//Veritabanı ve koleksiyonların değişkenlere tanımlanması
	urlDB = *client.Database("shorturls")
	urlCollection = *urlDB.Collection("url")
	userCollection = *urlDB.Collection("user")

	fmt.Println("Successfully connected and pinged.")
}

// Veritabanı ve koleksiyon bağlatısını diğer dosyalarda kullanmak için fonksiyonlar
func getDB() *mongo.Database {
	return &urlDB
}

func getUrlCollection() *mongo.Collection {
	return &urlCollection
}

func getUserCollection() *mongo.Collection {
	return &userCollection
}
