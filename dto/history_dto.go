package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// History struct penuh dari domain history, ID di isikan dari luar karena ada keperluan
// penamaan foto. bukan generate dari database
type History struct {
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	CreatedBy   string `json:"created_by" bson:"created_by"`
	CreatedByID string `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string `json:"updated_by" bson:"updated_by"`
	UpdatedByID string `json:"updated_by_id" bson:"updated_by_id"`
	Category    string `json:"category" bson:"category"`
	Branch      string `json:"branch" bson:"branch"`

	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
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
}

type HistoryResponse struct {
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	CreatedBy string `json:"created_by" bson:"created_by"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
	UpdatedBy string `json:"updated_by" bson:"updated_by"`
	Category  string `json:"category" bson:"category"`
	Branch    string `json:"branch" bson:"branch"`

	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
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
}

type HistoryResponseMinList []HistoryResponseMin

type HistoryResponseMin struct {
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	CreatedBy string `json:"created_by" bson:"created_by"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
	UpdatedBy string `json:"updated_by" bson:"updated_by"`
	Category  string `json:"category" bson:"category"`
	Branch    string `json:"branch" bson:"branch"`

	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
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
}

type HistoryCountList []HistoryCountResponse

// HistoryCountResponse tipe return dari aggregate, bson _id dirubah ke json menjadi branch
type HistoryCountResponse struct {
	Branch string `json:"branch" bson:"_id"`
	Count  int    `json:"count" bson:"count"`
}

type HistoryEdit struct {
	FilterBranch    string `json:"filter_branch"`
	FilterTimestamp int64  `json:"filter_timestamp"`

	Status         string   `json:"status" bson:"status"`
	Problem        string   `json:"problem" bson:"problem"`
	ProblemResolve string   `json:"problem_resolve" bson:"problem_resolve"`
	CompleteStatus int      `json:"complete_status" bson:"complete_status"`
	DateEnd        int64    `json:"date_end" bson:"date_end"`
	Tag            []string `json:"tag" bson:"tag"`
	UpdatedAt      int64    `json:"updated_at" bson:"updated_at"`
	UpdatedBy      string   `json:"updated_by" bson:"updated_by"`
	UpdatedByID    string   `json:"updated_by_id" bson:"updated_by_id"`
}

type FilterBranchCatComplete struct {
	Branch         string
	Category       string
	CompleteStatus int
}

type FilterIDBranchTime struct {
	ID     primitive.ObjectID
	Branch string
	Time   int64
}

type FilterTimeRangeLimit struct {
	Start int64
	End   int64
	Limit int64
}
