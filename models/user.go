package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Kullanıcı yapısı
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserName string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
	UrlCount int                `bson:"click_count,omitempty"`
}

// Kullanıcıyı veri tabanına kaydetme
func (usr User) CreateUser() (User, error) {
	ctx := context.TODO()
	db := getDB()
	response, err := db.Collection("user").InsertOne(ctx, &usr)
	if err != nil {
		return User{}, err
	}
	oid, _ := response.InsertedID.(primitive.ObjectID)
	usr.ID = oid
	return usr, nil
}

// Kullanıcıyı kullanıcı adı ile alma
func (user User) FindByUserName(username string) (User, error) {
	db := getUserCollection()
	ctx := context.TODO()
	filter := bson.M{"username": username}

	var result User
	err := db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return User{}, err
	}

	return result, nil
}

// Kullancıları en çok link oluşturanlara göre alma
func (user User) FindAllUsers() ([]User, error) {
	ctx := context.TODO()
	db := getDB()
	//En çok link oluşturan göre almamızı sağlayan ayar
	opts := options.Find().SetSort(bson.D{{"click_count", -1}})
	cursor, err := db.Collection("user").Find(ctx, bson.D{}, opts)
	if err != nil {
		return []User{}, err
	}

	var results []User
	if err = cursor.All(ctx, &results); err != nil {
		return []User{}, err
	}

	for i, _ := range results {
		results[i].Password = ""
	}
	return results, err
}

// Yeni link oluşturduktan sonra libk oluşturma sayısı 1 artıram fonkisyon
func (user User) NewLinkCount(username string) error {
	db := getUserCollection()
	ctx := context.TODO()
	filter := bson.M{"username": username}
	var result User
	err := db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return err
	}
	count := result.UrlCount + 1

	update := bson.D{{"$set", bson.D{{"click_count", count}}}}

	_, err = db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}
	return nil
}
