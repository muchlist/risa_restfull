package category

const (
	Cctv  = "CCTV"
	PC    = "PC"
	Stock = "Stock"
)

func GetCategoryAvailable() []string {
	return []string{Cctv, PC, Stock}
}
