package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ServerFile struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy string             `json:"updated_by" bson:"updated_by"`
	Branch    string             `json:"branch" bson:"branch"`
	Title     string             `json:"title" bson:"title"`
	Note      string             `json:"note" bson:"note"`
	Diff      string             `json:"diff" bson:"diff"`
	Image     string             `json:"image" bson:"image"`
}
