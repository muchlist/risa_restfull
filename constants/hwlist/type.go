package hwlist

const (
	cctvFixed = "Fixed"
	cctvPTZ   = "PTZ"

	pcDesktop  = "Desktop"
	pcLaptop   = "Laptop"
	pcAllInOne = "All in one"
	pcServer   = "Server"
	pcMini     = "Mini PC"

	other = "Lainnya"
)

func GetCctvTypeAvailable() []string {
	return []string{cctvFixed, cctvPTZ, other}
}

func GetComputerTypeAvailable() []string {
	return []string{pcDesktop, pcLaptop, pcAllInOne, pcServer, pcMini, other}
}
