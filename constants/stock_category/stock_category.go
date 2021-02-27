package stock_category

const (
	Computer     = "KOMPUTER"
	ComputerPart = "PART KOMPUTER"
	Ups          = "UPS"
	Printer      = "PRINTER"
	Cctv         = "CCTV"
	Jaringan     = "JARINGAN"
	Server       = "SERVER"
	Handheld     = "HANDHELD"
	Display      = "DISPLAY"
	Presensi     = "PRESENSI"
	Lainnya      = "LAINNYA"
)

func GetStockCategoryAvailable() []string {
	return []string{
		Computer,
		ComputerPart,
		Ups,
		Printer,
		Cctv,
		Jaringan,
		Server,
		Handheld,
		Display,
		Presensi,
		Lainnya,
	}
}
