package hwlist

const (
	win1064 = "Windows 10 64"
	win1032 = "Windows 10 32"
	win864  = "Windows 8 64"
	win832  = "Windows 8 32"
	win764  = "Windows 7 64"
	win732  = "Windows 7 32"
	winXp   = "Windows XP 32"
	winServ = "Windows Server"
	mac     = "Mac Os"
	linux   = "Linux"

	osC2duo = "Core 2 Duo"
	osi3    = "Core i3"
	osi5    = "Core i5"
	osi7    = "Core i7"
	osAmd   = "AMD"

	ram0     = 0
	ram1000  = 1000
	ram2000  = 2000
	ram3000  = 3000
	ram4000  = 4000
	ram6000  = 6000
	ram8000  = 8000
	ram12000 = 1200
	ram16000 = 1600
	ram32000 = 3200

	hdd0     = 0
	hdd128   = 128
	hdd256   = 256
	hdd512   = 512
	hdd1000  = 1000
	hdd2000  = 2000
	hdd3000  = 3000
	hdd4000  = 4000
	hdd8000  = 8000
	hdd16000 = 16000
	hdd32000 = 32000
)

func GetPCOSAvailable() []string {
	return []string{
		win1064,
		win1032,
		win864,
		win832,
		win764,
		win732,
		winXp,
		winServ,
		mac,
		linux,
		other,
	}
}

func GetPCProcessor() []string {
	return []string{
		osC2duo,
		osi3,
		osi5,
		osi7,
		osAmd,
	}
}

func GetPCRam() []int {
	return []int{
		ram0,
		ram1000,
		ram2000,
		ram3000,
		ram4000,
		ram6000,
		ram8000,
		ram12000,
		ram16000,
		ram32000,
	}
}
func GetPCHDD() []int {
	return []int{
		hdd0,
		hdd128,
		hdd256,
		hdd512,
		hdd1000,
		hdd2000,
		hdd3000,
		hdd4000,
		hdd8000,
		hdd16000,
		hdd32000,
	}
}
