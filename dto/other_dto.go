package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Other struct penuh dari domain other
type Other struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt   int64              `json:"created_at" bson:"created_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedByID string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID string             `json:"updated_by_id" bson:"updated_by_id"`
	Branch      string             `json:"branch" bson:"branch"`
	Disable     bool               `json:"disable" bson:"disable"`

	Name        string `json:"name" bson:"name"`
	Detail      string `json:"detail" bson:"detail"`
	Division    string `json:"division" bson:"division"`
	SubCategory string `json:"sub_category" bson:"sub_category"`

	IP              string   `json:"ip" bson:"ip"`
	InventoryNumber string   `json:"inventory_number" bson:"inventory_number"`
	Location        string   `json:"location" bson:"location"`
	LocationLat     string   `json:"location_lat" bson:"location_lat"`
	LocationLon     string   `json:"location_lon" bson:"location_lon"`
	Date            int64    `json:"date" bson:"date"`
	Tag             []string `json:"tag" bson:"tag"`
	Image           string   `json:"image" bson:"image"`
	Brand           string   `json:"brand" bson:"brand"`
	Type            string   `json:"type" bson:"type"`
	Note            string   `json:"note" bson:"note"`
	Extra           GenExtra `json:"extra" bson:"extra,omitempty"`
}

// OtherRequest user input, id tidak diinput oleh user
type OtherRequest struct {
	ID              string   `json:"-" bson:"-"`
	Name            string   `json:"name" bson:"name"`
	Detail          string   `json:"detail" bson:"detail"`
	Division        string   `json:"division" bson:"division"`
	SubCategory     string   `json:"sub_category" bson:"sub_category"`
	IP              string   `json:"ip" bson:"ip"`
	InventoryNumber string   `json:"inventory_number" bson:"inventory_number"`
	Location        string   `json:"location" bson:"location"`
	LocationLat     string   `json:"location_lat" bson:"location_lat"`
	LocationLon     string   `json:"location_lon" bson:"location_lon"`
	Date            int64    `json:"date" bson:"date"`
	Tag             []string `json:"tag" bson:"tag"`
	Image           string   `json:"image" bson:"image"`

	Brand string `json:"brand" bson:"brand"`
	Type  string `json:"type" bson:"type"`
	Note  string `json:"note" bson:"note"`
}

type OtherEdit struct {
	ID                primitive.ObjectID
	FilterBranch      string
	FilterSubCategory string
	FilterTimestamp   int64

	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string

	Name            string
	Detail          string
	Division        string
	IP              string
	InventoryNumber string
	Location        string
	LocationLat     string
	LocationLon     string
	Date            int64
	Tag             []string

	Brand string
	Type  string
	Note  string
}

// OtherEditRequest user input
type OtherEditRequest struct {
	FilterTimestamp   int64  `json:"filter_timestamp"`
	FilterSubCategory string `json:"filter_sub_category" bson:"filter_sub_category"`

	Name     string `json:"name" bson:"name"`
	Detail   string `json:"detail" bson:"detail"`
	Division string `json:"division" bson:"division"`

	IP              string   `json:"ip" bson:"ip"`
	InventoryNumber string   `json:"inventory_number" bson:"inventory_number"`
	Location        string   `json:"location" bson:"location"`
	LocationLat     string   `json:"location_lat" bson:"location_lat"`
	LocationLon     string   `json:"location_lon" bson:"location_lon"`
	Date            int64    `json:"date" bson:"date"`
	Tag             []string `json:"tag" bson:"tag"`

	Brand string `json:"brand" bson:"brand"`
	Type  string `json:"type" bson:"type"`
	Note  string `json:"note" bson:"note"`
}

type OtherResponseMinList []OtherResponseMin

// OtherResponse
type OtherResponseMin struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Branch      string             `json:"branch" bson:"branch"`
	Disable     bool               `json:"disable" bson:"disable"`
	Name        string             `json:"name" bson:"name"`
	Detail      string             `json:"detail" bson:"detail"`
	Division    string             `json:"division" bson:"division"`
	SubCategory string             `json:"sub_category" bson:"sub_category"`
	IP          string             `json:"ip" bson:"ip"`
	Location    string             `json:"location" bson:"location"`
	Tag         []string           `json:"tag" bson:"tag"`
}
