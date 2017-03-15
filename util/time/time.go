package time

import (
	"time"
)

// StrToUnixMs converts time string with layout to unix (epoch) value in millis.
func StrToUnixMs(layout, timeStr string) (int, error) {
	t, err := time.Parse(layout, timeStr)

	if err != nil {
		return -1, err
	}

	return int(t.UnixNano() / 1000000), nil
}
