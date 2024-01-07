package common

import "time"

func GetCurrentYear() uint16 {
	year, _, _ := time.Now().Date()
	return uint16(year)
}