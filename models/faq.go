package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Faq struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Question string             `bson:"question"`
	Answer   string             `bson:"answer"`
}

func (faq Faq) NewFaq() (Faq, error) {
	db, ctx := getBalanceCollection()
	result, err := db.InsertOne(ctx, faq)
	if err != nil {
		return Faq{}, err
	}
	faq.ID = result.InsertedID.(primitive.ObjectID)
	return faq, nil
}

func (faq Faq) FindAllFaqs() ([]Faq, error) {
	db, ctx := getBalanceCollection()

	cursor, err := db.Find(ctx, bson.D{})
	if err != nil {
		return []Faq{}, err
	}
	var faqs []Faq
	if err = cursor.All(ctx, &faqs); err != nil {
		return []Faq{}, err
	}
	return faqs, nil
}

// find faq by id
func (faq Faq) FindFaqByID() (Faq, error) {
	db, ctx := getBalanceCollection()
	err := db.FindOne(ctx, bson.D{{"_id", faq.ID}}).Decode(&faq)
	if err != nil {
		return Faq{}, err
	}
	return faq, nil
}

// Delete Faq byID
func (faq Faq) DeleteFaqByID() error {
	db, ctx := getBalanceCollection()
	_, err := db.DeleteOne(ctx, bson.D{{"_id", faq.ID}})
	if err != nil {
		return err
	}
	return nil
}

func (faq Faq) UpdateFaqByID() error {
	db, ctx := getBalanceCollection()
	_, err := db.UpdateOne(ctx, bson.D{{"_id", faq.ID}}, bson.D{{"$set", bson.D{{"question", faq.Question}, {"answer", faq.Answer}}}})
	if err != nil {
		return err
	}
	return nil
}
