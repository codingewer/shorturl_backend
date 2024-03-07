package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaparaNo struct {
	ID       primitive.ObjectID `bson:"_id"`
	PaparaNo string             `bson:"papara_no"`
	UserId   primitive.ObjectID `bson:"user_id"`
}

// find all Papara info
func (paparaNo PaparaNo) FindAllPaparaNo() ([]PaparaNo, error) {
	db, ctx := getPaparaNoCollection()
	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return []PaparaNo{}, err
	}
	var results []PaparaNo
	if err = cursor.All(ctx, &results); err != nil {
		return []PaparaNo{}, err
	}
	return results, nil
}

func (paparaNo PaparaNo) CreatePaparaNo(userId primitive.ObjectID) (PaparaNo, error) {
	paparaNo.ID = primitive.NewObjectID()
	paparaNo.UserId = userId
	db, ctx := getPaparaNoCollection()
	response, err := db.InsertOne(ctx, &paparaNo)
	if err != nil {
		return PaparaNo{}, err
	}
	paparaNo.ID = response.InsertedID.(primitive.ObjectID)
	return paparaNo, nil
}

// update Papara info by User id
func (paparaNo PaparaNo) UpdatePaparaNo(userId primitive.ObjectID) (PaparaNo, error) {
	db, ctx := getPaparaNoCollection()
	_, err := db.UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{
		"papara_no": paparaNo.PaparaNo,
	}})
	if err != nil {
		return PaparaNo{}, err
	}
	return paparaNo, nil
}

// find Papara info by User id
func (paparaNo PaparaNo) FindPaparaNoByUserId(userId primitive.ObjectID) (PaparaNo, error) {
	db, ctx := getPaparaNoCollection()
	filter := bson.M{"user_id": userId}
	err := db.FindOne(ctx, filter).Decode(&paparaNo)
	if err != nil {
		return PaparaNo{}, err
	}
	return paparaNo, nil
}

// find Papara info by id
func (paparaNo PaparaNo) FindPaparaNoById(id primitive.ObjectID) (PaparaNo, error) {
	db, ctx := getPaparaNoCollection()
	err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&paparaNo)
	if err != nil {
		return PaparaNo{}, err
	}
	return paparaNo, nil
}

// delete Papara info by id
func (PaparaNo PaparaNo) DeletePaparaNoById(id primitive.ObjectID) error {
	db, ctx := getPaparaNoCollection()
	_, err := db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
