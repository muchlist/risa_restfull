package sfunc

import (
	"testing"
	"time"
)

func TestGetDayName(t *testing.T) {
	timenow := time.Now().Unix()
	hari := GetDayName(timenow)
	println(hari)
}
