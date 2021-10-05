package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type PendingReport struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	CreatedByID    string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt      int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy      string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID    string             `json:"updated_by_id" bson:"updated_by_id"`
	Branch         string             `json:"branch" bson:"branch"`
	Number         string             `json:"number" bson:"number"`
	Title          string             `json:"title" bson:"title"`
	Descriptions   []PRDescription    `json:"descriptions" bson:"descriptions"`
	Date           int64              `json:"date" bson:"date"`
	Participants   []Participant      `json:"participants" bson:"participants"`
	Approvers      []Participant      `json:"approvers" bson:"approvers"`
	Equipments     []PREquipment      `json:"equipments" bson:"equipments"`
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
	Location       string             `json:"location" bson:"location"`
	Images         []string           `json:"images" bson:"images"`
}

type PendingReportRequest struct {
	Branch       string          `json:"branch" bson:"branch"`
	Number       string          `json:"number" bson:"number"`
	Title        string          `json:"title" bson:"title"`
	Descriptions []PRDescription `json:"descriptions" bson:"descriptions"`
	Date         int64           `json:"date" bson:"date"`
	Equipments   []PREquipment   `json:"equipments" bson:"equipments"`
	Location     string          `json:"location" bson:"location"`
}

type PRDescription struct {
	Description     string `json:"description" bson:"description"`
	DescriptionType string `json:"description_type" bson:"description_type"`
	Position        int    `json:"position" bson:"position"`
}

type PREquipment struct {
	ID            string `json:"id" bson:"id"`
	EquipmentName string `json:"equipment_name" bson:"equipment_name"`
	AttachTo      string `json:"attach_to" bson:"attach_to"`
	Description   string `json:"description" bson:"description"`
	Qty           int    `json:"qty" bson:"qty"`
}

type Participant struct {
	ID       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Position string `json:"position" bson:"position"`
	Division string `json:"division" bson:"division"`
	UserID   string `json:"user_id" bson:"user_id"`
	Sign     string `json:"sign" bson:"sign"`
	SignAt   int64  `json:"sign_at" bson:"sign_at"`
}
