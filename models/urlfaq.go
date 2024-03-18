package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UrlFaq struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Question string             `bson:"question"`
	Answer   string             `bson:"answer"`
}

func (urlFaq UrlFaq) NewUrlFaq() (UrlFaq, error) {
	db, ctx := getUrlFaqCollection()
	result, err := db.InsertOne(ctx, urlFaq)
	if err != nil {
		return UrlFaq{}, err
	}
	urlFaq.ID = result.InsertedID.(primitive.ObjectID)
	return urlFaq, nil
}

func (urlFaq UrlFaq) FindAllUrlFaqs() ([]UrlFaq, error) {
	db, ctx := getUrlFaqCollection()
	cursor, err := db.Find(ctx, bson.D{})
	if err != nil {
		return []UrlFaq{}, err
	}
	var urlFaqs []UrlFaq
	if err = cursor.All(ctx, &urlFaqs); err != nil {
		return []UrlFaq{}, err
	}
	return urlFaqs, nil
}

// find UrlFaq by id
func (urlFaq UrlFaq) FindUrlFaqByID() (UrlFaq, error) {
	db, ctx := getUrlFaqCollection()
	err := db.FindOne(ctx, bson.D{{"_id", urlFaq.ID}}).Decode(&urlFaq)
	if err != nil {
		return UrlFaq{}, err
	}
	return urlFaq, nil
}

// Delete UrlFaq byID
func (urlFaq UrlFaq) DeleteUrlFaqByID() error {
	db, ctx := getUrlFaqCollection()
	_, err := db.DeleteOne(ctx, bson.D{{"_id", urlFaq.ID}})
	if err != nil {
		return err
	}
	return nil
}

func (urlFaq UrlFaq) UpdateUrlFaqByID() error {
	db, ctx := getUrlFaqCollection()
	_, err := db.UpdateOne(ctx, bson.D{{"_id", urlFaq.ID}}, bson.D{{"$set", bson.D{{"question", urlFaq.Question}, {"answer", urlFaq.Answer}}}})
	if err != nil {
		return err
	}
	return nil
}
