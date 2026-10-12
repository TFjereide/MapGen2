package rlconsole

import (
	common "common"
	rl "github.com/gen2brain/raylib-go/raylib"
	//"math"
	rand "math/rand"
	worldgen "worldgenerator"
)

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

	VectorMap      []rl.Vector2
	VectorStrength []float32

	RenderFG      bool
	RenderBG      bool
	RenderVectors bool

	RandColors []rl.Color
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
		Size:           size,
		X:              width,
		Y:              height,
		Font:           font,
		Symbol:         make([]byte, size),
		FgCol:          make([]rl.Color, size),
		BgCol:          make([]rl.Color, size),
		VectorMap:      make([]rl.Vector2, size),
		VectorStrength: make([]float32, size),
		RandWeight:     make([]float32, size),
		RandColors:     make([]rl.Color, 1000),
	}

	for i := 0; i < size; i++ {
		clayer.RandWeight[i] = 0.0 + (rand.Float32() * 0.18)
	}
	for i, _ := range clayer.RandColors {
		clayer.RandColors[i] = rl.NewColor(uint8(rand.Float32()*255.0), uint8(rand.Float32()*255.0), uint8(rand.Float32()*255.0), 255)
	}
	clayer.RandColors[0] = rl.Black

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
	if clayer.RenderVectors {
		for y := offsetY; y < (offsetY + sizeY); y++ {
			if y > clayer.Y {
				continue
			}

			for x := offsetX; x < (offsetX + sizeX); x++ {
				if x > clayer.X {
					continue
				}

				idx := getIdx(x, y, clayer.X)

				if clayer.VectorMap[idx].X == 0.0 && clayer.VectorMap[idx].Y == 0.0 {
					originPoint := rl.NewVector2(
						float32(x*cellWidth),
						float32(y*cellHeight))

					rl.DrawRectangleV(originPoint, rl.Vector2{float32(cellWidth), float32(cellHeight)}, rl.Green)
				}

				originPoint := rl.NewVector2(
					float32(x*cellWidth)+float32(cellWidth/2),
					float32(y*cellHeight)+float32(cellHeight/2))

				endPoint := rl.NewVector2(
					originPoint.X+(clayer.VectorMap[idx].X*3.0),
					originPoint.Y+(clayer.VectorMap[idx].Y*3.0))

				rl.DrawPixelV(originPoint, rl.Red)
				rl.DrawLineV(originPoint, endPoint, rl.SkyBlue)
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
	lightMap := make([]float32, world.X*world.Y)
	clayer.VectorMap = world.FlowDirection
	clayer.VectorStrength = world.FlowFactor

	//mixCol := rl.NewColor(20, 10, 25, 255)

	//determine lightness
	for i := 0; i < (world.X * world.Y); i++ {

		// height := float64(world.HeightMap[i] / 255.0)

		// if height < float64(world.WaterLine)/255.0 {
		// 	continue
		// }

		//x, y := posFromIdx(i, world.X)
		// leftHeight := 0.0
		//if x < 0 || y < 0 || y > world.Y || x > world.X {
		//	lightMap[i] = 0.0
		// 	continue
		// }
		// idx := getIdx(x, y, world.X)
		// leftHeight = float64(world.HeightMap[idx] / 255.0)

		// lightMap[i] = float32(height - leftHeight)
		//lightMap[i] *= 1.0

		//lightMap[i] = 1.0 - ((1.0 - float32(FAKE_NORMALS_2D[clayer.VectorMap[i]])) * float32(math.Pow(float64(world.HeightMap[i])/255, 1.2)))
		lightMap[i] = float32(FAKE_NORMALS_2D[clayer.VectorMap[i]])
	}

	for i := 0; i < (world.X * world.Y); i++ {

		//-----Determine Colors ----
		if world.HeightMap[i] > world.WaterLine {
			// Land
			waterHeight := float32(world.WaterLine) / 255.0
			fHeight := float32(world.HeightMap[i]) / 255
			waterRatioCorrected := (fHeight - waterHeight) * (1.0 / waterHeight)

			shadowFactor := waterRatioCorrected
			waterRatioCorrected *= 0.5

			brightness := (1.0 - lightMap[i]) * float32(shadowFactor*0.5)
			brightness = 1.0 - brightness
			brightness *= brightness

			//brightness := 1.0 - waterRatioCorrected // brighter further down
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
			waterRatioCorrected *= 0.45
			brightness := waterRatioCorrected + 0.55

			//heightCol := common.ByteToGrayColor(world.HeightMap[i])
			biomeColBg := worldgen.BIOME_BG_COLORS[world.BiomeMap[i]]
			biomeColFg := worldgen.BIOME_FG_COLORS[world.BiomeMap[i]]

			bgCol := common.ColorFloat32Multiply(biomeColBg, brightness)
			bgColormap[i] = common.ColorFloat32Multiply(bgCol, 1.0-clayer.RandWeight[i])

			fgCol := common.ColorFloat32Multiply(biomeColFg, brightness)
			fgColormap[i] = common.ColorFloat32Multiply(fgCol, 1.0+(clayer.RandWeight[i]*2.0))
		}

		//-----Determine Symbols ----
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

			//River tiles
			smallRiverThreshold := 0.02
			bigRiverThreshold := 0.25
			var bitMask byte = 0
			up := i - world.X
			down := i + world.X
			left := i - 1
			right := i + 1

			//		1
			//	2	0	4
			//		8
			if intInSlice(up, world.FlowsFrom[i]) || world.FlowsTo[i] == up { // Up
				if world.FlowStrength[up] > smallRiverThreshold {
					bitMask += 1
				}
			}
			if intInSlice(down, world.FlowsFrom[i]) || world.FlowsTo[i] == down { // Down
				if world.FlowStrength[down] > smallRiverThreshold {
					bitMask += 8
				}
			}
			if intInSlice(left, world.FlowsFrom[i]) || world.FlowsTo[i] == left { //Left
				if world.FlowStrength[left] > smallRiverThreshold {
					bitMask += 2
				}
			}
			if intInSlice(right, world.FlowsFrom[i]) || world.FlowsTo[i] == right { // Right
				if world.FlowStrength[right] > smallRiverThreshold {
					bitMask += 4
				}
			}
			if bitMask != 0 {

				if world.FlowStrength[i] > bigRiverThreshold {
					fgColormap[i] = worldgen.RIVER_FG_COLOR
					symMap[i] = IBM_437_ASCII_MAP[BitMaskThick[bitMask]]
				} else if world.FlowStrength[i] > smallRiverThreshold {
					fgColormap[i] = worldgen.RIVER_FG_COLOR
					symMap[i] = IBM_437_ASCII_MAP[BitMaskThin[bitMask]]
				}
			}

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
			fgColormap[i] = worldgen.COAST_FG_COLOR
			//symMap[i] = IBM_437_ASCII_MAP["#"]
		}
	}

	//TODO: Fix edge errors
	//clayer.BgCol = bgColormap
	//clayer.FgCol = fgColormap
	clayer.BgCol = common.BlurColorSlice(bgColormap, world.X, world.Y, 0.2)
	clayer.FgCol = common.BlurColorSlice(fgColormap, world.X, world.Y, 0.05)
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
func easeOutQuad(x float64) float64 {
	return 1.0 - (1.0-x)*(1.0-x)

}
func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
