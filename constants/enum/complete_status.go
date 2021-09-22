package enum

const (
	HDataInfo = iota - 1
	HInfo
	HProgress
	HRequestPending
	HPending
	HComplete
)

func GetProgressString(status int) string {

	if status == HDataInfo {
		return "Data"
	}

	listCStatus := make([]string, 6)
	listCStatus[HInfo] = "Info"
	listCStatus[HProgress] = "Progress"
	listCStatus[HRequestPending] = "Pending"
	listCStatus[HPending] = "Pending"
	listCStatus[HComplete] = "Complete"

	return listCStatus[status]
}
