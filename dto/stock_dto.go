package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Stock struct penuh dari domain cctv
type Stock struct {
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	CreatedBy   string `json:"created_by" bson:"created_by"`
	CreatedByID string `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string `json:"updated_by" bson:"updated_by"`
	UpdatedByID string `json:"updated_by_id" bson:"updated_by_id"`
	Branch      string `json:"branch" bson:"branch"`
	Disable     bool   `json:"disable" bson:"disable"`

	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	StockCategory string             `json:"stock_category" bson:"stock_category"`
	Unit          string             `json:"unit" bson:"unit"`
	Qty           int                `json:"qty" bson:"qty"`
	Location      string             `json:"location" bson:"location"`
	Threshold     int                `json:"threshold" bson:"threshold"`
	Increment     []StockChange      `json:"increment" bson:"increment"`
	Decrement     []StockChange      `json:"decrement" bson:"decrement"`
	Tag           []string           `json:"tag" bson:"tag"`
	Image         string             `json:"image" bson:"image"`
	Note          string             `json:"note" bson:"note"`
}

type StockChange struct {
	DummyID  int64  `json:"dummy_id" bson:"dummy_id"`
	Author   string `json:"author" bson:"author"`
	Qty      int    `json:"qty" bson:"qty"`
	BaNumber string `json:"ba_number" bson:"ba_number"`
	Note     string `json:"note" bson:"note"`
	Time     int64  `json:"time" bson:"time"`
}
