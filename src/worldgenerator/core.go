package worldgenerator

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Define terraintypes in the form of biomes
type BiomeType uint8

//TODO: add glacial type
const (
	OPEN_SEA BiomeType = iota
	SHALLOW_SEA
	FRESHWATER
	OCEANIC_GLACIER
	GLACIER
	MOUNTAIN
	BARE
	TUNDRA
	TEMPERATE_DESERT
	SUBTROPICAL_DESERT
	TAIGA
	TEMPERATE_RAINFOREST
	TROPICAL_RAINFOREST
	TROPICAL_SEASONAL_FOREST
	TEMPERATE_DECIDUOUS_FOREST
	SHRUBLAND
	GRASSLAND
	SAVANNA
	WETLAND
)

var BIOME_NAMES = [...]string{
	"open sea",
	"shallow sea",
	"freshwater",
	"oceanic glacier",
	"glacier",
	"mountain",
	"bare",
	"tundra",
	"temperate desert",
	"subropical desert",
	"taiga",
	"temperate rainforest",
	"tropical rainforest",
	"tropical seasonal forest",
	"temperate deciduous forest",
	"shrubland",
	"grassland",
	"savannah",
	"wetland"}

var BIOME_BG_COLORS = [...]rl.Color{
	rl.NewColor(17, 14, 65, 255),    // open w
	rl.NewColor(8, 34, 73, 255),     //  shallow w
	rl.NewColor(8, 24, 73, 255),     //freshwater
	rl.NewColor(233, 235, 255, 255), // oceanic Glacier
	rl.NewColor(233, 235, 255, 255), // Glacier
	rl.NewColor(156, 155, 134, 255), // mountain
	rl.NewColor(139, 129, 118, 255), // bare
	rl.NewColor(50, 82, 64, 255),    //Tundra
	rl.NewColor(189, 162, 132, 255), //temp desert
	rl.NewColor(189, 162, 132, 255), //subtropical desert
	rl.NewColor(16, 74, 41, 255),    // taiga
	rl.NewColor(30, 60, 34, 255),    // temp rainf
	rl.NewColor(28, 64, 5, 255),     // trop rainf
	rl.NewColor(86, 107, 55, 255),   // trop seas f
	rl.NewColor(42, 78, 38, 255),    // temp dec for
	rl.NewColor(50, 92, 44, 255),    // shrubland
	rl.NewColor(62, 92, 47, 255),    // grassland
	rl.NewColor(85, 100, 50, 255),   // savannah
	rl.NewColor(62, 63, 55, 255)}    // wetland

var BIOME_FG_COLORS = [...]rl.Color{
	rl.NewColor(31, 46, 143, 255),   // open w
	rl.NewColor(40, 72, 160, 255),   // shallow w
	rl.NewColor(46, 62, 165, 255),   // freshwater
	rl.NewColor(233, 235, 255, 255), // Oceanic Glacier
	rl.NewColor(233, 235, 255, 255), // Glacier
	rl.NewColor(156, 155, 134, 255), // Mountain
	rl.NewColor(115, 129, 118, 255), // bare
	rl.NewColor(50, 82, 64, 255),    //Tundra
	rl.NewColor(189, 162, 132, 255), //temp desert
	rl.NewColor(189, 162, 132, 255), // subtropical desert
	rl.NewColor(20, 60, 50, 255),    // taiga
	rl.NewColor(15, 40, 30, 255),    // temp rainf
	rl.NewColor(17, 32, 7, 255),     // trop rainf
	rl.NewColor(47, 67, 22, 255),    // trop  seasonal f
	rl.NewColor(5, 40, 20, 255),     // temp dec for
	rl.NewColor(50, 92, 44, 255),    // shrubland
	rl.NewColor(62, 92, 47, 255),    // grassland
	rl.NewColor(90, 110, 57, 255),   // savannah
	rl.NewColor(62, 92, 47, 255)}    // wetland

var BIOME_SYMBOLS = [...]string{
	`≈`,               // open w 226
	`≈`,               // shallow w
	`≈≈`,              // freshwater
	`░░`,              // Oceanic Glacier
	`░░`,              // Glacier
	`▲∩∩`,             // Mountain
	`∩`,               // bare
	`░░░▒▓""`,         // Tundra
	`░≈`,              // temp desert
	`░≈`,              // subtropical desert
	`♣¶♠♠♠♠♠t↑îƒ§¥↕⌠`, // taiga
	`♣¶♠t↑îƒ§¥↕⌠`,     // temp rainf
	`♣¶♠t↑îƒ§¥⌠`,      // trop rainf
	`♣¶♠t↑îƒ§¥░░▒▓`,   // trop  seasonal f
	`♣¶♠t↑îƒ§¥`,       // temp dec for
	`,.,.""",."░░▒▓`,  // shrubland
	`,.,.,.,.,.''"░░▒▓♣♠♣♠♣`, // grassland
	`░░░░░░░░░▒▓♣♣♠"`,        // savannah
	`≈░░¾≈`} // wetland

// TODO: implement temp scaling to real world celsius
const MAX_TEMP_C = 38
const MIN_TEMP_c = 18 // in minus

//TODO: Moisture as humidity in real world scales?
const MAX_PRECIPITATION = 450
const MIN_PRECIPITATION = 0

// Container for all the data in a world
type World struct {
	X, Y            int
	WaterLine       byte
	MountainLine    byte
	Equator         int
	HeightMap       []byte
	MoistureMap     []byte
	TempMap         []byte
	SoilMap         []byte
	BiomeMap        []BiomeType
	DistanceToLand  []byte
	DistanceToWater []byte
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
	return BIOME_NAMES[world.BiomeMap[getIdx(x, y, world.X)]]
}
func (world *World) GetHeightAtPoint(x, y int) byte {
	return world.HeightMap[getIdx(x, y, world.X)]
}
func (world *World) GetTemperatureAtPoint(x, y int) byte {
	return world.TempMap[getIdx(x, y, world.X)]
}
func (world *World) GetMoistureAtPoint(x, y int) byte {
	return world.MoistureMap[getIdx(x, y, world.X)]
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

func getPoint(idx, width int) (int, int) {
	return int(idx % width), int(idx / width)
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
