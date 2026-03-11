package timeutil

import "time"

const (
	// jstLocationName は日本標準時のタイムゾーン名です。
	jstLocationName = "Asia/Tokyo"
)

// JSTLocation は日本標準時のLocationを返します。
func JSTLocation() *time.Location {
	location, err := time.LoadLocation(jstLocationName)
	if err != nil {
		return time.FixedZone("JST", 9*60*60)
	}
	return location
}

// NowJST は日本標準時の現在時刻を返します。
func NowJST() time.Time {
	return time.Now().In(JSTLocation())
}
