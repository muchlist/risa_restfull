package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AltaiPhyCheck struct {
	ID                 primitive.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	QuarterlyMode      bool                     `json:"quarterly_mode" bson:"quarterly_mode"`
	Name               string                   `json:"name" bson:"name"`
	CreatedAt          int64                    `json:"created_at" bson:"created_at"`
	CreatedBy          string                   `json:"created_by" bson:"created_by"`
	CreatedByID        string                   `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt          int64                    `json:"updated_at" bson:"updated_at"`
	UpdatedBy          string                   `json:"updated_by" bson:"updated_by"`
	UpdatedByID        string                   `json:"updated_by_id" bson:"updated_by_id"`
	Branch             string                   `json:"branch" bson:"branch"`
	TimeStarted        int64                    `json:"time_started" bson:"time_started"`
	TimeEnded          int64                    `json:"time_ended" bson:"time_ended"`
	IsFinish           bool                     `json:"is_finish" bson:"is_finish"`
	AltaiPhyCheckItems []AltaiPhyCheckItemEmbed `json:"altai_phy_check_items" bson:"altai_phy_check_items"`
	Note               string                   `json:"note" bson:"note"`
}

type AltaiPhyCheckItemEmbed struct {
	ID           string `json:"id" bson:"id"` // same as ID CCTV or general Unit cctv
	Name         string `json:"name" bson:"name"`
	Location     string `json:"location" bson:"location"`
	CheckedAt    int64  `json:"checked_at" bson:"checked_at"`
	CheckedBy    string `json:"checked_by" bson:"checked_by"`
	IsChecked    bool   `json:"is_checked" bson:"is_checked"`
	IsMaintained bool   `json:"is_maintained" bson:"is_maintained"`
	IsOffline    bool   `json:"is_offline" bson:"is_offline"`
	ImagePath    string `json:"image_path" bson:"image_path"`
}

type AltaiPhyCheckEdit struct {
	FilterIDBranch
	Name        string
	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string
	TimeStarted int64
	TimeEnded   int64
	IsFinish    bool
	Note        string
}

type AltaiPhyCheckItemUpdateRequest struct {
	ParentID     string `json:"parent_id"`
	ChildID      string `json:"child_id"`
	IsChecked    bool   `json:"is_checked"`
	IsMaintained bool   `json:"is_maintained"`
	IsOffline    bool   `json:"is_offline"`
}

type AltaiPhyCheckItemUpdate struct {
	FilterParentIDChildIDBranch
	CheckedAt    int64
	CheckedBy    string
	IsChecked    bool
	IsMaintained bool
	IsOffline    bool
}

type BulkAltaiPhyCheckUpdateRequest struct {
	Items []AltaiPhyCheckItemUpdateRequest `json:"items"`
}
