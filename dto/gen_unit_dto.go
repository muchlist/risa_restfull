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

// GenUnit struct lengkap dari document all_unit di Mongodb, ID harus sama dengan ID unitDetail
type GenUnit struct {
	ID         string      `json:"id,omitempty" bson:"_id,omitempty"`
	Category   string      `json:"category" bson:"category"`
	Name       string      `json:"name" bson:"name"`
	IP         string      `json:"ip" bson:"ip"`
	Branch     string      `json:"branch" bson:"branch"`
	Disable    bool        `json:"disable" bson:"disable"`
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
	Disable  bool
	Pings    bool
	LastPing string
}

// GenUnitCaseRequest,
// UnitID adalah FilterID dari gen_unit,
// CaseID sama dengan historyID.hex()
// CaseNote tidak masalah dikosongkan jika ingin menghapus Case
type GenUnitCaseRequest struct {
	UnitID       string
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

// GenUnitResponseList tipe slice dari GenUnitResponse
type GenUnitResponseList []GenUnitResponse

type GenUnitResponse struct {
	ID         string      `json:"id,omitempty" bson:"_id,omitempty"`
	Category   string      `json:"category" bson:"category"`
	Name       string      `json:"name" bson:"name"`
	IP         string      `json:"ip" bson:"ip"`
	Branch     string      `json:"branch" bson:"branch"`
	Cases      []Case      `json:"cases" bson:"cases"`
	CasesSize  int         `json:"cases_size" bson:"cases_size"`
	PingsState []PingState `json:"pings_state" bson:"pings_state"`
	LastPing   string      `json:"last_ping" bson:"last_ping"`
	Disable    bool        `json:"-" bson:"disable"`
}

type GenUnitIPList []IPAddressContainer
type IPAddressContainer struct {
	IP string `json:"ip" bson:"ip"`
}

type GenUnitPingStateRequest struct {
	Branch      string   `json:"branch"`
	Category    string   `json:"category"`
	IPAddresses []string `json:"ip_addresses"`
	PingCode    int      `json:"ping_code"`
}
