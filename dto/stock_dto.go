package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Stock struct penuh dari domain cctv
type Stock struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt     int64              `json:"created_at" bson:"created_at"`
	CreatedBy     string             `json:"created_by" bson:"created_by"`
	CreatedByID   string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt     int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy     string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID   string             `json:"updated_by_id" bson:"updated_by_id"`
	Branch        string             `json:"branch" bson:"branch"`
	Disable       bool               `json:"disable" bson:"disable"`
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

// StockChange disertakan di model penuh Stock
type StockChange struct {
	DummyID  int64  `json:"dummy_id" bson:"dummy_id"`
	Author   string `json:"author" bson:"author"`
	Qty      int    `json:"qty" bson:"qty"`
	BaNumber string `json:"ba_number" bson:"ba_number"`
	Note     string `json:"note" bson:"note"`
	Time     int64  `json:"time" bson:"time"`
}

// StockChangeRequest input user
type StockChangeRequest struct {
	DummyID  int64  `json:"-" bson:"dummy_id"`
	Author   string `json:"author" bson:"author"`
	Qty      int    `json:"qty" bson:"qty"`
	BaNumber string `json:"ba_number" bson:"ba_number"`
	Note     string `json:"note" bson:"note"`
	Time     int64  `json:"time" bson:"time"`
}

type StockRequest struct {
	Name          string   `json:"name" bson:"name"`
	StockCategory string   `json:"stock_category" bson:"stock_category"`
	Unit          string   `json:"unit" bson:"unit"`
	Qty           int      `json:"qty" bson:"qty"`
	Location      string   `json:"location" bson:"location"`
	Threshold     int      `json:"threshold" bson:"threshold"`
	Tag           []string `json:"tag" bson:"tag"`
	Note          string   `json:"note" bson:"note"`
}

type StockEdit struct {
	ID              primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64
	UpdatedAt       int64
	UpdatedBy       string
	UpdatedByID     string
	Name            string
	StockCategory   string
	Unit            string
	Location        string
	Threshold       int
	Tag             []string
	Note            string
}

type StockEditRequest struct {
	FilterTimestamp int64    `json:"filter_timestamp"`
	Name            string   `json:"name"`
	StockCategory   string   `json:"stock_category"`
	Unit            string   `json:"unit"`
	Location        string   `json:"location"`
	Threshold       int      `json:"threshold"`
	Tag             []string `json:"tag"`
	Note            string   `json:"note"`
}

type StockResponseMinList []StockResponseMin
type StockResponseMin struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Branch        string             `json:"branch" bson:"branch"`
	Disable       bool               `json:"disable" bson:"disable"`
	Name          string             `json:"name" bson:"name"`
	StockCategory string             `json:"stock_category" bson:"stock_category"`
	Unit          string             `json:"unit" bson:"unit"`
	Qty           int                `json:"qty" bson:"qty"`
	Location      string             `json:"location" bson:"location"`
	Threshold     int                `json:"threshold" bson:"threshold"`
	Tag           []string           `json:"tag" bson:"tag"`
	Image         string             `json:"image" bson:"image"`
	Note          string             `json:"note" bson:"note"`
}
