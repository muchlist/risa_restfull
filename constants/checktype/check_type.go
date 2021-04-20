package checktype

const (
	Outstanding = "OUTSTANDING"
	Koneksi     = "KONEKSI"
	Gate        = "GATE"
	Aplikasi    = "APLIKASI"
	Ups         = "UPS"
	Server      = "SERVER"
	Improvement = "IMPROVEMENT"
	Lainnya     = "LAINNYA"
)

func GetCheckTypeAvailable() []string {
	return []string{Outstanding, Koneksi, Gate, Aplikasi, Ups, Server, Improvement, Lainnya}
}
