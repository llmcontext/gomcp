package mcp_server_time

import "time"

func isDst(location *time.Location) bool {
	_, offset := time.Now().In(location).Zone()
	standardTime := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, location) // A known non-DST date
	_, standardOffset := standardTime.Zone()
	return offset != standardOffset
}
