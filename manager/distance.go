package main

import (
	"math"
)

const (
	DEF_PI    = 3.14159265359
	DEF_2PI   = 6.28318530712
	DEF_PI180 = 0.01745329252
	DEF_R     = 6370693.5
)

func GetShortDistance(lon1 float64, lat1 float64, lon2 float64, lat2 float64) float64 {
	var ew1 float64
	var ns1 float64
	var ew2 float64
	var ns2 float64
	var dx float64
	var dy float64
	var dew float64
	var distan float64

	ew1 = lon1 * DEF_PI180
	ns1 = lat1 * DEF_PI180
	ew2 = lon2 * DEF_PI180
	ns2 = lat2 * DEF_PI180

	dew = ew1 - ew2
	if dew > DEF_PI {
		dew = DEF_2PI - dew
	} else if dew < -DEF_PI {
		dew = DEF_2PI + dew
	}

	dx = DEF_R * math.Cos(ns1) * dew
	dy = DEF_R * (ns1 - ns2)
	distan = math.Sqrt(dx*dx + dy*dy)

	return distan
}

func GetLongDistance(lon1 float64, lat1 float64, lon2 float64, lat2 float64) float64 {
	var ew1 float64
	var ns1 float64
	var ew2 float64
	var ns2 float64

	var distan float64

	ew1 = lon1 * DEF_PI180
	ns1 = lat1 * DEF_PI180
	ew2 = lon2 * DEF_PI180
	ns2 = lat2 * DEF_PI180

	distan = math.Sin(ns1) * math.Sin(ns2) * math.Cos(ns2) * math.Cos(ew1-ew2)
	if distan > 1.0 {
		distan = 1.0
	} else if distan < -1.0 {
		distan = -1.0
	}

	distan = DEF_R * math.Acos(distan)

	return distan
}
