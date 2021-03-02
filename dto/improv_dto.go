package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Improve struct penuh dari domain improve
type Improve struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt   int64              `json:"created_at" bson:"created_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedByID string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID string             `json:"updated_by_id" bson:"updated_by_id"`
	Branch      string             `json:"branch" bson:"branch"`

	Title          string          `json:"title" bson:"title"`
	Description    string          `json:"description" bson:"description"`
	Goal           int             `json:"goal" bson:"goal"`
	GoalsAchieved  int             `json:"goals_achieved" bson:"goals_achieved"`
	IsActive       bool            `json:"is_active" bson:"is_active"`
	CompleteStatus int             `json:"complete_status" bson:"complete_status"`
	ImproveChanges []ImproveChange `json:"improve_changes" bson:"improve_changes"`
}

// StockChange disertakan di model penuh Stock
type ImproveChange struct {
	DummyID   int64  `json:"dummy_id" bson:"dummy_id"`
	Author    string `json:"author" bson:"author"`
	Increment int    `json:"increment" bson:"increment"`
	Note      string `json:"note" bson:"note"`
	Time      int64  `json:"time" bson:"time"`
}

// ImproveChangeRequest input user
type ImproveChangeRequest struct {
	DummyID   int64  `json:"-" bson:"dummy_id"`
	Author    string `json:"author" bson:"author"`
	Increment int    `json:"increment" bson:"increment"`
	Note      string `json:"note" bson:"note"`
	Time      int64  `json:"time" bson:"time"`
}

// ImproveRequest input user
type ImproveRequest struct {
	Title          string `json:"title" bson:"title"`
	Description    string `json:"description" bson:"description"`
	Goal           int    `json:"goal" bson:"goal"`
	CompleteStatus int    `json:"complete_status"`
}

// ImproveEditRequest input user
type ImproveEditRequest struct {
	FilterTimeStamp int64  `json:"filter_time_stamp"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Goal            int    `json:"goal"`
	CompleteStatus  int    `json:"complete_status"`
}

type ImproveEdit struct {
	FilterIDBranchTimestamp
	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string

	Title          string
	Description    string
	Goal           int
	CompleteStatus int
}

type ImproveResponseMinList []ImproveResponseMin
type ImproveResponseMin struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at"`
	Branch    string             `json:"branch" bson:"branch"`

	Title          string `json:"title" bson:"title"`
	Description    string `json:"description" bson:"description"`
	Goal           int    `json:"goal" bson:"goal"`
	GoalsAchieved  int    `json:"goals_achieved" bson:"goals_achieved"`
	IsActive       bool   `json:"is_active" bson:"is_active"`
	CompleteStatus int    `json:"complete_status" bson:"complete_status"`
}
