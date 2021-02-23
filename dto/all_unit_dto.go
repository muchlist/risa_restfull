package dto

// PingState hasil ping dari perangkat (di inputkan oleh sistem terpisah, pingers)
type PingState struct {
	Code   int    `json:"code" bson:"code"`
	Time   int64  `json:"time" bson:"time"`
	Status string `json:"status" bson:"status"`
}

//AllUnit struct lengkap dari document all_unit di Mongodb
type AllUnit struct {
	ID         string      `json:"id,omitempty" bson:"_id,omitempty"`
	Category   string      `json:"category" bson:"category"`
	Name       string      `json:"name" bson:"name"`
	IP         string      `json:"ip" bson:"ip"`
	Branch     string      `json:"branch" bson:"branch"`
	Cases      []string    `json:"cases" bson:"cases"`
	CasesSize  int         `json:"cases_size" bson:"cases_size"`
	PingsState []PingState `json:"pings_state,omitempty" bson:"pings_state,omitempty"`
	LastPing   string      `json:"last_ping,omitempty" bson:"last_ping,omitempty"`
}

type AllUnitRequest struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Category string `json:"category" bson:"category"`
	Name     string `json:"name" bson:"name"`
	IP       string `json:"ip" bson:"ip"`
	Branch   string `json:"branch" bson:"branch"`
}

type AllUnitEditRequest struct {
	Category string `json:"category" bson:"category"`
	Name     string `json:"name" bson:"name"`
	IP       string `json:"ip" bson:"ip"`
	Branch   string `json:"branch" bson:"branch"`
}

//AllUnitResponseList tipe slice dari AllUnitResponse
type AllUnitResponseList []AllUnitResponse

type AllUnitResponse struct {
	ID         string      `json:"id,omitempty" bson:"_id,omitempty"`
	Category   string      `json:"category" bson:"category"`
	Name       string      `json:"name" bson:"name"`
	IP         string      `json:"ip" bson:"ip"`
	Branch     string      `json:"branch" bson:"branch"`
	Cases      []string    `json:"cases" bson:"cases"`
	CasesSize  int         `json:"cases_size" bson:"cases_size"`
	PingsState []PingState `json:"pings_state,omitempty" bson:"pings_state,omitempty"`
	LastPing   string      `json:"last_ping,omitempty" bson:"last_ping,omitempty"`
}
