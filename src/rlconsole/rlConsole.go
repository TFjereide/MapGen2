package rlconsole

import (
	common "common"
	rl "github.com/gen2brain/raylib-go/raylib"
	rand "math/rand"
	worldgen "worldgenerator"
)

var IBM_437_ASCII_MAP = map[string]byte{
	"☺":    1,
	"☻":    2,
	"♥":    3,
	"♦":    4,
	"♣":    5,
	"♠":    6,
	"•":    7,
	"◘":    8,
	"○":    9,
	"◙":    10,
	"♂":    11,
	"♀":    12,
	"♪":    13,
	"♫":    14,
	"☼":    15,
	"►":    16,
	"◄":    17,
	"↕":    18,
	"‼":    19,
	"¶":    20,
	"§":    21,
	"▬":    22,
	"↨":    23,
	"↑":    24,
	"↓":    25,
	"→":    26,
	"←":    27,
	"∟":    28,
	"↔":    29,
	"▲":    30,
	"▼":    31,
	"SP":   32,
	"!":    33,
	`"`:    34,
	"#":    35,
	"$":    36,
	"%":    37,
	"&":    38,
	"'":    39,
	"(":    40,
	")":    41,
	"*":    42,
	"+":    43,
	",":    44,
	"-":    45,
	".":    46,
	"/":    47,
	"0":    48,
	"1":    49,
	"2":    50,
	"3":    51,
	"4":    52,
	"5":    53,
	"6":    54,
	"7":    55,
	"8":    56,
	"9":    57,
	":":    58,
	";":    59,
	"<":    60,
	"=":    61,
	">":    62,
	"?":    63,
	"@":    64,
	"A":    65,
	"B":    66,
	"C":    67,
	"D":    68,
	"E":    69,
	"F":    70,
	"G":    71,
	"H":    72,
	"I":    73,
	"J":    74,
	"K":    75,
	"L":    76,
	"M":    77,
	"N":    78,
	"O":    79,
	"P":    80,
	"Q":    81,
	"R":    82,
	"S":    83,
	"T":    84,
	"U":    85,
	"V":    86,
	"W":    87,
	"X":    88,
	"Y":    89,
	"Z":    90,
	"[":    91,
	"\\":   92,
	"]":    93,
	"^":    94,
	"_":    95,
	"`":    96,
	"a":    97,
	"b":    98,
	"c":    99,
	"d":    101,
	"e":    102,
	"f":    103,
	"g":    104,
	"h":    105,
	"i":    106,
	"j":    107,
	"k":    108,
	"l":    109,
	"m":    110,
	"n":    111,
	"o":    112,
	"p":    113,
	"q":    114,
	"r":    115,
	"s":    116,
	"t":    117,
	"u":    118,
	"v":    119,
	"w":    120,
	"x":    121,
	"y":    122,
	"z":    123,
	"{":    124,
	"|":    125,
	"}":    126,
	"~":    127,
	"Ç":    128,
	"ü":    129,
	"é":    130,
	"â":    131,
	"ä":    132,
	"à":    133,
	"å":    134,
	"ç":    135,
	"ê":    136,
	"ë":    137,
	"è":    138,
	"ï":    139,
	"î":    140,
	"ì":    141,
	"Ä":    142,
	"Å":    143,
	"É":    144,
	"æ":    145,
	"Æ":    146,
	"ô":    147,
	"ö":    148,
	"ò":    149,
	"û":    150,
	"ù":    151,
	"ÿ":    152,
	"Ö":    153,
	"Ü":    154,
	"ø":    155,
	"£":    156,
	"Ø":    157,
	"×":    158,
	"ƒ":    159,
	"á":    160,
	"í":    161,
	"ó":    162,
	"ú":    163,
	"ñ":    164,
	"Ñ":    165,
	"ª":    166,
	"º":    167,
	"¿":    168,
	"®":    169,
	"¬":    170,
	"½":    171,
	"¼":    172,
	"¡":    173,
	"«":    174,
	"»":    175,
	"░":    176,
	"▒":    177,
	"▓":    178,
	"│":    179,
	"┤":    180,
	"Á":    181,
	"Â":    182,
	"À":    183,
	"©":    184,
	"╣":    185,
	"║":    186,
	"╗":    187,
	"╝":    188,
	"¢":    189,
	"¥":    190,
	"┐":    191,
	"└":    192,
	"┴":    193,
	"┬":    194,
	"├":    195,
	"─":    196,
	"┼":    197,
	"ã":    198,
	"Ã":    199,
	"╚":    200,
	"╔":    201,
	"╩":    202,
	"╦":    203,
	"╠":    204,
	"═":    205,
	"╬":    206,
	"¤":    207,
	"ð":    208,
	"Ð":    209,
	"Ê":    210,
	"Ë":    211,
	"È":    212,
	"ı":    213,
	"Í":    214,
	"Î":    215,
	"Ï":    216,
	"┘":    217,
	"┌":    218,
	"█":    219,
	"▄":    220,
	"¦":    221,
	"Ì":    222,
	"▀":    223,
	"α":    224,
	"ß":    225,
	"Γ":    226,
	"π":    227,
	"Σ":    228,
	"σ":    229,
	"µ":    230,
	"τ":    231,
	"Φ":    232,
	"Θ":    233,
	"Ω":    234,
	"δ":    235,
	"∞":    236,
	"φ":    237,
	"ε":    238,
	"∩":    239,
	"≡":    240,
	"±":    241,
	"≥":    242,
	"≤":    243,
	"⌠":    244,
	"⌡":    245,
	"÷":    246,
	"≈":    247,
	"°":    248,
	"∙":    249,
	"·":    250,
	"√":    251,
	"ⁿ":    252,
	"²":    253,
	"■":    254,
	"NBSP": 255}

var BitMaskThick = [...]string{
	"≈", //0
	"≈", //1
	"≈", //2
	"╝", //3
	"≈", //4
	"╚", //5
	"═", //6
	"╩", //7
	"≈", //8
	"║", //9
	"╗", //10
	"╣", //11
	"╔", //12
	"╠", //13
	"╦", //14
	"╬"} //15

var BitMaskThin = [...]string{
	"≈", //0
	"≈", //1
	"≈", //2
	"┘", //3
	"≈", //4
	"└", //5
	"─", //6
	"┴", //7
	"≈", //8
	"│", //9
	"┐", //10
	"┤", //11
	"┌", //12
	"├", //13
	"┬", //14
	"┼"} //15

type SpriteSheetFont struct {
	Texture    rl.Texture2D
	Width      int32
	Height     int32
	CharWidth  int32
	CharHeight int32
}

type ConsoleLayer struct {
	Size       int
	X, Y       int
	Texture    rl.Texture2D
	Font       SpriteSheetFont
	Symbol     []byte
	FgCol      []rl.Color
	BgCol      []rl.Color
	RandWeight []float32

	RenderFG bool
	RenderBG bool
}

func NewSpriteSheetFont(tex rl.Texture2D, width, height int32) SpriteSheetFont {
	return SpriteSheetFont{
		Texture:    tex,
		Width:      tex.Width,
		Height:     tex.Height,
		CharWidth:  width,
		CharHeight: height}
}

// Returns a new *ConsoleLayer, which is a primary surface to display graphics on.
func NewConsoleLayer(width, height int, font SpriteSheetFont) *ConsoleLayer {
	size := width * height
	clayer := ConsoleLayer{
		Size:       size,
		X:          width,
		Y:          height,
		Font:       font,
		Symbol:     make([]byte, size),
		FgCol:      make([]rl.Color, size),
		BgCol:      make([]rl.Color, size),
		RandWeight: make([]float32, size)}

	for i := 0; i < size; i++ {
		clayer.RandWeight[i] = 0.0 + (rand.Float32() * 0.22)
	}
	return &clayer
}

// Unloads every texture from gpu memory
func (clayer *ConsoleLayer) Unload() {
	rl.UnloadTexture(clayer.Texture)
	rl.UnloadTexture(clayer.Font.Texture)
}

// Sets Symbol, fgcolor and bgcolor of coordinate
// Does not update texture - used in batch calls
func (clayer *ConsoleLayer) SetCell(x, y int, sym string, fg, bg rl.Color) {
	idx := getIdx(x, y, clayer.X)
	clayer.Symbol[idx] = []byte(sym)[0]
	clayer.FgCol[idx] = fg
	clayer.BgCol[idx] = bg
}

// Loads in bg color data from image struct
func (clayer *ConsoleLayer) SetBG(pixels []rl.Color) {
	if len(pixels) == clayer.Size {
		clayer.BgCol = pixels
	} else {
		panic("Unreasonable pixel slice")
	}
}

// Updates bgtexture from current data
func (clayer *ConsoleLayer) UpdateTexture() {
	rl.UnloadTexture(clayer.Texture)
	img := rl.LoadImageEx(clayer.BgCol, int32(clayer.X), int32(clayer.Y))
	clayer.Texture = rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
}

// Draws ConsoleLayer to screen
func (clayer *ConsoleLayer) Draw(offsetX, offsetY, sizeX, sizeY int) {

	// The bg colors
	if clayer.RenderBG {
		rl.DrawTextureEx(clayer.Texture, rl.NewVector2(0, 0), 0, float32(clayer.Font.CharWidth), rl.White)
	}

	cellWidth := int(clayer.Font.CharWidth)
	cellHeight := int(clayer.Font.CharHeight)

	// Symbols currently drawn all individually
	if clayer.RenderFG {
		for y := offsetY; y < (offsetY + sizeY); y++ {
			if y > clayer.Y {
				continue
			}
			for x := offsetX; x < (offsetX + sizeX); x++ {
				if x > clayer.X {
					continue
				}

				idx := getIdx(x, y, clayer.X)
				fgc := clayer.FgCol[idx]
				sym := clayer.Symbol[idx]
				pos := rl.NewVector2(float32(x*cellWidth), float32(y*cellHeight))

				symRect := byteToSymbolRect(
					sym,
					int(clayer.Font.Width),
					int(clayer.Font.Height),
					int(clayer.Font.CharWidth),
					int(clayer.Font.CharHeight))

				rl.DrawTextureRec(clayer.Font.Texture, symRect, pos, fgc)
				//rl.DrawText(string(sym), int32(pos.X), int32(pos.Y), 8, fgc)
			}
		}
	}
}

// Makes world data graphical, by importing data and setting the console
func (clayer *ConsoleLayer) ExtractVisuals(world *worldgen.World) {

	// Convert world info to bg and fg color
	bgColormap := make([]rl.Color, world.X*world.Y)
	fgColormap := make([]rl.Color, world.X*world.Y)
	symMap := make([]byte, world.X*world.Y)

	//mixCol := rl.NewColor(20, 10, 25, 255)

	for i := 0; i < (world.X * world.Y); i++ {

		if world.HeightMap[i] > world.WaterLine {
			// Land
			waterHeight := float32(world.WaterLine) / 255.0
			fHeight := float32(world.HeightMap[i]) / 255
			waterRatioCorrected := (fHeight - waterHeight) * (1.0 / waterHeight)
			waterRatioCorrected *= 0.55
			brightness := 1.0 - waterRatioCorrected // brighter further down
			//brightness := (fHeight * 0.5) + 0.5 // brighter further up

			biomeColBg := worldgen.BIOME_BG_COLORS[world.BiomeMap[i]]
			biomeColFg := worldgen.BIOME_FG_COLORS[world.BiomeMap[i]]

			bgCol := common.ColorFloat32Multiply(biomeColBg, brightness)
			bgColormap[i] = common.ColorFloat32Multiply(bgCol, 1.0-clayer.RandWeight[i])

			//fgcolnew := common.ColorMixWeighted(mixCol, biomeColFg, clayer.RandWeight[i])
			fgCol := common.ColorFloat32Multiply(biomeColFg, brightness)
			fgColormap[i] = common.ColorFloat32Multiply(fgCol, 1.0+clayer.RandWeight[i])

		} else {
			// Water
			waterHeight := float32(world.WaterLine) / 255.0
			fHeight := float32(world.HeightMap[i]) / 255
			waterRatioCorrected := fHeight * (1.0 / waterHeight)
			waterRatioCorrected *= 0.5
			brightness := waterRatioCorrected + 0.5

			//heightCol := common.ByteToGrayColor(world.HeightMap[i])
			biomeColBg := worldgen.BIOME_BG_COLORS[world.BiomeMap[i]]
			biomeColFg := worldgen.BIOME_FG_COLORS[world.BiomeMap[i]]

			bgCol := common.ColorFloat32Multiply(biomeColBg, brightness)
			bgColormap[i] = common.ColorFloat32Multiply(bgCol, 1.0-clayer.RandWeight[i])

			fgCol := common.ColorFloat32Multiply(biomeColFg, brightness)
			fgColormap[i] = common.ColorFloat32Multiply(fgCol, 1.0+(clayer.RandWeight[i]*1.5))
		}

		if world.DistanceToLand[i] != 1 {
			str := worldgen.BIOME_SYMBOLS[world.BiomeMap[i]]
			runes := []rune(str)
			var runechar rune
			if len(runes) > 1 {
				runechar = runes[rand.Intn(len(runes))]
			} else {
				runechar = runes[0]
			}

			asciibyte := IBM_437_ASCII_MAP[string(runechar)]

			symMap[i] = asciibyte

		} else {
			var bitMask byte = 0

			if (i-world.X) > 0 && world.DistanceToLand[i-world.X] == 1 { // Up
				bitMask += 1
			}
			if (i+world.X) < world.X*world.Y && world.DistanceToLand[i+world.X] == 1 { // Down
				bitMask += 8
			}
			if (i-1) > 0 && world.DistanceToLand[i-1] == 1 { //Left
				bitMask += 2
			}
			if (i+1) < world.X*world.Y && world.DistanceToLand[i+1] == 1 { // Right
				bitMask += 4
			}

			symMap[i] = IBM_437_ASCII_MAP[BitMaskThick[bitMask]]
			//symMap[i] = IBM_437_ASCII_MAP["#"]
		}
	}

	//TODO: Fix edge errors
	clayer.BgCol = common.BlurColorSlice(bgColormap, world.X, world.Y, 0.2)
	clayer.FgCol = common.BlurColorSlice(fgColormap, world.X, world.Y, 0.15)
	clayer.Symbol = symMap
	clayer.UpdateTexture()

}

// Helper - to get slice index easily
func getIdx(x, y, width int) int {
	return ((y * width) + x)
}

func posFromIdx(i, width int) (int, int) {
	return int(i % width), int(i / width)
}

// Helper - to get right rect settings to display a character from a SpriteSheetFont
func getSymbolRect(symbol string, texWidth, texHeight, fSizeX, fSizeY int) rl.Rectangle {

	characterlist := []byte(symbol)
	character := byte(characterlist[0])

	width := byte(texWidth / fSizeX)

	xPos := int(character % width)
	yPos := int(character / width)

	return rl.Rectangle{
		X:      float32(xPos * fSizeX),
		Y:      float32(yPos * fSizeY),
		Width:  float32(fSizeX),
		Height: float32(fSizeY)}
}
func byteToSymbolRect(symbol byte, texWidth, texHeight, fSizeX, fSizeY int) rl.Rectangle {

	character := symbol

	width := byte(texWidth / fSizeX)

	xPos := int(character % width)
	yPos := int(character / width)

	return rl.Rectangle{
		X:      float32(xPos * fSizeX),
		Y:      float32(yPos * fSizeY),
		Width:  float32(fSizeX),
		Height: float32(fSizeY)}
}

// Positive modulo, returns non negative solution to x % d
func pmod(x, b int) int {
	x = x % b
	if x >= 0 {
		return x
	}
	if b < 0 {
		return x - b
	}
	return x + b
}
