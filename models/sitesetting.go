package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Settings struct {
	ID               primitive.ObjectID `bson:"_id"`
	SiteName         string             `bson:"site_name"`
	AdSlot           string             `bson:"ad_slot"`
	AdClient         string             `bson:"ad_client"`
	RevenuePerClick  float64            `bson:"revenue_per_click"`
	WithdrawnBalance float64            `bson:"withdrawn_balance"`
}

func (s Settings) NewSettings() (*Settings, error) {
	s.ID = primitive.NewObjectID()
	db, ctx := getSiteSettingsCollection()
	result, err := db.InsertOne(ctx, &s)
	if err != nil {
		return &Settings{}, err
	}
	s.ID = result.InsertedID.(primitive.ObjectID)
	return &s, nil
}

// find by site name
func (s Settings) FindBySiteName(siteName string) (*Settings, error) {
	db, ctx := getSiteSettingsCollection()
	filer := bson.M{"site_name": siteName}
	err := db.FindOne(ctx, filer).Decode(&s)
	if err != nil {
		return &Settings{}, err
	}
	return &s, nil
}