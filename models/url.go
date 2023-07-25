package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Link objesi
type Url struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	OrginalUrl   string             `bson:"orginal_url,omitempty" json:"OrginalUrl"`
	ShortenedUrl string             `bson:"shortened_url,CreatedBy"`
	CreatedBy    string             `bson:"created_by,CreatedBy"`
	CreatedAt    primitive.DateTime `bson:"created_at,omitempty"`
	ValidityDays int                `bson:"validity_days,omitempty"`
	ClickCount   int                `bson:"click_count,CreatedBy"`
	RemainingDay int                `json:"RemainingDay"`
}

// Kısaltılan linki veri tabanına kaydeden fonsikyon
func (url *Url) ShortLink() (Url, error) {
	ctx := context.TODO()
	db := getDB()
	response, err := db.Collection("url").InsertOne(ctx, &url)
	if err != nil {
		return Url{}, err
	}
	now := time.Now()
	createdAt := time.Unix(time.Now().Unix(), 0)
	difference := now.Sub(createdAt).Hours() / 24
	url.RemainingDay = url.ValidityDays - int(difference)
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

	//Süresi geçmiş linklerin silinmesi
	for i, _ := range results {
		now := time.Now()
		createdAt := time.Unix(results[i].CreatedAt.Time().Unix(), 0)
		difference := now.Sub(createdAt).Hours() / 24
		results[i].RemainingDay = results[i].ValidityDays - int(difference)
		if int(difference) > results[i].ValidityDays {
			err := url.DeleteByID(results[i].ID)
			if err != nil {
				return []Url{}, err
			}
		}
	}
	return results, err
}

// Kısaltılan link başlığına sahip veriyi çeken fonksiyon
func (url Url) FindByUrl(shortenedurl string) (Url, error) {
	db := getUrlCollection()
	ctx := context.TODO()
	filter := bson.M{"shortened_url": shortenedurl}

	var result Url
	err := db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return Url{}, err
	}

	//Link çağırıldıktan sonra tıklanma sayısını güncelleme
	click := result.ClickCount + 1

	update := bson.D{{"$set", bson.D{{"click_count", click}}}}

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
func (url Url) FindByCreatedBy(username string) ([]Url, error) {
	db := getUrlCollection()
	ctx := context.TODO()
	filter := bson.M{"created_by": username}

	cursor, err := db.Find(ctx, filter)
	if err != nil {
		return []Url{}, err
	}

	var results []Url
	if err = cursor.All(ctx, &results); err != nil {
		return []Url{}, err
	}
	for i, _ := range results {
		now := time.Now()
		createdAt := time.Unix(results[i].CreatedAt.Time().Unix(), 0)
		difference := now.Sub(createdAt).Hours() / 24
		results[i].RemainingDay = results[i].ValidityDays - int(difference)
		if int(difference) > results[i].ValidityDays {
			err := url.DeleteByID(results[i].ID)
			if err != nil {
				return []Url{}, err
			}
		}
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
