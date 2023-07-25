package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Veri tabanı bağalantı linki ujkbqpruk3Q2cnXo
const uri = "mongodb+srv://yucelatl:ujkbqpruk3Q2cnXo@short-url.vwcugln.mongodb.net/?retryWrites=true&w=majority"

// Değişkenler
var urlDB mongo.Database
var urlCollection mongo.Collection
var userCollection mongo.Collection

func init() {

	// Veri tabanına bağlanmak için bir istemci oluştur
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
