package worldgenerator

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	N byte = iota
	NE
	E
	SE
	S
	SW
	W
	NW
)

var DIRECTION_VECTORS = map[rl.Vector2]byte{
	rl.Vector2{0, -1}:  N,
	rl.Vector2{1, -1}:  NE,
	rl.Vector2{1, 0}:   E,
	rl.Vector2{1, 1}:   SE,
	rl.Vector2{0, 1}:   S,
	rl.Vector2{-1, 1}:  SW,
	rl.Vector2{-1, 0}:  W,
	rl.Vector2{-1, -1}: NW,
}

//TODO: add glacial type
const (
	OPEN_SEA BiomeType = iota
	SHALLOW_SEA
	FRESHWATER
	OCEANIC_GLACIER
	GLACIER
	SNOW
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
	"snow",
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
	rl.NewColor(233, 235, 255, 255), // SNOW
	rl.NewColor(156, 155, 134, 255), // mountain
	rl.NewColor(100, 100, 90, 255),  // bare
	rl.NewColor(50, 82, 64, 255),    //Tundra
	rl.NewColor(185, 147, 100, 255), //temp desert
	rl.NewColor(220, 180, 132, 255), //subtropical desert
	rl.NewColor(16, 74, 41, 255),    // taiga
	rl.NewColor(30, 60, 34, 255),    // temp rainf
	rl.NewColor(28, 64, 5, 255),     // trop rainf
	rl.NewColor(86, 107, 55, 255),   // trop seas f
	rl.NewColor(42, 78, 38, 255),    // temp dec for
	rl.NewColor(50, 92, 44, 255),    // shrubland
	rl.NewColor(79, 92, 60, 255),    // grassland
	rl.NewColor(85, 100, 50, 255),   // savannah
	rl.NewColor(62, 63, 55, 255)}    // wetland

var BIOME_FG_COLORS = [...]rl.Color{
	rl.NewColor(31, 46, 143, 255),   // open w
	rl.NewColor(40, 72, 160, 255),   // shallow w
	rl.NewColor(46, 62, 165, 255),   // freshwater
	rl.NewColor(233, 235, 255, 255), // Oceanic Glacier
	rl.NewColor(243, 235, 255, 255), // Glacier
	rl.NewColor(243, 235, 255, 255), // SNOW
	rl.NewColor(160, 160, 160, 255), // Mountain
	rl.NewColor(60, 39, 20, 255),    // bare
	rl.NewColor(60, 60, 34, 255),    //Tundra
	rl.NewColor(185, 147, 100, 255), //temp desert
	rl.NewColor(220, 180, 132, 255), // subtropical desert
	rl.NewColor(20, 40, 50, 255),    // taiga
	rl.NewColor(15, 40, 30, 255),    // temp rainf
	rl.NewColor(17, 32, 7, 255),     // trop rainf
	rl.NewColor(47, 67, 22, 255),    // trop  seasonal f
	rl.NewColor(5, 40, 20, 255),     // temp dec for
	rl.NewColor(40, 65, 32, 255),    // shrubland
	rl.NewColor(79, 92, 60, 255),    // grassland
	rl.NewColor(90, 110, 57, 255),   // savannah
	rl.NewColor(62, 92, 47, 255)}    // wetland

var BIOME_SYMBOLS = [...]string{
	`≈`,                 // open w 226
	`≈`,                 // shallow w
	`≈≈`,                // freshwater
	`░░`,                // Oceanic Glacier
	`░░`,                // Glacier
	`░░`,                // SNOW
	`▲▲^`,               // Mountain
	`▲^^^^`,             // bare
	`░░░▒▓""`,           // Tundra
	`░░░▒▓,.,.`,         // temp desert
	`░░░▒▓.,.`,          // subtropical desert
	`♣¶♠♠♠♠♠↑î¥`,        // taiga
	`♣¶♠t↑îƒ§¥↕⌠`,       // temp rainf
	`♣¶♠t↑îƒ¥⌠`,         // trop rainf
	`♣¶♠t↑îƒ¥░░▒▓`,      // trop  seasonal f
	`♣¶♠t↑îƒ¥`,          // temp dec for
	`,.,.***""",."░░▒▓`, // shrubland
	`,.,.,.,.,.''"░░░░░░░░░░░▒▓♣♠♣♠♣`, // grassland
	`░░░░░░░░░▒▓♣♣♠"`,                 // savannah
	`≈░░¾≈`} // wetland

var RIVER_FG_COLOR = rl.NewColor(20, 140, 170, 150)
var COAST_FG_COLOR = rl.NewColor(20, 110, 140, 255)

// TODO: implement temp scaling to real world celsius
const MAX_TEMP_C = 40
const MIN_TEMP_c = 40 // in minus

//TODO: Moisture as humidity in real world scales?
const MAX_PRECIPITATION = 450
const MIN_PRECIPITATION = 0
