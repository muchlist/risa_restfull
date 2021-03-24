package hwlist

const (
	cctvFixed = "Fixed"
	cctvPTZ   = "PTZ"
)

func GetCctvTypeAvailable() []string {
	return []string{cctvFixed, cctvPTZ}
}
