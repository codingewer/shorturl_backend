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
var balanceCollection mongo.Collection
var seenCollection mongo.Collection
var helpCollection mongo.Collection

func init() {

	//.env dosyasını yükle
	err := godotenv.Load()
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
	balanceCollection = *urlDB.Collection("balance")
	seenCollection = *urlDB.Collection("seen")
	helpCollection = *urlDB.Collection("help")

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

func getBalanceCollection() (*mongo.Collection, context.Context) {
	return &balanceCollection, context.TODO()
}

func getSeenCollection() (*mongo.Collection, context.Context) {
	return &seenCollection, context.TODO()
}

func getHelpCollection() (*mongo.Collection, context.Context) {
	return &helpCollection, context.TODO()
}
