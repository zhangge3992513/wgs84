package wgsconvert

import (
	"math"
	"reflect"
	"testing"
)

func TestGCJ02_ToWGS84(t *testing.T) {
	// go test -v -timeout 30s -run ^TestPoint_GCJ02ToWGS84$ github.com/zhangge3992513/wgsconvert

	tests := []struct {
		name      string
		p         GCJ02
		wantWgs84 WGS84
	}{
		// TODO: Add test cases.
		{
			name: "cgj-02坐标转1换为高德wgs84",
			p: GCJ02{
				Point{
					Longitude: 113.942831, // wgs84: 113.93795567290164   113.93796085386425  113.93796084819188 113.93796084819807 113.93796084819807
					Latitude:  22.530282,  //wgs84: 22.533298820468712
				},
			},
			wantWgs84: WGS84{
				Point{
					Longitude: 113.937961,
					Latitude:  22.533307,
				},
			},
		},
		{
			name: "cgj-02新疆乌恰县换为wgs84",
			p: GCJ02{
				Point{
					Longitude: 75.25905448676403, //
					Latitude:  39.71932164456469,
				},
			},
			wantWgs84: WGS84{
				Point{
					Longitude: 75.255949, // √
					Latitude:  39.719029, // √
				},
			},
		},
		{
			name: "cgj-02黑龙江佳木斯抚远镇wgs84",
			p: GCJ02{
				Point{
					Longitude: 134.30364464514864, //  134.29595,48.36794,
					Latitude:  48.3703309490954,
				},
			},
			wantWgs84: WGS84{
				Point{
					Longitude: 134.29595, // √ 结果: 134.29595000000018
					Latitude:  48.36794,  // √ 结果: 48.36794000000012
				},
			},
		},
	}
	t.Log(math.Pi)
	t.Log(math.Pi + 1.1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWgs84Point := tt.p.ToWGS84(); !reflect.DeepEqual(gotWgs84Point, tt.wantWgs84) {
				t.Errorf("Point.GCJ02ToWGS84() = %v, want %v", gotWgs84Point, tt.wantWgs84)
			}
		})
	}
}

// 11次 BenchmarkPoint_GCJ02ToWGS84-4   	  293467	      3643 ns/op	       0 B/op	       0 allocs/op	       0 allocs/op
// 4 次 BenchmarkPoint_GCJ02ToWGS84-4         925525	      1445 ns/op	       0 B/op	       0 allocs/op
// go test -benchmem -run=^$ -bench ^(BenchmarkPoint_GCJ02ToWGS84)$ github.com/zhangge3992513/wgsconvert
func BenchmarkGCJ02_ToWGS84(t *testing.B) {
	p := GCJ02{
		Point{
			Longitude: 134.30364464514864, //  134.29595,48.36794,
			Latitude:  48.3703309490954,
		},
	}
	t.ResetTimer()
	t.StartTimer()
	for i := 0; i < t.N; i++ {
		p.ToWGS84()
	}
	t.StopTimer()
}

func TestWGS84_ToGCJ02(t *testing.T) {
	tests := []struct {
		name     string
		p        WGS84
		wantGcjP GCJ02
	}{
		// TODO: Add test cases.
		{
			name: "wgs84坐标转1换为高德cgj-02",
			p: WGS84{
				Point{
					Longitude: 116.46706996, // cgj: 116.4732071375548
					Latitude:  39.99188446,  // cgj: 39.9932029943646
				},
			},
			wantGcjP: GCJ02{
				Point{
					Longitude: 116.4732071375548,
					Latitude:  39.9932029943646,
				},
			},
		}, {
			name: "wgs84坐标转2换为高德cgj-02",
			p: WGS84{
				Point{
					Longitude: 113.937961, // cgj: 113.94283115195938  gps: 113.937961
					Latitude:  22.533307,  // cgj: 22.53028209851657   gps: 22.533307
				},
			},
			wantGcjP: GCJ02{
				Point{
					Longitude: 113.94283115195938,
					Latitude:  22.53028209851657,
				},
			},
		}, {
			name: "黑龙江乌恰县",
			p: WGS84{
				Point{
					Longitude: 75.255949, //
					Latitude:  39.719029, //
				},
			},
			wantGcjP: GCJ02{
				Point{
					Longitude: 75.25905680338599,
					Latitude:  39.719323459202,
				},
			},
		},
		{
			name: "黑龙江佳木斯抚远镇", // 结果: 134.30364464514864,48.3703309490954, 偏差600多米!
			p: WGS84{
				Point{
					Longitude: 134.29595,
					Latitude:  48.36794,
				},
			},
			wantGcjP: GCJ02{
				Point{
					Longitude: 134.311357,
					Latitude:  48.372734,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotGcjP := tt.p.ToGCJ02(); !reflect.DeepEqual(gotGcjP, tt.wantGcjP) {
				t.Errorf("WGS84.ToGCJ02() = %v, want %v", gotGcjP, tt.wantGcjP)
			}
		})
	}
}

func TestPoint_DistanceFrom(t *testing.T) {
	type args struct {
		desp Point
	}
	tests := []struct {
		name         string
		p            Point
		args         args
		wantDistance float64
	}{
		// TODO: Add test cases.
		{
			name: "点1与点2之间的距离",
			p: Point{
				Longitude: 113.00458767361198,
				Latitude:  9.998094075521,
			},
			args: args{
				desp: Point{
					Longitude: 114.00443684895902,
					Latitude:  11.998577473959,
				},
			},
			wantDistance: 248048.606567,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDistance := tt.p.DistanceFrom(tt.args.desp); gotDistance != tt.wantDistance {
				t.Errorf("Point.DistanceFrom() = %v, want %v", gotDistance, tt.wantDistance)
			}
		})
	}
}
