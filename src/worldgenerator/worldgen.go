package worldgenerator

import (
	common "common"
	"fmt"
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

	mountainMask := *goNoise.GenFBM(width, height, mountainMaskParams, scale)
	continentMask := *goNoise.GenFBM(width, height, continentMaskParams, scale)

	heightMap := *goNoise.GetLayeredNoise(width, height, landNoiseParams, scale)
	mountainNoise := *goNoise.GenRidgedFBM(width, height, mountainNoiseParams, scale, true)
	tempMap := *goNoise.GenFBM(width, height, landNoiseParams[1], scale)
	moistMap := *goNoise.GenFBM(width, height, goNoise.NoiseParams{
		ScaleX:      30.0,
		ScaleY:      15.0,
		Octaves:     4,
		Persistance: 0.5,
		Lacunarity:  2.0,
		Gain:        0.3,
		Bias:        0.3,
		Seed:        params.Seed + 7},
		scale)

	t := time.Now().UnixNano() / 1000000
	fmt.Println("Noise Maps creation msec: ", t-start)
	start = time.Now().UnixNano() / 1000000

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
			heightMap[i] *= 0.5
			heightMap[i] += fWaterHeight

		}
		if heightMap[i] > highest_num {
			highest_num = heightMap[i]
		}

	}
	fmt.Println("Highest num is ", highest_num)

	t = time.Now().UnixNano() / 1000000
	fmt.Println("Height operations msec: ", t-start)

	start = time.Now().UnixNano() / 1000000

	// Find distance to water and land
	maxDistR := 8
	for idx := range heightMap {
		x, y := getPoint(idx, width)
		currElev := heightMap[idx]

		//Water or land value for tile that is checked
		water := true

		if currElev > fWaterHeight {
			water = false
			world.DistanceToLand[idx] = 0
			world.DistanceToWater[idx] = byte(maxDistR)
		} else {
			world.DistanceToWater[idx] = 0
			world.DistanceToLand[idx] = byte(maxDistR)
		}

	distanceFound:

		// Check within a radius of r
		for r := 1; r < maxDistR; r++ {

			for dy := -r; dy <= r; dy++ {
				for dx := -r; dx <= r; dx++ {

					if dx < r && dy < r && dx > -r && dy > -r {
						continue
					}

					didx := getIdx((x + dx), (y + dy), width)
					if didx < 0 || didx > len(heightMap)-1 {
						continue
					}

					locElev := heightMap[didx]

					if water {
						// Check  found land
						if locElev > fWaterHeight {
							world.DistanceToLand[idx] = byte(r)
							break distanceFound
						}

					} else {
						//Check found water
						if locElev < fWaterHeight {
							world.DistanceToWater[idx] = byte(r)
							break distanceFound
						}
					}

				}
			}
		}
	}

	for i := range heightMap {
		distToWater := easeOutQuad(float64(world.DistanceToWater[i]) / float64(maxDistR+2))
		heightMap[i] += (mountainNoise[i] * mountainMask[i] * distToWater) * 0.4
		//mountainNoise[i] = mountainMask[i]
	}

	t = time.Now().UnixNano() / 1000000
	fmt.Println("Distance calc msec: ", t-start)
	start = time.Now().UnixNano() / 1000000

	//Temperature
	for i := 0; i < len(tempMap); i++ {
		distToCentre := math.Pow(math.Abs(float64(midheight-(i/width)))/float64(midheight), 2)
		centreDistanceFactor := goNoise.GetBias(1.0-distToCentre, 0.5)
		if centreDistanceFactor > 1.0 {
			centreDistanceFactor = 1.0
		}

		distToWater := easeOutQuad(float64(world.DistanceToWater[i]) / float64(maxDistR))
		height := heightMap[i]
		if height < fWaterHeight {
			tempMap[i] = goNoise.GetBias((centreDistanceFactor+height)*0.5, 0.8) * 0.8
		} else {
			noiseFactor := 0.2
			distWFactor := 0.3
			hFactor := 0.8

			temp := (centreDistanceFactor)

			temp = temp + ((1.0 - ((height) - fWaterHeight)) * hFactor) + (distToWater * distWFactor) + (tempMap[i] * noiseFactor)

			tempMap[i] = (temp / (1.0 + hFactor + distWFactor + noiseFactor))

		}

	}
	t = time.Now().UnixNano() / 1000000
	fmt.Println("temp calc msec: ", t-start)

	// Moisture
	for i := 0; i < len(moistMap); i++ {
		if heightMap[i] < fWaterHeight {
			moistMap[i] = 1.0
		} else {
			noise := moistMap[i]
			temp := 1.0 - tempMap[i]
			height := 1.0 - heightMap[i]
			waterdist := easeOutQuad(float64(world.DistanceToWater[i]) / float64(maxDistR))

			distToWaterFactor := 0.25
			noisefactor := 0.3
			tempFactor := 3.0
			heightFactor := 0.4
			moistMap[i] = ((waterdist * distToWaterFactor) + (noise * noisefactor) + (temp * tempFactor) + (height * heightFactor)) / (distToWaterFactor + noisefactor + tempFactor + heightFactor)
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
			if cTemp < 1.0 {
				biome = OCEANIC_GLACIER
			} else if height < (waterHeight - 35) {
				biome = OPEN_SEA
			} else {
				biome = SHALLOW_SEA
			}
		} else if height > world.MountainLine+20 {
			if moist < 90 {
				biome = MOUNTAIN
			} else {
				biome = GLACIER
			}
		} else if height > world.MountainLine {
			if moist < 125 {
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
				} else if precipitation < 75.0 {
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
