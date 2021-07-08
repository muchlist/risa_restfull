package timegen

import (
	"time"
)

func GetTimeWITA(timestampSec int64) (string, error) {
	witaTimeZone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		return "", err
	}
	return time.Unix(timestampSec, 0).In(witaTimeZone).Format("02 Jan 15:04"), nil
}

func GetTimeWithYearWITA(timestampSec int64) (string, error) {
	witaTimeZone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		return "", err
	}
	return time.Unix(timestampSec, 0).In(witaTimeZone).Format("02 Jan 2006 15:04"), nil
}

func GetTimeAsName(timestampSec int64) (string, error) {
	witaTimeZone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		return "", err
	}
	return time.Unix(timestampSec, 0).In(witaTimeZone).Format("02-01-2006-15-04"), nil
}

func GetHourWITA(timestampSec int64) (string, error) {
	witaTimeZone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		return "", err
	}
	return time.Unix(timestampSec, 0).In(witaTimeZone).Format("15:04"), nil
}
