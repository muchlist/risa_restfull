package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// PdfFile struct penuh dari domain other
type PdfFile struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	Branch    string             `json:"branch" bson:"branch"`
	Name      string             `json:"name" bson:"name"`
	Type      string             `json:"type" bson:"type"`
	FileName  string             `json:"file_name" bson:"file_name"`
}
