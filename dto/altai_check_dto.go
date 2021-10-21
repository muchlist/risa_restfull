package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// AltaiCheck full struct pengecekan altai
type AltaiCheck struct {
	ID              primitive.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt       int64                 `json:"created_at" bson:"created_at"`
	CreatedBy       string                `json:"created_by" bson:"created_by"`
	CreatedByID     string                `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt       int64                 `json:"updated_at" bson:"updated_at"`
	UpdatedBy       string                `json:"updated_by" bson:"updated_by"`
	UpdatedByID     string                `json:"updated_by_id" bson:"updated_by_id"`
	Branch          string                `json:"branch" bson:"branch"`
	TimeStarted     int64                 `json:"time_started" bson:"time_started"`
	TimeEnded       int64                 `json:"time_ended" bson:"time_ended"`
	IsFinish        bool                  `json:"is_finish" bson:"is_finish"`
	AltaiCheckItems []AltaiCheckItemEmbed `json:"altai_check_items" bson:"altai_check_items"`
	Note            string                `json:"note" bson:"note"`
}

type AltaiCheckItemEmbed struct {
	ID        string `json:"id" bson:"id"` // same as ID CCTV or general Unit cctv
	Name      string `json:"name" bson:"name"`
	Location  string `json:"location" bson:"location"`
	CheckedAt int64  `json:"checked_at" bson:"checked_at"`
	CheckedBy string `json:"checked_by" bson:"checked_by"`
	IsChecked bool   `json:"is_checked" bson:"is_checked"`
	IsOffline bool   `json:"is_offline" bson:"is_offline"`
	ImagePath string `json:"image_path" bson:"image_path"`
	DisVendor bool   `json:"dis_vendor" bson:"dis_vendor"` // if true, disable from report vendor
}

type AltaiCheckEdit struct {
	FilterIDBranch
	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string
	TimeStarted int64
	TimeEnded   int64
	IsFinish    bool
	Note        string
}

type AltaiCheckItemUpdateRequest struct {
	ParentID  string `json:"parent_id"`
	ChildID   string `json:"child_id"`
	IsChecked bool   `json:"is_checked"`
	IsOffline bool   `json:"is_offline"`
}

type AltaiCheckItemUpdate struct {
	FilterParentIDChildIDBranch
	CheckedAt int64
	CheckedBy string
	IsChecked bool
	IsOffline bool
}

type BulkAltaiCheckUpdateRequest struct {
	Items []AltaiCheckItemUpdateRequest `json:"items"`
}
