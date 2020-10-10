package worldgenerator

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Define terraintypes in the form of biomes
type BiomeType uint8
type Vector2Int32 struct {
	X, Y int32
}

func NewVector2Int32(x, y int32) Vector2Int32 {
	return Vector2Int32{X: x, Y: y}
}

// Container for all the data in a world
type World struct {
	// World Parameters
	X, Y         int
	WaterLine    byte
	MountainLine byte
	Equator      int

	// Maps/2D arrays of data
	HeightMap   []byte
	MoistureMap []byte
	TempMap     []byte
	SoilMap     []byte
	BiomeMap    []BiomeType

	DistanceToLand  []byte
	DistanceToWater []byte

	// Fake normals, should be renamed
	FlowDirection []rl.Vector2
	FlowFactor    []float32

	// Watershed data
	// indexes of tiles to and from which this tile flows
	FlowsTo      []int
	FlowsFrom    [][]int
	FlowStrength []float64
}
type WorldParams struct {
	Seed          int64
	Scale         float64
	Height, Width int
	WaterHeight   byte
	FWaterHeight  float64
	MountainLine  byte
	FMountainLine float64
}

func (world *World) GetBiomeAtPoint(x, y int) string {
	if x < 0 || y < 0 || x > world.X || y > world.Y {
		return "NONE"
	}
	return BIOME_NAMES[world.BiomeMap[getIdx(x, y, world.X)]]
}

func (world *World) GetHeightAtPoint(x, y int) byte {
	if x < 0 || y < 0 || x > world.X || y > world.Y {
		return 0
	}
	return world.HeightMap[getIdx(x, y, world.X)]
}

func (world *World) GetTemperatureAtPoint(x, y int) byte {
	if x < 0 || y < 0 || x > world.X || y > world.Y {
		return 0
	}
	return world.TempMap[getIdx(x, y, world.X)]
}

func (world *World) GetMoistureAtPoint(x, y int) byte {
	if x < 0 || y < 0 || x > world.X || y > world.Y {
		return 0
	}
	return world.MoistureMap[getIdx(x, y, world.X)]
}
func (world *World) GetFlowToAtPoint(x, y int) int {
	if x < 0 || y < 0 || x > world.X || y > world.Y {
		return 0
	}
	return world.FlowsTo[getIdx(x, y, world.X)]
}

func (world *World) GetFlowFromAtPoint(x, y int) []int {
	if x < 0 || y < 0 || x > world.X || y > world.Y {
		return []int{}
	}
	return world.FlowsFrom[getIdx(x, y, world.X)]
}

func clampF64(val, bottom, top float64) float64 {
	if val > top {
		return top
	} else if val < bottom {
		return bottom
	}
	return val
}

func colorMix(a, b rl.Color) rl.Color {

	col := rl.NewColor(0, 0, 0, 0)

	col.R = a.R + b.R
	col.G = a.G + b.G
	col.B = a.B + b.B
	col.A = a.A + b.A

	return col

}

func getIdx(x, y, width int) int {
	return ((y * width) + x)
}
func easeOutQuad(x float64) float64 {
	return 1.0 - (1.0-x)*(1.0-x)

}

func GetPoint(idx, width int) (int, int) {
	return int(idx % width), int(idx / width)
}
func GetPointI32(idx, width int) (int32, int32) {
	return int32(idx % width), int32(idx / width)
}

func blurByteSlice(bytes *[]byte, width, height, iterations int, holdout byte) *[]byte {

	readBuffer := *bytes
	writeBuffer := readBuffer

	for it := 0; it < iterations; it++ {

		for y := 1; y < height-1; y++ {
			for x := 1; x < width-1; x++ {

				idx := getIdx(x, y, width)
				currByte := readBuffer[idx]
				if currByte == holdout {
					writeBuffer[idx] = currByte
					continue
				}

				for dy := -1; dy < 2; dy++ {
					for dx := -1; dx < 2; dx++ {
						if dx != 0 && dy != 0 {

							if x+dx < 0 || x+dx >= width || y+dy < 0 || y+dy >= height {
								continue
							}

							_idx := getIdx(x+dx, y+dy, width)
							mixByte := readBuffer[_idx]
							currByte = byte((int64(mixByte) + int64(currByte)) / 2)
						}
					}
				}
				writeBuffer[idx] = currByte
			}
		}
		readBuffer = writeBuffer
	}

	return &readBuffer
}

// Finds distance to water and distance to land for each tile
func findShoreDistances(heightMap []float64, world *World, maxDistR int) {

	fWaterHeight := float64(world.WaterLine) / 255
	width := world.X
	for idx := range heightMap {
		x, y := GetPoint(idx, width)
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

}
func cubicInterpolate(y0, y1, y2, y3, mu float64) float64 {
	var a0, a1, a2, a3 float64
	mu2 := mu * mu

	a0 = y3 - y2 - y0 + y1
	a1 = y0 - y1 - a0
	a2 = y2 - y0
	a3 = y1

	return clampF64(a0*mu*mu2+a1*mu2+a2*mu+a3, 0.0, 1.0)
}

// Small jobs helper functions
// func getDistanceRow(buffer *sync.WaitGroup, wdist, ldist chan []byte, y, width, maxDistR int, heightMap []float64, fWaterHeight float64) {

// 	defer buffer.Done()

// 	ldistrow := make([]byte, width+1)
// 	wdistrow := make([]byte, width+1)
// 	row[0] = byte(y)
// 	for i := 1; i < width+1; i++ {

// 	distanceFound:
// 		// Check within a radius of r
// 		for r := 1; r < maxDistR; r++ {

// 			for dy := -r; dy <= r; dy++ {
// 				for dx := -r; dx <= r; dx++ {

// 					if dx < r && dy < r && dx > -r && dy > -r {
// 						continue
// 					}

// 					didx := getIdx((x + dx), (y + dy), width)
// 					if didx < 0 || didx > len(heightMap)-1 {
// 						continue
// 					}

// 					locElev := heightMap[didx]

// 					if water {
// 						// Check  found land
// 						if locElev > fWaterHeight {
// 							row[i] = byte(r)
// 							break distanceFound
// 						}

// 					} else {
// 						//Check found water
// 						if locElev < fWaterHeight {
// 							row[i] = byte(r)
// 							break distanceFound
// 						}
// 					}

// 				}
// 			}
// 		}

// 	}

// 	c <- row
// }
