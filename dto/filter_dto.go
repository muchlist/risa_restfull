package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type FilterBranchCatComplete struct {
	FilterBranch         string
	FilterCategory       string // ex "CCTV,ALTAI"
	FilterCompleteStatus []int
}

type FilterBranchCatInCompleteIn struct {
	FilterBranch         string
	FilterCategory       string // ex "CCTV,ALTAI"
	FilterCompleteStatus string // ex "0,1,2"
}

type FilterBranchLocIPNameDisable struct {
	FilterBranch   string
	FilterLocation string
	FilterIP       string
	FilterName     string
	FilterDisable  bool
}

type FilterBranchNameDisable struct {
	FilterBranch  string
	FilterName    string
	FilterDisable bool
}

type FilterBranchNameCatDisable struct {
	FilterBranch   string
	FilterName     string
	FilterCategory string
	FilterDisable  bool
}

type FilterBranchCatDisable struct {
	FilterBranch   string
	FilterCategory string
	FilterDisable  bool
}

// FilterIDBranchCreateGte
// field FilterCreateGTE digunakan untuk memberikan batas waktu, misalnya saat akan menghapus dokumen
// maka FilterCreateGTE di isi dengan tanggal sekarang kurang sekian waktu (misalnya 1 hari)
type FilterIDBranchCreateGte struct {
	FilterID        primitive.ObjectID
	FilterBranch    string
	FilterCreateGTE int64
}

// FilterIDBranchCategoryCreateGte
// field FilterCreateGTE digunakan untuk memberikan batas waktu, misalnya saat akan menghapus dokumen
// maka FilterCreateGTE di isi dengan tanggal sekarang kurang sekian waktu (misalnya 1 hari)
type FilterIDBranchCategoryCreateGte struct {
	FilterID          primitive.ObjectID
	FilterBranch      string
	FilterSubCategory string
	FilterCreateGTE   int64
}

type FilterIDBranchTimestamp struct {
	FilterID        primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64
}

type FilterIDBranch struct {
	FilterID     primitive.ObjectID
	FilterBranch string
}

type FilterIDBranchAuthor struct {
	FilterID       primitive.ObjectID
	FilterBranch   string
	FilterAuthorID string
}

type FilterParentIDChildIDAuthor struct {
	FilterParentID primitive.ObjectID
	FilterChildID  string
	FilterAuthorID string
}

type FilterParentIDChildIDBranch struct {
	FilterParentID primitive.ObjectID
	FilterChildID  string
	FilterBranch   string
}

type FilterTimeRangeLimit struct {
	FilterStart int64
	FilterEnd   int64
	Limit       int64
}

type FilterBranchCompleteTimeRangeLimit struct {
	FilterBranch         string
	FilterCompleteStatus int
	FilterStart          int64
	FilterEnd            int64
	Limit                int64
}

type FilterBranchCategory struct {
	FilterBranch   string
	FilterCategory string
}

type FilterComputer struct {
	FilterBranch         string
	FilterLocation       string
	FilterDivision       string
	FilterIP             string
	FilterName           string
	FilterDisable        bool
	FilterSeatManagement int // -1 all , 0 false, 1 true
}

type FilterOther struct {
	FilterBranch      string
	FilterSubCategory string
	FilterLocation    string
	FilterDivision    string
	FilterIP          string
	FilterName        string
	FilterDisable     bool
}
