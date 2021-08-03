package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// History struct penuh dari domain history, ID di isikan dari luar karena ada keperluan
// penamaan foto. bukan generate dari database
type History struct {
	Version        int                `json:"version" bson:"version"`
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	CreatedByID    string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt      int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy      string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID    string             `json:"updated_by_id" bson:"updated_by_id"`
	Category       string             `json:"category" bson:"category"`
	Branch         string             `json:"branch" bson:"branch"`
	ParentID       string             `json:"parent_id" bson:"parent_id"`
	ParentName     string             `json:"parent_name" bson:"parent_name"`
	Status         string             `json:"status" bson:"status"`
	Problem        string             `json:"problem" bson:"problem"`
	ProblemResolve string             `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
	DateStart      int64              `json:"date_start" bson:"date_start"`
	DateEnd        int64              `json:"date_end" bson:"date_end"`
	Tag            []string           `json:"tag" bson:"tag"`
	Image          string             `json:"image" bson:"image"`
	Updates        []HistoryUpdate    `json:"updates" bson:"updates"`
}

type HistoryUpdate struct {
	Time           int64  `json:"time" bson:"time"`
	UpdatedBy      string `json:"updated_by" bson:"updated_by"`
	UpdatedByID    string `json:"updated_by_id" bson:"updated_by_id"`
	Problem        string `json:"problem" bson:"problem"`
	ProblemResolve string `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus int    `json:"complete_status" bson:"complete_status"`
	Vendor         bool   `json:"vendor" bson:"vendor"`
}

// HistoryRequest user input
type HistoryRequest struct {
	ID             string   `json:"id,omitempty" bson:"_id,omitempty"`
	ParentID       string   `json:"parent_id" bson:"parent_id"`
	Status         string   `json:"status" bson:"status"`
	Problem        string   `json:"problem" bson:"problem"`
	ProblemResolve string   `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus int      `json:"complete_status" bson:"complete_status"`
	DateStart      int64    `json:"date_start" bson:"date_start"`
	DateEnd        int64    `json:"date_end" bson:"date_end"`
	Tag            []string `json:"tag" bson:"tag"`
}

type HistoryResponse struct {
	Version        int                `json:"version" bson:"version"`
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	UpdatedAt      int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy      string             `json:"updated_by" bson:"updated_by"`
	Category       string             `json:"category" bson:"category"`
	Branch         string             `json:"branch" bson:"branch"`
	ParentID       string             `json:"parent_id" bson:"parent_id"`
	ParentName     string             `json:"parent_name" bson:"parent_name"`
	Status         string             `json:"status" bson:"status"`
	Problem        string             `json:"problem" bson:"problem"`
	ProblemResolve string             `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
	DateStart      int64              `json:"date_start" bson:"date_start"`
	DateEnd        int64              `json:"date_end" bson:"date_end"`
	Tag            []string           `json:"tag" bson:"tag"`
	Image          string             `json:"image" bson:"image"`
	Updates        []HistoryUpdate    `json:"updates" bson:"updates"`
}

type HistoryUnwindResponseList []HistoryUnwindResponse
type HistoryUnwindResponse struct {
	Version        int                `json:"version" bson:"version"`
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	UpdatedAt      int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy      string             `json:"updated_by" bson:"updated_by"`
	Category       string             `json:"category" bson:"category"`
	Branch         string             `json:"branch" bson:"branch"`
	ParentID       string             `json:"parent_id" bson:"parent_id"`
	ParentName     string             `json:"parent_name" bson:"parent_name"`
	Status         string             `json:"status" bson:"status"`
	Problem        string             `json:"problem" bson:"problem"`
	ProblemResolve string             `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
	DateStart      int64              `json:"date_start" bson:"date_start"`
	DateEnd        int64              `json:"date_end" bson:"date_end"`
	Tag            []string           `json:"tag" bson:"tag"`
	Image          string             `json:"image" bson:"image"`
	Updates        HistoryUpdate      `json:"updates" bson:"updates"`
}

type HistoryResponseMinList []HistoryResponseMin

type HistoryResponseMin struct {
	Version        int                `json:"version" bson:"version"`
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	UpdatedAt      int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy      string             `json:"updated_by" bson:"updated_by"`
	Category       string             `json:"category" bson:"category"`
	Branch         string             `json:"branch" bson:"branch"`
	ParentID       string             `json:"parent_id" bson:"parent_id"`
	ParentName     string             `json:"parent_name" bson:"parent_name"`
	Status         string             `json:"status" bson:"status"`
	Problem        string             `json:"problem" bson:"problem"`
	ProblemResolve string             `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
	DateStart      int64              `json:"date_start" bson:"date_start"`
	DateEnd        int64              `json:"date_end" bson:"date_end"`
	Tag            []string           `json:"tag" bson:"tag"`
	Image          string             `json:"image" bson:"image"`
	Updates        []HistoryUpdate    `json:"-" bson:"updates"`
}

type HistoryCountList []HistoryCountResponse

// HistoryCountResponse tipe return dari aggregate, bson _id dirubah ke json menjadi branch
type HistoryCountResponse struct {
	Branch string `json:"branch" bson:"_id"`
	Count  int    `json:"count" bson:"count"`
}

type HistoryEdit struct {
	FilterBranch    string   `json:"filter_branch"`
	FilterTimestamp int64    `json:"filter_timestamp"`
	Status          string   `json:"status" bson:"status"`
	Problem         string   `json:"problem" bson:"problem"`
	ProblemResolve  string   `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus  int      `json:"complete_status" bson:"complete_status"`
	DateEnd         int64    `json:"date_end" bson:"date_end"`
	Tag             []string `json:"tag" bson:"tag"`
	UpdatedAt       int64    `json:"updated_at" bson:"updated_at"`
	UpdatedBy       string   `json:"updated_by" bson:"updated_by"`
	UpdatedByID     string   `json:"updated_by_id" bson:"updated_by_id"`
}

// HistoryEditRequest user input
type HistoryEditRequest struct {
	FilterTimestamp int64    `json:"filter_timestamp"`
	Status          string   `json:"status" bson:"status"`
	Problem         string   `json:"problem" bson:"problem"`
	ProblemResolve  string   `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus  int      `json:"complete_status" bson:"complete_status"`
	DateEnd         int64    `json:"date_end" bson:"date_end"`
	Tag             []string `json:"tag" bson:"tag"`
}
