package enum

const (
	Draft = iota
	NeedSign
	CompletedSign
)

// GetDocProgress mengembalikan string dari enum progress status
func GetDocProgress(status int) string {
	switch status {
	case Draft:
		return "Draft"
	case NeedSign:
		return "Perlu-tanda-tangan"
	case CompletedSign:
		return "Selesai"
	default:
		return "Unknown"
	}
}
