package enum

const (
	HDataInfo = iota - 1
	HInfo
	HProgress // 1
	HRequestPending
	HPending
	HComplete // 4
	HRequestComplete
	HCompleteWithBA
)

// GetProgressString mengembalikan string dari enum progress status
func GetProgressString(status int) string {
	switch status {
	case HDataInfo:
		return "Data"
	case HProgress:
		return "Progress"
	case HRequestPending:
		return "Req-Pending"
	case HPending:
		return "Pending"
	case HRequestComplete:
		return "Req-Complete"
	case HCompleteWithBA:
		return "Complete"
	case HComplete:
		return "Complete"
	default:
		return "Unknown"
	}
}
