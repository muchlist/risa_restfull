package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Cctv struct penuh dari domain cctv
type Cctv struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt       int64              `json:"created_at" bson:"created_at"`
	CreatedBy       string             `json:"created_by" bson:"created_by"`
	CreatedByID     string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt       int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy       string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID     string             `json:"updated_by_id" bson:"updated_by_id"`
	Branch          string             `json:"branch" bson:"branch"`
	Disable         bool               `json:"disable" bson:"disable"`
	Name            string             `json:"name" bson:"name"`
	IP              string             `json:"ip" bson:"ip"`
	InventoryNumber string             `json:"inventory_number" bson:"inventory_number"`
	Location        string             `json:"location" bson:"location"`
	LocationLat     string             `json:"location_lat" bson:"location_lat"`
	LocationLon     string             `json:"location_lon" bson:"location_lon"`
	Date            int64              `json:"date" bson:"date"`
	Tag             []string           `json:"tag" bson:"tag"`
	Image           string             `json:"image" bson:"image"`
	Brand           string             `json:"brand" bson:"brand"`
	Type            string             `json:"type" bson:"type"`
	Note            string             `json:"note" bson:"note"`
	Extra           GenExtra           `json:"extra" bson:"extra,omitempty"`
	DisVendor       bool               `json:"dis_vendor" bson:"dis_vendor"` // if true, disable from report vendor
}

// CctvRequest user input, id tidak diinput oleh user
type CctvRequest struct {
	ID              string   `json:"-" bson:"-"`
	Name            string   `json:"name" bson:"name"`
	IP              string   `json:"ip" bson:"ip"`
	InventoryNumber string   `json:"inventory_number" bson:"inventory_number"`
	Location        string   `json:"location" bson:"location"`
	LocationLat     string   `json:"location_lat" bson:"location_lat"`
	LocationLon     string   `json:"location_lon" bson:"location_lon"`
	Date            int64    `json:"date" bson:"date"`
	Tag             []string `json:"tag" bson:"tag"`
	Image           string   `json:"image" bson:"image"`
	DisVendor       bool     `json:"dis_vendor" bson:"dis_vendor"` // if true, disable from report vendor

	Brand string `json:"brand" bson:"brand"`
	Type  string `json:"type" bson:"type"`
	Note  string `json:"note" bson:"note"`
}

type CctvEdit struct {
	ID              primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64

	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string

	Name            string
	IP              string
	InventoryNumber string
	Location        string
	LocationLat     string
	LocationLon     string
	Date            int64
	Tag             []string
	DisVendor       bool `json:"dis_vendor" bson:"dis_vendor"` // if true, disable from report vendor

	Brand string
	Type  string
	Note  string
}

// CctvEditRequest user input
type CctvEditRequest struct {
	FilterTimestamp int64 `json:"filter_timestamp"`

	Name            string   `json:"name" bson:"name"`
	IP              string   `json:"ip" bson:"ip"`
	InventoryNumber string   `json:"inventory_number" bson:"inventory_number"`
	Location        string   `json:"location" bson:"location"`
	LocationLat     string   `json:"location_lat" bson:"location_lat"`
	LocationLon     string   `json:"location_lon" bson:"location_lon"`
	Date            int64    `json:"date" bson:"date"`
	Tag             []string `json:"tag" bson:"tag"`
	DisVendor       bool     `json:"dis_vendor" bson:"dis_vendor"` // if true, disable from report vendor

	Brand string `json:"brand" bson:"brand"`
	Type  string `json:"type" bson:"type"`
	Note  string `json:"note" bson:"note"`
}

type CctvResponseMinList []CctvResponseMin

type CctvResponseMin struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Branch    string             `json:"branch" bson:"branch"`
	Disable   bool               `json:"disable" bson:"disable"`
	Name      string             `json:"name" bson:"name"`
	IP        string             `json:"ip" bson:"ip"`
	Location  string             `json:"location" bson:"location"`
	Tag       []string           `json:"tag" bson:"tag"`
	DisVendor bool               `json:"dis_vendor" bson:"dis_vendor"` // if true, disable from report vendor
}
