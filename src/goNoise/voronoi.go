package goNoise

import (
	"math"
	"math/rand"
)

type voronoiPoint struct {
	X, Y int
}

func (vp *voronoiPoint) DistanceTo(x, y int) int {
	return int(math.Sqrt(float64((vp.X-x)*(vp.X-x) + (vp.Y-y)*(vp.Y-y))))
}

// Generates a flat voronoi map,
func GenVoronoiCells(width, height, numPoints int) []byte {

	rand.Seed(1)
	hm := make([]byte, width*height)

	vPoints := make([]voronoiPoint, numPoints)
	vHeights := make([]byte, numPoints)
	for i := 0; i < numPoints; i++ {
		vPoints[i] = voronoiPoint{rand.Intn(width), rand.Intn(height)}
	}

	for i, _ := range vHeights {
		//flag := rand.Intn(255)
		//vHeights[i] = byte(flag)
		//continue
		if (i % 2) == 0 {
			vHeights[i] = 255
		} else {
			vHeights[i] = 50
		}

	}
	for i, _ := range hm {
		currPosX, currPosY := getPoint(i, width)
		idx := 0
		prevDist := vPoints[0].DistanceTo(currPosX, currPosY)
		for j := 1; j < numPoints; j++ {
			dist := vPoints[j].DistanceTo(currPosX, currPosY)
			if dist < prevDist {
				prevDist = dist
				idx = j
			}
		}
		hm[i] = byte(vHeights[idx])
	}

	// for i, _ := range hm {
	// 	currPosX, currPosY := getPoint(i, width)
	// 	idx := 0
	// 	prevDist := vPoints[0].DistanceTo(currPosX, currPosY)
	// 	for j := 1; j < numPoints; j++ {
	// 		dist := vPoints[j].DistanceTo(currPosX, currPosY)
	// 		if dist < prevDist {
	// 			prevDist = dist
	// 			idx = j
	// 		}
	// 	}
	// 	hm[i] = vHeights[idx]

	// 	fdist := byteToF64(byte(
	// 		clampInt(prevDist, 255)),
	// 	)

	// 	hm[i] = f64Tobyte(fdist * fdist)
	// }

	return hm
}

func getPoint(idx, width int) (int, int) {
	return int(idx % width), int(idx / width)
}

func byteToF64(b byte) float64 {

	return float64(b) * 255
}
func f64Tobyte(f float64) byte {

	return byte(f / 255)
}
func clampInt(i, max int) int {
	if i > max {
		return max
	}
	return i
}
