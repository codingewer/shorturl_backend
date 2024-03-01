package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HelpRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"ID"`
	Title     string             `bson:"title" json:"Title"`
	Content   string             `bson:"content" json:"Content"`
	Answer    string             `bson:"answer" json:"Answer"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userID"`
	User      ResponseUser       `json:"user"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	Status    bool               `bson:"status" json:"status"`
}

func (help HelpRequest) NewHelpRequest(userId primitive.ObjectID) (HelpRequest, error) {
	db, ctx := getHelpCollection()
	help.ID = primitive.NewObjectID()
	help.UserID = userId
	help.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	help.Status = false

	user := User{}
	response, err := db.InsertOne(ctx, help)
	if err != nil {
		return help, err
	}
	responseUser, err := user.FindResposeUserByID(userId)
	if err != nil {
		return help, err
	}
	help.ID = response.InsertedID.(primitive.ObjectID)
	help.User = responseUser
	return help, nil
}

// Find by status
func (help HelpRequest) FindByStatus(status bool) ([]HelpRequest, error) {
	db, ctx := getHelpCollection()
	var helpRequests []HelpRequest
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	cursor, err := db.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return helpRequests, err
	}
	err = cursor.All(ctx, &helpRequests)
	if err != nil {
		return helpRequests, err
	}
	return helpRequests, nil
}

//Change status

func (help HelpRequest) ChangeStatus() error {
	db, ctx := getHelpCollection()
	filter := bson.M{"_id": help.ID}
	update := bson.M{"$set": bson.M{"status": help.Status, "answer": help.Answer}}
	_, err := db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

//fin by user id

func (help HelpRequest) FindByUserId(userId primitive.ObjectID) ([]HelpRequest, error) {
	db, ctx := getHelpCollection()
	var helpRequests []HelpRequest
	//find new to old
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	filter := bson.M{"user_id": userId}
	cursor, err := db.Find(ctx, filter, opts)
	if err != nil {
		return helpRequests, err
	}
	err = cursor.All(ctx, &helpRequests)
	if err != nil {
		return helpRequests, err
	}
	return helpRequests, nil
}
