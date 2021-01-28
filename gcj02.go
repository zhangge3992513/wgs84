package wgs84

import "math"

// GCJ02 火星坐标
type GCJ02 struct {
	Point
}

// ToBD09 gcj-02 convert to BD-09
// 火星坐标转百度坐标
func (p GCJ02) ToBD09() (bd09Point BD09) {
	var z = math.Sqrt(p.Longitude*p.Longitude+p.Latitude*p.Latitude) + 0.00002*math.Sin(p.Latitude*xPI)
	var theta = math.Atan2(p.Latitude, p.Longitude) + 0.000003*math.Cos(p.Longitude*xPI)
	bd09Point.Longitude = z*math.Cos(theta) + 0.0065
	bd09Point.Latitude = z*math.Sin(theta) + 0.006
	return
}

// ToWGS84 gcj-02 convert to wgs84
// 火星坐标转wgs84
func (p GCJ02) ToWGS84() (wgs84Point WGS84) {
	wgs84Point.Latitude = p.Latitude
	wgs84Point.Longitude = p.Longitude
	// 多次求偏移量得到结果基本就是最佳结果
	for i := 0; i < N; i++ {
		dLat, dLon := delta(wgs84Point.Latitude, wgs84Point.Longitude)
		wgs84Point.Latitude = p.Latitude - dLat
		wgs84Point.Longitude = p.Longitude - dLon
	}
	return
}
