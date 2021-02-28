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

type FilterBranchNameCatDisable struct {
	Branch   string
	Name     string
	Category string
	Disable  bool
}

// FilterIDBranchTime
//field CreateGTE digunakan untuk memberikan batas waktu, misalnya saat akan menghapus dokumen
//maka CreateGTE di isi dengan tanggal sekarang kurang sekian waktu (misalnya 1 hari)
type FilterIDBranchTime struct {
	ID        primitive.ObjectID
	Branch    string
	CreateGTE int64
}

type FilterIDBranch struct {
	ID     primitive.ObjectID
	Branch string
}

type FilterTimeRangeLimit struct {
	Start int64
	End   int64
	Limit int64
}

type FilterBranchCategory struct {
	Branch   string
	Category string
}
