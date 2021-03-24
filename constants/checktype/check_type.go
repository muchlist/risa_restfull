package checktype

const (
	Ups    = "UPS"
	Server = "SERVER"
	Cctv   = "CCTV"
)

func GetCheckTypeAvailable() []string {
	return []string{Ups, Server, Cctv}
}
