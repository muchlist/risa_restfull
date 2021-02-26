package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type FilterBranchCatComplete struct {
	Branch         string
	Category       string
	CompleteStatus int
}

type FilterBranchLocIPNameDisable struct {
	Branch   string
	Location string
	IP       string
	Name     string
	Disable  bool
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
