package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Settings struct {
	ID               primitive.ObjectID `bson:"_id"`
	SiteName         string             `bson:"site_name"`
	AboutUs          string             `bson:"about_us"`
	ContactNumber    string             `bson:"contact_number"`
	ContactEmail     string             `bson:"contact_email"`
	Address          string             `bson:"address"`
	PrivacyPolicy    string             `bson:"privacy_policy"`
	TermsConditions  string             `bson:"terms_conditions"`
	AdSlot           string             `bson:"ad_slot"`
	AdClient         string             `bson:"ad_client"`
	ReChapchaCode    string             `bson:"rechapcha_code"`
	SmtpMail         string             `bson:"smtp_mail"`
	SmtpPassword     string             `bson:"smtp_password"`
	RevenuePerClick  float64            `bson:"revenue_per_click"`
	WithdrawnBalance float64            `bson:"withdrawn_balance"`
	AllUsersLenght   int64              `bson:"all_users_length"`
	AllClicksLenght  int64              `bson:"all_clicks_length"`
	AllLinksLenght   int64              `bson:"all_links_length"`
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

// update site settings by id
func (s Settings) UpdateSettings(siteName string) (*Settings, error) {
	db, ctx := getSiteSettingsCollection()
	filer := bson.M{"site_name": siteName}
	sitedataFromDb := Settings{}
	err := db.FindOne(ctx, filer).Decode(&sitedataFromDb)
	if err != nil {
		return &Settings{}, err
	}
	if s.SmtpPassword == "" {
		s.SmtpPassword = sitedataFromDb.SmtpPassword
	}
	if s.SmtpMail == "" {
		s.SmtpMail = sitedataFromDb.SmtpMail
	}
	update := bson.M{
		"about_us":          s.AboutUs,
		"contact_number":    s.ContactNumber,
		"contact_email":     s.ContactEmail,
		"address":           s.Address,
		"ad_slot":           s.AdSlot,
		"ad_client":         s.AdClient,
		"revenue_per_click": s.RevenuePerClick,
		"withdrawn_balance": s.WithdrawnBalance,
		"privacy_policy":    s.PrivacyPolicy,
		"terms_conditions":  s.TermsConditions,
		"rechapcha_code":    s.ReChapchaCode,
		"smtp_mail":         s.SmtpMail,
		"smtp_password":     s.SmtpPassword,
	}
	_, err = db.UpdateOne(ctx, filer, bson.M{"$set": update})
	if err != nil {
		return &Settings{}, err
	}
	return &s, nil
}
