package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Seen struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	UrlID     primitive.ObjectID `bson:"url_id"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type ChartData struct {
	Balance map[string]float64 `json:"balance_chart"`
	Views   map[string]int     `json:"views_chart"`
}

func (seen Seen) NewSeen(userID, urlID primitive.ObjectID) error {
	date := time.Now()
	seen = Seen{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		UrlID:     urlID,
		CreatedAt: primitive.NewDateTimeFromTime(date),
	}
	db, ctx := getSeenCollection()
	response, err := db.InsertOne(ctx, seen)
	if err != nil {
		return err
	}
	fmt.Println("response: ", userID)
	seen.ID = response.InsertedID.(primitive.ObjectID)
	return nil
}

// Seen objesini kullanıcı idsi ile alıp haftalık veya istenilen günler kadar gün gün verileri array şekilde verilecek biçimde çekme
func (seen Seen) GetSeenData(userID primitive.ObjectID, days int) (ChartData, error) {
	db, ctx := getSeenCollection()
	siteData := Settings{}
	setdata, err := siteData.FindBySiteName("short-url")
	if err != nil {
		return ChartData{}, err
	}
	viewCounts := make(map[string]int)
	balanceCounts := make(map[string]float64)
	// sadece integerlardan oluşan bir array ver
	filter := bson.M{"user_id": userID, "created_at": bson.M{"$gte": time.Now().AddDate(0, 0, -days)}}
	cursor, err := db.Find(ctx, filter)
	if err != nil {
		return ChartData{}, err
	}
	for cursor.Next(ctx) {
		var seen Seen
		err := cursor.Decode(&seen)
		if err != nil {
			return ChartData{}, err
		}
		viewCounts[seen.CreatedAt.Time().Format("2006-01-02")]++
		balanceCounts[seen.CreatedAt.Time().Format("2006-01-02")] += setdata.RevenuePerClick
	}
	if err := cursor.Err(); err != nil {
		return ChartData{}, err
	}
	data := ChartData{
		Balance: balanceCounts,
		Views:   viewCounts,
	}
	return data, nil
}

func (seen Seen) FindAllSeenLength() (int64, error) {
	db, ctx := getSeenCollection()
	cursor, err := db.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return cursor, nil
}
