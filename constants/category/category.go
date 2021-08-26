package category

const (
	Cctv  = "CCTV"
	PC    = "PC"
	Stock = "STOCK"

	Application = "APPLICATION"
	Printer     = "PRINTER"
	Altai       = "ALTAI"
	Handheld    = "HANDHELD"
	RouterSW    = "ROUTER_SWITCH"
	Server      = "SERVER"
	Gate        = "GATE"
	Ups         = "UPS"
	Other       = "OTHER"
	OtherV      = "OTHER-V"
)

func GetCategoryAvailable() []string {
	return []string{Cctv, PC, Stock}
}

func GetSubCategoryAvailable() []string {
	return []string{
		Application,
		Printer,
		Altai,
		Handheld,
		RouterSW,
		Server,
		Gate,
		Ups,
		Other,
		OtherV,
	}
}
