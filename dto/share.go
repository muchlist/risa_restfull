package dto

// GenExtra informasi tambahan dari genUnit dengan ID yang sama dengan ID unit
type GenExtra struct {
	Cases      []Case      `json:"cases" bson:"cases"`
	CasesSize  int         `json:"cases_size" bson:"cases_size"`
	PingsState []PingState `json:"pings_state" bson:"pings_state"`
	LastPing   string      `json:"last_ping" bson:"last_ping"`
}
