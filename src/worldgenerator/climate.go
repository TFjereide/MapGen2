package worldgenerator

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

func generateMoisture(pheightMap, ptempMap *[]float64, width, height int) *[]float64 {

	headings := []rl.Vector2{
		rl.NewVector2(0, 1),
		rl.NewVector2(1, -1),
		rl.NewVector2(-1, 1),
		rl.NewVector2(-1, -1),
		rl.NewVector2(1, -1),
		rl.NewVector2(0, -1),
	}

	heightMap := *pheightMap
	tempMap := *ptempMap
	moistureMap := make([]float64, len(tempMap))
	windMap := make([]rl.Vector2, len(tempMap))
	//windStrengthMap := make([]float64, len(tempMap))

	// setup big wind
	for i, _ := range heightMap {
		_, py := GetPoint(i, width)
		hidx := int(math.Min(float64(py/(height/len(headings))), 5.0))
		windMap[i] = headings[hidx]

		//Debug
		moistureMap[i] = float64(hidx) / float64(len(headings))

	}
	//tick through

	return &moistureMap
}
func getNext(idx, width int, vec rl.Vector2) int {
	x, y := GetPointI32(idx, width)
	dx := int32(vec.X)
	dy := int32(vec.Y)
	return getIdx(int(x+dx), int(y+dy), width)
}
