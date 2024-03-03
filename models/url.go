package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Link objesi
type Url struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty"`
	Description  string             `bson:"description,omitempty"`
	OrginalUrl   string             `bson:"orginal_url,omitempty" json:"OrginalUrl"`
	ShortenedUrl string             `bson:"shortened_url,CreatedBy"`
	CreatedBy    string             `bson:"created_by,CreatedBy"`
	CreatedAt    primitive.DateTime `bson:"created_at,omitempty"`
	ClickCount   int                `bson:"click_count,CreatedBy"`
	ClickEarning float64            `bson:"click_earning,CreatedBy"`
}

// Kısaltılan linki veri tabanına kaydeden fonsikyon
func (url *Url) ShortLink() (Url, error) {
	ctx := context.TODO()
	db := getDB()
	response, err := db.Collection("url").InsertOne(ctx, &url)
	if err != nil {
		return Url{}, err
	}
	oid, _ := response.InsertedID.(primitive.ObjectID)
	url.ID = oid
	return *url, nil
}

// Bütün kısaltılan linkleri veri tabanında çeken fonksiyon en çok tıklanana göre
func (url Url) FindAllUrl() ([]Url, error) {
	ctx := context.TODO()
	db := getDB()
	//En çok tıklanana göre sıralayan ayar
	opts := options.Find().SetSort(bson.D{{"click_count", -1}})
	cursor, err := db.Collection("url").Find(ctx, bson.D{}, opts)
	if err != nil {
		return []Url{}, err
	}

	var results []Url
	if err = cursor.All(ctx, &results); err != nil {
		return []Url{}, err
	}
	return results, err
}

// Kısaltılan link başlığına sahip veriyi çeken fonksiyon
func (url Url) FindByUrl(shortenedurl string) (Url, error) {
	db := getUrlCollection()
	ctx := context.TODO()
	filter := bson.M{"shortened_url": shortenedurl}
	//Linki veri tabanından çekme
	var result Url
	err := db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return Url{}, err
	}
	siteData := Settings{}

	data, _ := siteData.FindBySiteName("short-url")
	clickearning := result.ClickEarning + data.RevenuePerClick
	//Link çağırıldıktan sonra tıklanma sayısını güncelleme
	click := result.ClickCount + 1
	update := bson.D{{"$set", bson.D{{"click_count", click}, {"click_earning", clickearning}}}}
	_, err = db.UpdateOne(ctx, filter, update)
	if err != nil {
		return Url{}, err
	}

	return result, nil
}

// Kısaltılan link idsine sahip veriyi çeken fonksiyon
func (url Url) FindByID(id primitive.ObjectID) (Url, error) {
	db := getUrlCollection()
	ctx := context.TODO()
	filter := bson.M{"_id": id}

	var result Url
	err := db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return Url{}, err
	}
	return result, nil
}

// Link sahibinin kullanıcı adına göre linkleri çekme
func (url Url) FindByCreatedBy(id primitive.ObjectID) ([]Url, error) {
	db := getUrlCollection()
	ctx := context.TODO()
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	filter := bson.M{"user_id": id}

	cursor, err := db.Find(ctx, filter, opts)
	if err != nil {
		return []Url{}, err
	}

	var results []Url
	if err = cursor.All(ctx, &results); err != nil {
		return []Url{}, err
	}
	return results, nil
}

// İd ile linki silme
func (url Url) DeleteByID(id primitive.ObjectID) error {

	db := getUrlCollection()
	ctx := context.TODO()
	filter := bson.M{"_id": id}
	_, err := db.DeleteOne(ctx, filter)
	if err != nil {
		return err

	}
	return nil
}

// update url
func (url Url) Update(id primitive.ObjectID) error {
	db := getUrlCollection()
	ctx := context.TODO()
	filter := bson.M{"_id": id}
	update := bson.D{{"$set", bson.D{{"shortened_url", url.ShortenedUrl}, {"description", url.Description}}}}
	_, err := db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
