/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:13:54
 * @Last Modified by:   Casso-Wong
 * @Last Modified time: 2021-06-05 10:13:54
 */
package distance

import "math"

func GetDistance(lng1, lat1, lng2, lat2 float64) float64 {
	var radius float64 = 6371000 // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}
