package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// ConfigCheck full struct pengecekan config
type ConfigCheck struct {
	ID               primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt        int64                  `json:"created_at" bson:"created_at"`
	CreatedBy        string                 `json:"created_by" bson:"created_by"`
	CreatedByID      string                 `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt        int64                  `json:"updated_at" bson:"updated_at"`
	UpdatedBy        string                 `json:"updated_by" bson:"updated_by"`
	UpdatedByID      string                 `json:"updated_by_id" bson:"updated_by_id"`
	Branch           string                 `json:"branch" bson:"branch"`
	TimeStarted      int64                  `json:"time_started" bson:"time_started"`
	TimeEnded        int64                  `json:"time_ended" bson:"time_ended"`
	IsFinish         bool                   `json:"is_finish" bson:"is_finish"`
	ConfigCheckItems []ConfigCheckItemEmbed `json:"config_check_items" bson:"config_check_items"`
	Note             string                 `json:"note" bson:"note"`
}

type ConfigCheckItemEmbed struct {
	ID        string `json:"id" bson:"id"` // same as ID CCTV or general Unit cctv
	Name      string `json:"name" bson:"name"`
	CheckedAt int64  `json:"checked_at" bson:"checked_at"`
	CheckedBy string `json:"checked_by" bson:"checked_by"`
	IsUpdated bool   `json:"is_updated" bson:"is_updated"`
}

type ConfigCheckEdit struct {
	FilterIDBranch
	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string
	TimeStarted int64
	TimeEnded   int64
	IsFinish    bool
	Note        string
}

type ConfigCheckItemUpdateRequest struct {
	ParentID  string `json:"parent_id"`
	ChildID   string `json:"child_id"`
	IsUpdated bool   `json:"is_updated"`
}

type ConfigCheckItemUpdate struct {
	FilterParentIDChildIDBranch
	CheckedAt int64
	CheckedBy string
	IsUpdated bool
}

type ConfigCheckUpdateManyRequest struct {
	ParentID       string   `json:"parent_id"`
	ChildUpdate    []string `json:"child_update"`
	ChildNotUpdate []string `json:"child_not_update"`
}

type ConfigCheckUpdateMany struct {
	ParentID       string
	ChildIDsUpdate []primitive.ObjectID
	UpdatedValue   bool
	Branch         string
	Updater        string
}
