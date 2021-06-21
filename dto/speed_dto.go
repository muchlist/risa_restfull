package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpeedTest struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Time      int64              `json:"time" bson:"time"`
	LatencyMs int64              `json:"latency_ms" bson:"latency_ms"`
	Upload    float64            `json:"upload" bson:"upload"`
	Download  float64            `json:"download" bson:"download"`
}

type SpeedTestList []SpeedTest
