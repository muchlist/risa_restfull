package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

// PendingReportModel struct utama Pending reports
type PendingReportModel struct {
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
	Equipments     []PREquipment      `json:"equipments" bson:"equipments"` // equipment position akan ditentukan di description dengan type equip
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
	Location       string             `json:"location" bson:"location"`
	Images         []string           `json:"images" bson:"images"`
	DocType        string             `json:"doc_type" bson:"doc_type"`
}

// NormalizeValue digunakan untuk mencegah ada nilai nil pada struct, terutama saat dimasukkan ke database mongodb yang bisa
// menyebabkan error
func (pd *PendingReportModel) NormalizeValue() {
	if pd.Descriptions == nil {
		pd.Descriptions = make([]PRDescription, 0)
	}
	if pd.Participants == nil {
		pd.Participants = make([]Participant, 0)
	}
	if pd.Approvers == nil {
		pd.Approvers = make([]Participant, 0)
	}
	if pd.Equipments == nil {
		pd.Equipments = make([]PREquipment, 0)
	}
	if pd.Images == nil {
		pd.Images = make([]string, 0)
	}

	pd.Title = strings.ToUpper(pd.Title)
	pd.Number = strings.ToUpper(pd.Number)
	pd.Branch = strings.ToUpper(pd.Branch)

}

type PendingReportEditRequest struct {
	FilterTimestamp int64 `json:"filter_timestamp"`

	Number       string          `json:"number"`
	Title        string          `json:"title"`
	Descriptions []PRDescription `json:"descriptions"`
	Date         int64           `json:"date"`
	Equipments   []PREquipment   `json:"equipments"`
	Location     string          `json:"location"`
}

type PendingReportEditModel struct {
	FilterID        primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64

	UpdatedAt    int64           `json:"updated_at" bson:"updated_at"`
	UpdatedBy    string          `json:"updated_by" bson:"updated_by"`
	UpdatedByID  string          `json:"updated_by_id" bson:"updated_by_id"`
	Number       string          `json:"number" bson:"number"`
	Title        string          `json:"title" bson:"title"`
	Descriptions []PRDescription `json:"descriptions" bson:"descriptions"`
	Date         int64           `json:"date" bson:"date"`
	Equipments   []PREquipment   `json:"equipments" bson:"equipments"`
	Location     string          `json:"location" bson:"location"`
}

func (pd *PendingReportEditModel) NormalizeValue() {
	if pd.Descriptions == nil {
		pd.Descriptions = make([]PRDescription, 0)
	}
	if pd.Equipments == nil {
		pd.Equipments = make([]PREquipment, 0)
	}

	pd.Title = strings.ToUpper(pd.Title)
	pd.Number = strings.ToUpper(pd.Number)
	pd.FilterBranch = strings.ToUpper(pd.FilterBranch)
}

type PendingReportRequest struct {
	Branch       string          `json:"branch" bson:"branch"`
	Number       string          `json:"number" bson:"number"`
	Title        string          `json:"title" bson:"title"`
	Descriptions []PRDescription `json:"descriptions" bson:"descriptions"`
	Date         int64           `json:"date" bson:"date"`
	Equipments   []PREquipment   `json:"equipments" bson:"equipments"`
	Location     string          `json:"location" bson:"location"`
	DocType      string          `json:"doc_type" bson:"doc_type"`
}

type PendingReportTempOneRequest struct {
	Branch     string        `json:"branch"`
	Number     string        `json:"number"`
	Title      string        `json:"title"`
	Actions    []string      `json:"actions"`
	Date       int64         `json:"date"`
	Equipments []PREquipment `json:"equipments"`
	Location   string        `json:"location"`
}

type PRDescription struct {
	Position        int    `json:"position" bson:"position"`                 // posisi urutan
	Description     string `json:"description" bson:"description"`           // isi dari suratnya
	DescriptionType string `json:"description_type" bson:"description_type"` // tipe tampilan, [???]
}

type PREquipment struct {
	ID            string `json:"id" bson:"id"`                         // id alat sesuai database
	EquipmentName string `json:"equipment_name" bson:"equipment_name"` // equip penamaan [] maybe stok
	AttachTo      string `json:"attach_to" bson:"attach_to"`           // dipasang di mesin mana ?
	Description   string `json:"description" bson:"description"`       // deskripsi alat
	Qty           int    `json:"qty" bson:"qty"`                       // jumlah stok yang berkurang
}

type Participant struct {
	ID       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Position string `json:"position" bson:"position"`
	Division string `json:"division" bson:"division"`
	UserID   string `json:"user_id" bson:"user_id"`
	Sign     string `json:"sign" bson:"sign"`
	SignAt   int64  `json:"sign_at" bson:"sign_at"`
	Alias    string `json:"alias" bson:"alias"`
}

// PendingReportResponse struct
type PendingReportResponse struct {
	ID             string          `json:"id"`
	CreatedAt      int64           `json:"created_at"`
	CreatedBy      string          `json:"created_by"`
	CreatedByID    string          `json:"created_by_id"`
	UpdatedAt      int64           `json:"updated_at"`
	UpdatedBy      string          `json:"updated_by"`
	UpdatedByID    string          `json:"updated_by_id"`
	Branch         string          `json:"branch"`
	Number         string          `json:"number"`
	Title          string          `json:"title"`
	Descriptions   []PRDescription `json:"descriptions"`
	Date           int64           `json:"date"`
	Participants   []Participant   `json:"participants"`
	Approvers      []Participant   `json:"approvers"`
	Equipments     []PREquipment   `json:"equipments"`
	CompleteStatus int             `json:"complete_status"`
	Location       string          `json:"location"`
	Images         []string        `json:"images"`
	DocType        string          `json:"doc_type"`
}

type PendingReportMin struct {
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
	Date           int64              `json:"date" bson:"date"`
	Participants   []Participant      `json:"participants" bson:"participants"`
	Approvers      []Participant      `json:"approvers" bson:"approvers"`
	CompleteStatus int                `json:"complete_status" bson:"complete_status"`
	Location       string             `json:"location" bson:"location"`
	Images         []string           `json:"images" bson:"images"`
	DocType        string             `json:"doc_type" bson:"doc_type"`
}
