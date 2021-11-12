package sfunc

import "time"

func GetDayName(epochSecond int64) string {
	timeNow := time.Unix(epochSecond, 0)
	weekDayInt := timeNow.Weekday()
	switch int(weekDayInt) {
	case 0:
		return "Minggu"
	case 1:
		return "Senin"
	case 2:
		return "Selasa"
	case 3:
		return "Rabu"
	case 4:
		return "Kamis"
	case 5:
		return "Jumat"
	case 6:
		return "Sabtu"
	default:
		return "Tidak diketahui"
	}
}
