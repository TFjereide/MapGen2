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
	continentMaskParams, mountainMaskParams, mountainNoiseParams goNoise.NoiseParams,
	landNoiseParams []goNoise.NoiseParams) *World {

	waterHeight := params.WaterHeight
	fWaterHeight := params.FWaterHeight
	mountainLine := params.MountainLine
	//fMountainLine := params.FMountainLine

	continentMaskParams.Seed = params.Seed + 20
	mountainMaskParams.Seed = params.Seed + 25
	mountainNoiseParams.Seed = params.Seed + 30
	for i := range landNoiseParams {
		landNoiseParams[i].Seed = params.Seed
	}

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

	// Noise Creation
	start := time.Now().UnixNano() / 1000000

	mountainMask := *goNoise.GenFBMParalell(width, height, mountainMaskParams, scale)
	continentMask := *goNoise.GenFBMParalell(width, height, continentMaskParams, scale)

	heightMap := *goNoise.GetLayeredNoise(width, height, landNoiseParams, scale)
	mountainNoise := *goNoise.GenRidgedFBMParalell(width, height, mountainNoiseParams, scale, true, 0.1)
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
		centreDistanceFactor := (1.0 - math.Pow(latdistToCentre, 3)) * (1.0 - math.Pow(longdistToCentre, 4))
		centreDistanceFactor = goNoise.GetGain(clampF64(centreDistanceFactor, 0.0, 1.0), 0.2)
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
	maxDistR := 5
	findShoreDistances(heightMap, &world, maxDistR)
	t = time.Now().UnixNano() / 1000000
	fmt.Println("Distance calc msec: ", t-start)

	// Add Mountains
	for i := range heightMap {
		distToWater := easeOutQuad(float64(world.DistanceToWater[i]) / float64(maxDistR+2))
		heightMap[i] += (mountainNoise[i] * mountainMask[i] * distToWater) * (1.0 - highest_num)
	}

	generateRivers(&world, params, &heightMap)

	//Temperature
	for i := 0; i < len(tempMap); i++ {
		distToCentre := math.Pow(math.Abs(float64(midheight-(i/width)))/float64(midheight), 1.5)

		//centreDistanceFactor := float64(easings.CubicInOut(float32(1.0-distToCentre), 0.0, 0.65, 1.0))
		//centreDistanceFactor := cubicInterpolate(0.6, 0.0, 1.0, 2.1, 1.0-distToCentre)
		centreDistanceFactor := goNoise.GetBias(1.0-distToCentre, 0.1)

		height := heightMap[i]

		if height < fWaterHeight {
			tempMap[i] = goNoise.GetBias((centreDistanceFactor+height)*0.5, 0.97) * 1.0
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
		distWFactor := 0.2
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

	// Moisture
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
			waterdist := easeOutQuad(float64(world.DistanceToWater[i]) / float64(maxDistR))

			distToWaterFactor := 0.25
			noisefactor := 3.0
			tempFactor := 0.2
			heightFactor := 0.3
			moistMap[i] = ((waterdist * distToWaterFactor) + (noise * noisefactor) + (temp * tempFactor) + (height * heightFactor)) / (distToWaterFactor + noisefactor + tempFactor + heightFactor)
			moistMap[i] *= temp
		}
	}

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
		//temp := world.TempMap[i]
		cTemp := (tempMap[i] * float64(temperatureRange)) - float64(MIN_TEMP_c)

		biome := SHRUBLAND

		if height < world.WaterLine {
			if cTemp < -10.0 {
				biome = OCEANIC_GLACIER
			} else if height < (waterHeight - 35) {
				biome = OPEN_SEA
			} else {
				biome = SHALLOW_SEA
			}
		} else if height > world.MountainLine+20 {
			if moist < 180 {
				biome = MOUNTAIN
			} else {
				biome = GLACIER
			}
		} else if height > world.MountainLine {
			if moist < 180 {
				biome = BARE
			} else {
				biome = GLACIER
			}
		} else {

			if cTemp < 0.0 {
				// glacial zone
				if precipitation < 10.0 {
					biome = TEMPERATE_DESERT
				} else if precipitation < 100.0 {
					biome = TUNDRA
				} else {
					biome = GLACIER
				}
			} else if cTemp < 5.0 {
				// boreal zone
				if precipitation < 30.0 {
					biome = SHRUBLAND
				} else {
					biome = TAIGA
				}
			} else if cTemp < 10.0 {
				// warm boreal zone
				if precipitation < 10.0 {
					biome = TEMPERATE_DESERT
				} else if precipitation < 25.0 {
					biome = SHRUBLAND
				} else if precipitation < 200.0 {
					biome = TEMPERATE_DECIDUOUS_FOREST
				} else {
					biome = TEMPERATE_RAINFOREST
				}
			} else if cTemp < 23.0 {
				// temperate zone
				if precipitation < 50.0 {
					biome = TEMPERATE_DESERT
				} else if precipitation < 70.0 {
					biome = SAVANNA
				} else if precipitation < 120.0 {
					biome = SHRUBLAND
				} else if precipitation < 200.0 {
					biome = TEMPERATE_DECIDUOUS_FOREST
				} else if precipitation < 300.0 {
					biome = TEMPERATE_RAINFOREST
				} else {
					biome = WETLAND
				}
			} else {
				// tropical zone
				if precipitation < 100.0 {
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
		world.BiomeMap[i] = biome

	}
	t = time.Now().UnixNano() / 1000000
	fmt.Println("biome map msec: ", t-start)

	return &world

}

func generateRivers(world *World, parameters WorldParams, pheightMap *[]float64) {

	width := parameters.Width
	height := parameters.Height
	heightMap := *pheightMap

	//riverFlowGen
	world.FlowDirection = make([]rl.Vector2, width*height)
	world.FlowFactor = make([]float32, width*height)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {

			idx := getIdx(x, y, width)
			currHeight := heightMap[idx]
			//lowIdx := idx
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
							//lowIdx = _idx
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
