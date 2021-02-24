package enum

const (
	HProgress = iota + 1
	HRequestPending
	HPending
	HComplete
)

func GetProgressString(status int) string {
	listCStatus := make([]string, 5)
	listCStatus[HProgress] = "Progress"
	listCStatus[HRequestPending] = "Pending"
	listCStatus[HPending] = "Pending"
	listCStatus[HComplete] = "Complete"

	return listCStatus[status]
}
