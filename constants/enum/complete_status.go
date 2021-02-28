package enum

const (
	HInfo = iota
	HProgress
	HRequestPending
	HPending
	HComplete
)

func GetProgressString(status int) string {
	listCStatus := make([]string, 5)
	listCStatus[HInfo] = "Info"
	listCStatus[HProgress] = "Progress"
	listCStatus[HRequestPending] = "Pending"
	listCStatus[HPending] = "Pending"
	listCStatus[HComplete] = "Complete"

	return listCStatus[status]
}
