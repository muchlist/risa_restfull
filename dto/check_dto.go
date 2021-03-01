package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Check struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt   int64              `json:"created_at" bson:"created_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedByID string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID string             `json:"updated_by_id" bson:"updated_by_id"`
	Branch      string             `json:"branch" bson:"branch"`

	Shift      int              `json:"shift" bson:"shift"`
	IsFinish   bool             `json:"is_finish" bson:"is_finish"`
	CheckItems []CheckItemEmbed `json:"check_items" bson:"check_items"`
	Note       string           `json:"note" bson:"note"`
}

type CheckItemEmbed struct {
	ID       string   `json:"id" bson:"id"`
	Name     string   `json:"name" bson:"name"`
	Location string   `json:"location" bson:"location"`
	Type     string   `json:"type" bson:"type"`
	Tag      []string `json:"tag" bson:"tag"`
	TagExtra []string `json:"tag_extra" bson:"tag_extra"`

	CheckedAt        int64  `json:"checked_at" bson:"checked_at"`
	IsChecked        bool   `json:"is_checked" bson:"is_checked"`
	TagSelected      string `json:"tag_selected" bson:"tag_selected"`
	TagExtraSelected string `json:"tag_extra_selected" bson:"tag_extra_selected"`
	ImagePath        string `json:"image_path" json:"image_path"`
	CheckedNote      string `json:"checked_note" bson:"checked_note"`
	HaveProblem      bool   `json:"have_problem" bson:"have_problem"`
	CompleteStatus   int    `json:"complete_status" bson:"complete_status"`
}

type CheckEdit struct {
	FilterIDBranchAuthor

	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string `json:"updated_by" bson:"updated_by"`
	UpdatedByID string `json:"updated_by_id" bson:"updated_by_id"`

	IsFinish bool   `json:"is_finish" bson:"is_finish"`
	Note     string `json:"note" bson:"note"`
}

type CheckChildUpdate struct {
	FilterParentIDChildIDAuthor

	UpdatedAt int64

	CheckedAt        int64
	IsChecked        bool
	TagSelected      string
	TagExtraSelected string
	CheckedNote      string
	HaveProblem      bool
	CompleteStatus   int
}
