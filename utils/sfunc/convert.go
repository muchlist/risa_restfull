package sfunc

import (
	"fmt"
	"strconv"
	"time"
)

func StrToInt(text string, defaultReturn int) int {
	number := defaultReturn
	if text != "" {
		var err error
		number, err = strconv.Atoi(text)
		if err != nil {
			number = defaultReturn
		}
	}
	return number
}

func IntToTime(second int64, returnIfZero string) string {
	if second <= 0 {
		return returnIfZero
	}
	var timeNeed time.Duration
	timeNeed = time.Duration(second)
	return fmt.Sprint(timeNeed * time.Second)
}

func IntToDateIndoFormat(second int64, returnIfZero string) string {
	if second <= 0 {
		return returnIfZero
	}
	return fmt.Sprint(time.Unix(second, 0).Format("02-01-2006"))
}
