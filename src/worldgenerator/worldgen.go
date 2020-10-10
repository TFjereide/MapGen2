package worldgenerator

import (
	common "common"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github/raylib-go-master/easings"

	goNoise "goNoise"
	"math"
	"time"
)

// Generates a proper world struct in several stages
func GenerateWorld(
	params WorldParams,
	continentMaskParams, mountainMaskParams goNoise.NoiseParams,
	mountainNoiseParams, landNoiseParams []goNoise.NoiseParams) *World {

	// Set seeds for all parameters. It needs to be derived from one master seed and then offset
	continentMaskParams.Seed = params.Seed + 20
	mountainMaskParams.Seed = params.Seed + 25
	for i := range mountainNoiseParams {
		mountainNoiseParams[i].Seed = params.Seed + 2
	}
	for i := range landNoiseParams {
		landNoiseParams[i].Seed = params.Seed + 1
	}
	waterHeight := params.WaterHeight
	fWaterHeight := params.FWaterHeight
	mountainLine := params.MountainLine
	//fMountainLine := params.FMountainLine
	width := params.Width
	height := params.Height
	scale := params.Scale
	midheight := height / 2
	midWidth := width / 2
	temperatureRange := MAX_TEMP_C + MIN_TEMP_c

	world := World{
		X:               width,
		Y:               height,
		WaterLine:       waterHeight,
		MountainLine:    mountainLine,
		Equator:         height / 2,
		DistanceToLand:  make([]byte, width*height),
		DistanceToWater: make([]byte, width*height)}

	// Noise map generation
	start := time.Now().UnixNano() / 1000000

	mountainMask := *goNoise.GenFBMParalell(width, height, mountainMaskParams, scale)
	continentMask := *goNoise.GenFBMParalell(width, height, continentMaskParams, scale)

	heightMap := *goNoise.GetLayeredNoise(width, height, landNoiseParams, scale)
	mountainNoise := *goNoise.GetLayeredNoise(width, height, mountainNoiseParams, scale)
	tempMap := *goNoise.GenFBMParalell(width, height, landNoiseParams[1], scale)
	moistMap := *goNoise.GenFBMParalell(width, height, goNoise.NoiseParams{
		ScaleX:      20.0,
		ScaleY:      15.0,
		Octaves:     5,
		Persistance: 0.5,
		Lacunarity:  2.5,
		Gain:        0.3,
		Bias:        0.4,
		Seed:        params.Seed + 7},
		scale)

	t := time.Now().UnixNano() / 1000000
	fmt.Println("Noise Maps creation msec: ", t-start)

	// Obtain Continents with minor noise
	highest_num := 0.0
	for i := 0; i < len(heightMap); i++ {
		latdistToCentre := math.Abs(float64(midheight-(i/width))) / float64(midheight)
		longdistToCentre := math.Abs(float64(midWidth-(i%width))) / float64(midWidth)

		// Obtain mask of "Strength" based on distance from edges
		centreDistanceFactor := (1.0 - math.Pow(latdistToCentre, 9)) * (1.0 - math.Pow(longdistToCentre, 9))
		centreDistanceFactor = goNoise.GetGain(clampF64(centreDistanceFactor, 0.0, 1.0), 0.4)
		heightMap[i] = heightMap[i] * centreDistanceFactor * continentMask[i]
		if heightMap[i] > fWaterHeight {
			heightMap[i] -= fWaterHeight
			heightMap[i] *= 0.33
			heightMap[i] += fWaterHeight
		}
		if heightMap[i] > highest_num {
			highest_num = heightMap[i]
		}

	}

	// Find distance to water and land
	start = time.Now().UnixNano() / 1000000
	maxDistR := 8
	findShoreDistances(heightMap, &world, maxDistR)
	t = time.Now().UnixNano() / 1000000
	fmt.Println("Distance calc msec: ", t-start)

	// Add Mountains
	for i := range heightMap {
		distToWater := easeOutQuad(float64(world.DistanceToWater[i]) / float64(maxDistR+2))
		heightMap[i] += (mountainNoise[i] * mountainMask[i] * distToWater) * (1.0 - highest_num)
	}

	//Temperature
	for i := 0; i < len(tempMap); i++ {
		heatCentre := height - (height / 5 * 2)
		distToCentre := math.Pow(math.Abs(float64(heatCentre-(i/width)))/float64(heatCentre), 2.0)
		//distToCentre := math.Pow(math.Abs(float64(midheight-(i/width)))/float64(midheight), 1.1)

		//centreDistanceFactor := float64(easings.CubicInOut(float32(1.0-distToCentre), 0.0, 0.65, 1.0))
		//centreDistanceFactor := cubicInterpolate(0.6, 0.0, 1.0, 2.1, 1.0-distToCentre)
		centreDistanceFactor := goNoise.GetBias(1.0-distToCentre, 0.1)

		height := heightMap[i]

		if height < fWaterHeight {
			tempMap[i] = ((goNoise.GetBias((centreDistanceFactor+height)*0.5, 0.97) * 0.7) + moistMap[i]*0.3)
			continue
		}

		distToWater := easeOutQuad(float64(world.DistanceToWater[i]) / float64(maxDistR))
		//height = clampF64(float64(easings.CubicInOut(float32(
		//	1.0-(height-fWaterHeight)*(1.0/fWaterHeight)),
		//	0.0,
		//	1.0,
		//	1.0)), 0.0, 1.0)
		height = 1.0 - (height-fWaterHeight)*(1.0/fWaterHeight)

		centreDistanceFactor = goNoise.GetGain(centreDistanceFactor, 0.35)
		noiseFactor := 0.2
		distWFactor := 0.4
		hFactor := 0.8

		//temp := centreDistanceFactor +
		//	(distToWater * distWFactor) +
		//	(tempMap[i] * noiseFactor)

		//tempMap[i] = (temp / (1.0 + distWFactor + noiseFactor)) * height

		temp := centreDistanceFactor +
			(height * hFactor) +
			(distToWater * distWFactor) +
			(tempMap[i] * noiseFactor)

		tempMap[i] = (temp / (1.0 + hFactor + distWFactor + noiseFactor))

	}

	//moistMap = *generateMoisture(&heightMap, &tempMap, width, height)

	//Moisture
	for i := 0; i < len(moistMap); i++ {
		if heightMap[i] < fWaterHeight {
			moistMap[i] = 1.0
		} else {
			noise := moistMap[i]
			temp := 1.0 - tempMap[i]
			height := heightMap[i]
			height = clampF64(float64(easings.CubicOut(float32(
				1.0-(height-fWaterHeight)*(1.0/fWaterHeight)),
				0.0,
				0.2,
				1.0)), 0.0, 1.0)
			waterdist := easeOutQuad(1.0 - float64(world.DistanceToWater[i])/float64(maxDistR))

			distToWaterFactor := 0.6
			noisefactor := 3.0
			tempFactor := 0.5
			heightFactor := 0.3
			moistMap[i] = ((waterdist * distToWaterFactor) + (noise * noisefactor) + (temp * tempFactor) + (height * heightFactor)) / (distToWaterFactor + noisefactor + tempFactor + heightFactor)
			moistMap[i] *= temp
		}
	}

	generateWatersheds(&world, params, &heightMap, &moistMap)
	generateNormals(&world, params, &heightMap)

	start = time.Now().UnixNano() / 1000000

	world.HeightMap = common.ConvertArrayFloat64ToByte(heightMap)
	world.MoistureMap = common.ConvertArrayFloat64ToByte(moistMap)
	world.TempMap = common.ConvertArrayFloat64ToByte(tempMap)
	world.BiomeMap = make([]BiomeType, world.X*world.Y)

	t = time.Now().UnixNano() / 1000000
	fmt.Println("mapconversion msec: ", t-start)

	start = time.Now().UnixNano() / 1000000

	//Biome
	for i := 0; i < (world.X * world.Y); i++ {
		height := world.HeightMap[i]
		moist := world.MoistureMap[i]
		precipitation := (float64(moist) / 255.0) * float64(MAX_PRECIPITATION)
		cTemp := (tempMap[i] * float64(temperatureRange)) - float64(MIN_TEMP_c)
		world.BiomeMap[i] = determineBiome(precipitation, cTemp, height, world.WaterLine, world.MountainLine)

	}
	t = time.Now().UnixNano() / 1000000
	fmt.Println("biome map msec: ", t-start)

	return &world

}

func determineBiome(precipitation, cTemp float64, height, waterLine, mountainLine byte) BiomeType {

	biome := SHRUBLAND
	if height < waterLine {
		if cTemp < -10.0 {
			biome = OCEANIC_GLACIER
		} else if height < (waterLine - 35) {
			biome = OPEN_SEA
		} else {
			biome = SHALLOW_SEA
		}
	} else if height > (255 - 20) {
		if precipitation > 100.0 && cTemp < 1.0 {
			biome = GLACIER
		} else {
			biome = MOUNTAIN
		}
	} else if height > mountainLine {
		if precipitation > 100.0 && cTemp < 1.0 {
			biome = GLACIER
		} else {
			biome = BARE
		}
	} else {

		if cTemp < -2.0 {
			// glacial zone
			if precipitation < 10.0 {
				biome = TEMPERATE_DESERT
			} else if precipitation > 150.0 && height > mountainLine-20 {
				biome = GLACIER
			} else {
				biome = TUNDRA
			}
		} else if cTemp < 5.0 {
			// boreal zone
			if precipitation < 30.0 {
				biome = SHRUBLAND
			} else {
				biome = TAIGA
			}
		} else if cTemp < 7.0 {
			// warm boreal zone
			if precipitation < 10.0 {
				biome = TEMPERATE_DESERT
			} else if precipitation < 50.0 {
				biome = SHRUBLAND
			} else if precipitation < 200.0 {
				biome = TEMPERATE_DECIDUOUS_FOREST
			} else {
				biome = TEMPERATE_RAINFOREST
			}
		} else if cTemp < 25.0 {
			// temperate zone
			if precipitation < 25.0 {
				biome = TEMPERATE_DESERT
			} else if precipitation < 45.0 {
				biome = SAVANNA
			} else if precipitation < 70.0 {
				biome = GRASSLAND
			} else if precipitation < 120.0 {
				biome = SHRUBLAND
			} else if precipitation < 250.0 {
				biome = TEMPERATE_DECIDUOUS_FOREST
			} else if precipitation < 320.0 {
				biome = TEMPERATE_RAINFOREST
			} else {
				biome = WETLAND
			}
		} else {
			// tropical zone
			if precipitation < 80.0 {
				biome = SUBTROPICAL_DESERT
			} else if precipitation < 150.0 {
				biome = SAVANNA
			} else if precipitation < 200.0 {
				biome = TROPICAL_SEASONAL_FOREST
			} else {
				biome = TROPICAL_RAINFOREST
			}
		}
	}
	return biome
}

func generateNormals(world *World, parameters WorldParams, pheightMap *[]float64) {

	width := parameters.Width
	height := parameters.Height
	heightMap := *pheightMap
	//fWaterHeight := parameters.FWaterHeight

	//riverFlowGen
	world.FlowDirection = make([]rl.Vector2, width*height)
	world.FlowFactor = make([]float32, width*height)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {

			idx := getIdx(x, y, width)
			currHeight := heightMap[idx]
			lowHeight := currHeight
			dir := rl.Vector2{0, 0}

			//Find lowestPoint adjacent
			for dy := -1; dy < 2; dy++ {
				for dx := -1; dx < 2; dx++ {
					if dx != 0 && dy != 0 {

						if x+dx < 0 || x+dx >= width || y+dy < 0 || y+dy >= height {
							continue
						}
						_idx := getIdx(x+dx, y+dy, width)
						checkHeight := heightMap[_idx]
						if checkHeight < lowHeight {
							lowHeight = checkHeight
							dir.X = float32(dx)
							dir.Y = float32(dy)
						}
					}
				}
			}
			world.FlowDirection[idx] = dir
			world.FlowFactor[idx] = float32(currHeight - lowHeight)

		}
	}

}

func getIdxNeighbours(idx, width int) []int {
	x, y := GetPointI32(idx, width)

	return []int{
		getIdx(int(x-1), int(y), width),
		getIdx(int(x), int(y-1), width),
		getIdx(int(x), int(y+1), width),
		getIdx(int(x+1), int(y), width),
	}

}
func getIdxNeighboursDiagonal(idx, width int) []int {
	point := NewVector2Int32(GetPointI32(idx, width))

	return []int{
		getIdx(int(point.X-1), int(point.Y-1), width),
		getIdx(int(point.X-1), int(point.Y), width),
		getIdx(int(point.X-1), int(point.Y+1), width),
		getIdx(int(point.X), int(point.Y-1), width),
		getIdx(int(point.X), int(point.Y+1), width),
		getIdx(int(point.X+1), int(point.Y-1), width),
		getIdx(int(point.X+1), int(point.Y), width),
		getIdx(int(point.X+1), int(point.Y+1), width),
	}

}
func getNeighbours(point Vector2Int32) []Vector2Int32 {

	return []Vector2Int32{
		NewVector2Int32(point.X-1, point.Y-1),
		NewVector2Int32(point.X-1, point.Y),
		NewVector2Int32(point.X-1, point.Y+1),
		NewVector2Int32(point.X, point.Y-1),
		NewVector2Int32(point.X, point.Y+1),
		NewVector2Int32(point.X+1, point.Y-1),
		NewVector2Int32(point.X+1, point.Y),
		NewVector2Int32(point.X+1, point.Y+1),
	}

}
func getNeighboursIdx(point Vector2Int32, width int) []int {
	return []int{
		getIdx(int(point.X-1), int(point.Y-1), width),
		getIdx(int(point.X-1), int(point.Y), width),
		getIdx(int(point.X-1), int(point.Y+1), width),
		getIdx(int(point.X), int(point.Y-1), width),
		getIdx(int(point.X), int(point.Y+1), width),
		getIdx(int(point.X+1), int(point.Y-1), width),
		getIdx(int(point.X+1), int(point.Y), width),
		getIdx(int(point.X+1), int(point.Y+1), width),
	}

}
