package models

import (
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
	Balance map[string]float32 `json:"balance_chart"`
	Views   map[string]int     `json:"views_chart"`
}

func (seen Seen) NewSeen(userID, urlID primitive.ObjectID) error {
	seen = Seen{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		UrlID:     urlID,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	db, ctx := getSeenCollection()
	response, err := db.InsertOne(ctx, seen)
	if err != nil {
		return err
	}
	seen.ID = response.InsertedID.(primitive.ObjectID)
	return nil
}

// Seen objesini kullanıcı idsi ile alıp haftalık veya istenilen günler kadar gün gün verileri array şekilde verilecek biçimde çekme
func (seen Seen) GetSeenData(userID primitive.ObjectID, days int) (ChartData, error) {
	db, ctx := getSeenCollection()

	viewCounts := make(map[string]int)
	balanceCounts := make(map[string]float32)
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
		balanceCounts[seen.CreatedAt.Time().Format("2006-01-02")] += 0.2
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
