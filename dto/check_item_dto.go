package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// CheckItem struct penuh dari domain checkItem
type CheckItem struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	CreatedByID    string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt      int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy      string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID    string             `json:"updated_by_id" bson:"updated_by_id"`
	Branch         string             `json:"branch" bson:"branch"`
	Disable        bool               `json:"disable" bson:"disable"`
	Name           string             `json:"name" bson:"name"`
	Location       string             `json:"location" bson:"location"`
	LocationLat    string             `json:"location_lat" bson:"location_lat"`
	LocationLon    string             `json:"location_lon" bson:"location_lon"`
	Type           string             `json:"type" bson:"type"`
	Tag            []string           `json:"tag" bson:"tag"`
	TagExtra       []string           `json:"tag_extra" bson:"tag_extra"`
	Note           string             `json:"note" bson:"note"`
	Shifts         []int              `json:"shifts" bson:"shifts"`
	CheckedNote    string             `json:"checked_note" bson:"checked_note"`
	HaveProblem    bool               `json:"have_problem" bson:"have_problem"`
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
}

type CheckItemRequest struct {
	ID          primitive.ObjectID `json:"-" bson:"-"`
	Name        string             `json:"name" bson:"name"`
	Location    string             `json:"location" bson:"location"`
	LocationLat string             `json:"location_lat" bson:"location_lat"`
	LocationLon string             `json:"location_lon" bson:"location_lon"`
	Type        string             `json:"type" bson:"type"`
	Tag         []string           `json:"tag" bson:"tag"`
	TagExtra    []string           `json:"tag_extra" bson:"tag_extra"`
	Note        string             `json:"note" bson:"note"`
	Shifts      []int              `json:"shifts" bson:"shifts"`
}

type CheckItemEditRequest struct {
	FilterTimestamp int64    `json:"filter_timestamp"`
	Name            string   `json:"name" bson:"name"`
	Location        string   `json:"location" bson:"location"`
	LocationLat     string   `json:"location_lat" bson:"location_lat"`
	LocationLon     string   `json:"location_lon" bson:"location_lon"`
	Type            string   `json:"type" bson:"type"`
	Tag             []string `json:"tag" bson:"tag"`
	TagExtra        []string `json:"tag_extra" bson:"tag_extra"`
	Note            string   `json:"note" bson:"note"`
	Shifts          []int    `json:"shifts" bson:"shifts"`
}

type CheckItemEdit struct {
	FilterIDBranchTimestamp
	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string
	Name        string
	Location    string
	LocationLat string
	LocationLon string
	Type        string
	Tag         []string
	TagExtra    []string
	Note        string
	Shifts      []int
}

type CheckItemEditBySys struct {
	FilterID       primitive.ObjectID
	UpdatedAt      int64
	CheckedNote    string
	HaveProblem    bool
	CompleteStatus int
}

type CheckItemResponseMinList []CheckItemResponseMin
type CheckItemResponseMin struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Branch         string             `json:"branch" bson:"branch"`
	Disable        bool               `json:"disable" bson:"disable"`
	Name           string             `json:"name" bson:"name"`
	Location       string             `json:"location" bson:"location"`
	Type           string             `json:"type" bson:"type"`
	Note           string             `json:"note" bson:"note"`
	Shifts         []int              `json:"shifts" bson:"shifts"`
	Tag            []string           `json:"tag" bson:"tag"`
	TagExtra       []string           `json:"tag_extra" bson:"tag_extra"`
	CheckedNote    string             `json:"checked_note" bson:"checked_note"`
	HaveProblem    bool               `json:"have_problem" bson:"have_problem"`
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
}
