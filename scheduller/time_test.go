package scheduller

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeInScheduller(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	timeNow := time.Now().In(loc)
	month := timeNow.Month()
	beforeMonth := int(month) - 1
	if beforeMonth == 0 {
		beforeMonth = 12
	}

	// tanggal awal
	timeStart := time.Date(timeNow.Year(), time.Month(beforeMonth), 1, 00, 00, 00, 0, loc)
	timeEnd := timeStart.AddDate(0, 1, 0).Add(time.Second * -1)

	timeStartUnix := timeStart.Unix()
	timeEndUnix := timeEnd.Unix()

	fmt.Println(timeStart)
	fmt.Println(timeEnd)
	fmt.Println(timeStartUnix)
	fmt.Println(timeEndUnix)
}
