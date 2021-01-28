package wgs84

import (
	"math"
)

// PI π
const pi float64 = math.Pi
const xPI = pi * 3000.0 / 180.0

// 长半径
// 卫星椭球坐标投影到平面地图坐标系的投影因子。
var semiMajorAxis = 6378245.0

// ee: 椭球的偏心率。
// 为什么不是: 2*(1/298.300000000000010000) 或 0.00335281006247
var ee = 0.00669342162296594323

// Point point in GPS Coordinate
type Point struct {
	Longitude float64
	Latitude  float64
}

// OutOfChina 判断是否在国内
// 在国内则需要转化为gcj02, 不再国内不用转换, 直接使用wgs84即可
// true 不在国内
// false 在国内
func (p Point) OutOfChina() bool {
	if p.Longitude < 72.004 || p.Longitude > 137.8347 {
		return true
	}
	if p.Latitude < 0.8293 || p.Latitude > 55.8271 {
		return true
	}
	return false
}

// EARTHRADIUS 地球半径
const EARTHRADIUS = 6378137

// DistanceFrom 点到点的距离
// 以米为单位
func (p Point) DistanceFrom(desp Point) (distance float64) {
	if p.Latitude > 90 || p.Latitude < -90 || desp.Latitude > 90 || desp.Latitude < -90 {
		// err = fmt.Errorf("illegal latitude: p.Latitude %f, desp.Latitude %f\n", p.Latitude, desp.Latitude)
		return
	}
	if p.Longitude > 180 || p.Longitude < -180 || desp.Longitude > 180 || desp.Longitude < -180 {
		// err = fmt.Errorf("illegal longitude: p.Longitude %f, desp.Longitude %f\n", p.Longitude, desp.Longitude)
		return
	}

	var rad = func(coordinate float64) (radian float64) {
		return coordinate * math.Pi / 180.0
	}

	differenceLocLat := rad(p.Latitude) - rad(desp.Latitude)
	differenceLocLng := rad(p.Longitude) - rad(desp.Longitude)

	distance = 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(differenceLocLat/2), 2)+
		math.Cos(rad(p.Latitude))*math.Cos(rad(desp.Latitude))*math.Pow(math.Sin(differenceLocLng/2), 2)))
	distance = distance * EARTHRADIUS
	return
}

// WGS84 wgs84坐标
type WGS84 struct {
	Point
}

// ToGCJ02 wgs84 convert to gcj-02
// wgs84 转换为火星坐标, 高德可用
func (p WGS84) ToGCJ02() (gcjP GCJ02) {
	// var mars_point={lon:0,lat:0};
	if p.OutOfChina() {
		gcjP.Latitude = p.Latitude
		gcjP.Longitude = p.Longitude
		return
	}

	dLat, dLon := delta(p.Latitude, p.Longitude)
	gcjP.Latitude = p.Latitude + dLat
	gcjP.Longitude = p.Longitude + dLon
	/* gcjP.Latitude = p.Latitude
	gcjP.Longitude = p.Longitude

	for i := 0; i < N; i++ {
		dLat, dLon := delta(gcjP.Latitude, gcjP.Longitude)
		gcjP.Latitude = p.Latitude + dLat
		gcjP.Longitude = p.Longitude + dLon
	} */

	return
}

// ToBD09 wgs84 convert to BD09
// wgs84 转换为BD09坐标, 高德可用
// 具体做法: 先将wgs84转GCJ02, 再转BD09
func (p WGS84) ToBD09() BD09 {
	return p.ToGCJ02().ToBD09()
}

// N 转换次数
// 经测试, 转换4次可到较优结果, 转换11次可得最优结果,只转换1次会有一些误差, 保证经度至少确定在0.000001即可
const N = 4

// delta 偏移量
// WGS与GCJ02之间的偏移量
func delta(lat, lon float64) (float64, float64) {
	dLat := transformLat(lon-105.0, lat-35.0)
	dLon := transformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((semiMajorAxis * (1 - ee)) / (magic * sqrtMagic) * pi)
	dLon = (dLon * 180.0) / (semiMajorAxis / sqrtMagic * math.Cos(radLat) * pi)
	return dLat, dLon
}

func transformLat(x, y float64) float64 {
	var ret = -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*pi) + 40.0*math.Sin(y/3.0*pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*pi) + 320*math.Sin(y*pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformLon(x, y float64) float64 {
	var ret = 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*pi) + 40.0*math.Sin(x/3.0*pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*pi) + 300.0*math.Sin(x/30.0*pi)) * 2.0 / 3.0
	return ret
}

// WGSToBeijing54

// WGSToXian80
