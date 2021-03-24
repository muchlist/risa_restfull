package statuses

const (
	Enable  = "enable"
	Disable = "disable"
)

func GetStatusAvailable() []string {
	return []string{Enable, Disable}
}
