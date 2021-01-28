package wgsconvert

import "math"

// BD09 百度坐标
type BD09 struct {
	Point
}

// ToGCJ02 BD-09 convert to gcj-02
// 百度坐标转火星坐标
func (p BD09) ToGCJ02() (gcj02Point GCJ02) {
	x := p.Longitude - 0.0065
	y := p.Latitude - 0.006

	z := math.Sqrt(p.Longitude*p.Longitude+p.Latitude*p.Latitude) - 0.00002*math.Sin(p.Latitude*xPI)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(p.Longitude*xPI)

	gcj02Point.Longitude = z * math.Cos(theta)
	gcj02Point.Latitude = z * math.Sin(theta)

	return
}

// ToWGS84 BD09 convert to wgs84
// BD09 转换为wgs84
// 具体做法: 先将BD09转为GCJ02, 再把GCJ02转为wgs84
func (p BD09) ToWGS84() WGS84 {
	return p.ToGCJ02().ToWGS84()
}
