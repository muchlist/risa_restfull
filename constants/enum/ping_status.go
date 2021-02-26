package enum

const (
	PingDown = iota
	PingHalf
	PingUp
)

func GetPingString(status int) string {
	listPingStatus := make([]string, 3)
	listPingStatus[PingDown] = "DOWN"
	listPingStatus[PingHalf] = "HALF"
	listPingStatus[PingUp] = "UP"

	return listPingStatus[status]
}
