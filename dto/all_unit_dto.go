package dto

// PingState hasil ping dari perangkat (di inputkan oleh sistem terpisah, pingers)
type PingState struct {
	Code   int    `json:"code" bson:"code"`
	Time   int64  `json:"time" bson:"time"`
	Status string `json:"status" bson:"status"`
}

// Case setiap history dibuat akan menambah case array di Allunit document
type Case struct {
	CaseID   string `json:"case_id" bson:"case_id"`
	CaseNote string `json:"case_note" bson:"case_note"`
}

//GenUnit struct lengkap dari document all_unit di Mongodb
type GenUnit struct {
	ID         string      `json:"id,omitempty" bson:"_id,omitempty"`
	Category   string      `json:"category" bson:"category"`
	Name       string      `json:"name" bson:"name"`
	IP         string      `json:"ip" bson:"ip"`
	Branch     string      `json:"branch" bson:"branch"`
	Cases      []Case      `json:"cases" bson:"cases"`
	CasesSize  int         `json:"cases_size" bson:"cases_size"`
	PingsState []PingState `json:"pings_state,omitempty" bson:"pings_state,omitempty"`
	LastPing   string      `json:"last_ping,omitempty" bson:"last_ping,omitempty"`
}

type GenUnitRequest struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Category string `json:"category" bson:"category"`
	Name     string `json:"name" bson:"name"`
	IP       string `json:"ip" bson:"ip"`
	Branch   string `json:"branch" bson:"branch"`
}

type GenUnitFilter struct {
	Branch   string
	Name     string
	Category string
	IP       string
}

type GenUnitCaseRequest struct {
	ID           string
	FilterBranch string
	CaseID       string
	CaseNote     string
}

type GenUnitEditRequest struct {
	Category string `json:"category" bson:"category"`
	Name     string `json:"name" bson:"name"`
	IP       string `json:"ip" bson:"ip"`
	Branch   string `json:"branch" bson:"branch"`
}

//GenUnitResponseList tipe slice dari GenUnitResponse
type GenUnitResponseList []GenUnitResponse

type GenUnitResponse struct {
	ID         string      `json:"id,omitempty" bson:"_id,omitempty"`
	Category   string      `json:"category" bson:"category"`
	Name       string      `json:"name" bson:"name"`
	IP         string      `json:"ip" bson:"ip"`
	Branch     string      `json:"branch" bson:"branch"`
	Cases      []Case      `json:"cases" bson:"cases"`
	CasesSize  int         `json:"cases_size" bson:"cases_size"`
	PingsState []PingState `json:"pings_state,omitempty" bson:"pings_state,omitempty"`
	LastPing   string      `json:"last_ping,omitempty" bson:"last_ping,omitempty"`
}
