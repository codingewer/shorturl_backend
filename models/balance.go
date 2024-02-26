package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BalanceRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	UserId    primitive.ObjectID `json:"userId"  bson:"user_id"`
	User      ResponseUser       `json:"user"  bson:"-"`
	Amount    float64            `json:"amount"  bson:"amount"`
	Status    bool               `json:"status"  bson:"status"`
	CreatedAt primitive.DateTime `json:"createdAt"  bson:"created_at"`
}

type BalanceInfo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	UserId    primitive.ObjectID `json:"userId"  bson:"user_id"`
	Iban      string             `json:"iban"  bson:"iban"`
	IbanOwner string             `json:"ibanOwner"  bson:"iban_owner"`
}

func (balanceReq BalanceRequest) CreateNewRequest() (BalanceRequest, error) {
	balanceReq.Status = false
	balanceReq.ID = primitive.NewObjectID()
	balanceReq.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	user := User{}
	err := user.UpdateBalance(balanceReq.UserId, balanceReq.Amount)
	if err != nil {
		fmt.Println("23")
		return BalanceRequest{}, err
	}
	db, ctx := getBalanceCollection()
	response, err := db.InsertOne(ctx, &balanceReq)
	if err != nil {
		return BalanceRequest{}, err
	}
	balanceReq.ID = response.InsertedID.(primitive.ObjectID)
	return balanceReq, nil
}

func (balancereq BalanceRequest) UpdateRequestStatus(status bool) error {
	db, ctx := getBalanceCollection()
	_, err := db.UpdateOne(ctx, bson.M{"_id": balancereq.ID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}
	return nil
}

func (balancereq BalanceRequest) FindRequestsByUserID(id primitive.ObjectID) ([]BalanceRequest, error) {
	db, ctx := getBalanceCollection()
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := db.Find(ctx, bson.M{"user_id": id}, opts)
	if err != nil {
		return []BalanceRequest{}, err
	}
	var results []BalanceRequest
	if err = cursor.All(ctx, &results); err != nil {
		return []BalanceRequest{}, err
	}
	return results, nil
}

func (balancereq BalanceRequest) FindRequestsByStatus(status bool) ([]BalanceRequest, error) {
	db, ctx := getBalanceCollection()
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := db.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return []BalanceRequest{}, err
	}
	var results []BalanceRequest
	if err = cursor.All(ctx, &results); err != nil {
		return []BalanceRequest{}, err
	}
	for i, _ := range results {
		user := User{}
		userr, err := user.FindResposeUserByID(results[i].UserId)
		if err != nil {
			return []BalanceRequest{}, err
		}
		results[i].User = userr
	}
	return results, nil
}

// find all balance info
func (balanceInfo BalanceInfo) FindAllBalanceInfo() ([]BalanceInfo, error) {
	db, ctx := getBalanceInfoCollection()
	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return []BalanceInfo{}, err
	}
	var results []BalanceInfo
	if err = cursor.All(ctx, &results); err != nil {
		return []BalanceInfo{}, err
	}
	return results, nil
}

func (balanceInfo BalanceInfo) CreateBalanceInfo(userId primitive.ObjectID) (BalanceInfo, error) {
	balanceInfo.ID = primitive.NewObjectID()
	balanceInfo.UserId = userId
	db, ctx := getBalanceInfoCollection()
	response, err := db.InsertOne(ctx, &balanceInfo)
	if err != nil {
		return BalanceInfo{}, err
	}
	balanceInfo.ID = response.InsertedID.(primitive.ObjectID)
	return balanceInfo, nil
}

// update balance info by User id
func (balanceInfo BalanceInfo) UpdateBalanceInfo(userId primitive.ObjectID) (BalanceInfo, error) {
	db, ctx := getBalanceInfoCollection()
	_, err := db.UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"iban": balanceInfo.Iban, "iban_owner": balanceInfo.IbanOwner}})
	if err != nil {
		return BalanceInfo{}, err
	}
	return balanceInfo, nil
}

// find balance info by User id
func (balanceInfo BalanceInfo) FindBalanceInfoByUserId(userId primitive.ObjectID) (BalanceInfo, error) {
	db, ctx := getBalanceInfoCollection()
	err := db.FindOne(ctx, bson.M{"user_id": userId}).Decode(&balanceInfo)
	if err != nil {
		return BalanceInfo{}, err
	}
	return balanceInfo, nil
}

// find balance info by id
func (balanceInfo BalanceInfo) FindBalanceInfoById(id primitive.ObjectID) (BalanceInfo, error) {
	db, ctx := getBalanceInfoCollection()
	err := db.FindOne(ctx, bson.M{"_id": id}).Decode(&balanceInfo)
	if err != nil {
		return BalanceInfo{}, err
	}
	return balanceInfo, nil
}
