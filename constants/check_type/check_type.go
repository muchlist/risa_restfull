package check_type

const (
	Ups    = "UPS"
	Server = "SERVER"
	Cctv   = "CCTV"
)

func GetCheckTypeAvailable() []string {
	return []string{Ups, Server, Cctv}
}
