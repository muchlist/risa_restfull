package checktype

const (
	Outstanding = "OUTSTANDING"
	Jaringan    = "JARINGAN"
	Gate        = "GATE"
	Aplikasi    = "APLIKASI"
	Ups         = "UPS"
	Server      = "SERVER"
	Improvement = "IMPROVEMENT"
	Lainnya     = "LAINNYA"
)

func GetCheckTypeAvailable() []string {
	return []string{Outstanding, Jaringan, Gate, Aplikasi, Ups, Server, Improvement, Lainnya}
}
