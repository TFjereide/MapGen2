package main

import (
	"common"
	"fmt"
	rayGui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	//imgui "github.com/inkyblackness/imgui-go"
	"goNoise"

	rlConsole "rlConsole"
	worldgen "worldgenerator"
)

type GameContainer struct {
	Console *rlConsole.ConsoleLayer
	World   *worldgen.World

	Camera      rl.Camera2D
	CameraSizeX int32
	CameraSizeY int32
}

//  Actions performed through normal inputs
type InputContext struct {
	MouseX, MouseY         int32
	MoveLeft               bool
	MoveRight              bool
	MoveUp                 bool
	MoveDown               bool
	WorldSizeIncrease      bool
	WorldSizeDecrease      bool
	SetModeColor           bool
	SetModeHeight          bool
	SetModeTemperature     bool
	SetModeMoisture        bool
	SetModeDistanceToWater bool
	SetModeFlow            bool
}

func newInputContext() InputContext {
	return InputContext{
		MoveLeft:               false,
		MoveRight:              false,
		MoveUp:                 false,
		MoveDown:               false,
		SetModeColor:           false,
		SetModeHeight:          false,
		SetModeTemperature:     false,
		SetModeMoisture:        false,
		SetModeDistanceToWater: false,
		SetModeFlow:            false}
}

// Other stuff the Gui needs to draw
type GuiData struct {
	MouseX, MouseY int32
}

// Changes returned by RayGui's lil' widgets
type GuiResponseBuffer struct {
	Changed      bool
	Seed         *int64
	WorldScale   *float64
	NoiseTypes   []*goNoise.NoiseType
	NoiseWeights []*float32
	NoiseBias    []*float64
	NoiseGain    []*float64
}

const (
	COLOR int = iota
	HEIGHT
	TEMPERATURE
	MOISTURE
	DISTANCE_TO_WATER
	FLOW
)

var MODE_NAMES = [...]string{
	"color",
	"height",
	"temperature",
	"moisture",
	"distance to water",
	"flow"}

var MODE int = COLOR

const mapWidth = 512
const mapHeight = 256
const fontWidth = 8
const fontHeight = 8

func moveCamera(game *GameContainer, input InputContext) {

	moveDirX := common.BToI(input.MoveRight) - common.BToI(input.MoveLeft)
	moveDirY := common.BToI(input.MoveDown) - common.BToI(input.MoveUp)

	game.Camera.Target.X += float32(moveDirX) * 20.0
	game.Camera.Target.Y += float32(moveDirY) * 20.0

	//Clamp
	if game.Camera.Target.X < 0 {
		game.Camera.Target.X = 0
	} else if int32(game.Camera.Target.X) > mapWidth*fontWidth-game.CameraSizeX {
		game.Camera.Target.X = float32(mapWidth*fontWidth - game.CameraSizeX)
	}
	if game.Camera.Target.Y < 0 {
		game.Camera.Target.Y = 0
	} else if int32(game.Camera.Target.Y) > mapHeight*fontHeight-game.CameraSizeY {
		game.Camera.Target.Y = float32(mapHeight*fontHeight - game.CameraSizeY)
	}

}
func setWorldScale(sizeX, sizeY float64, input InputContext) (float64, float64) {
	if input.WorldSizeIncrease {
		sizeX += 0.1
		sizeY += 0.1
	} else if input.WorldSizeDecrease {
		sizeX -= 0.1
		sizeY -= 0.1
	}
	return sizeX, sizeY
}
func setConsoleMode(input InputContext) {
	if input.SetModeColor {
		MODE = COLOR
	} else if input.SetModeHeight {
		MODE = HEIGHT
	} else if input.SetModeTemperature {
		MODE = TEMPERATURE
	} else if input.SetModeMoisture {
		MODE = MOISTURE
	} else if input.SetModeDistanceToWater {
		MODE = DISTANCE_TO_WATER
	} else if input.SetModeFlow {
		MODE = FLOW
	}
}

func updateConsoleVisuals(game *GameContainer) {
	switch MODE {
	case COLOR:
		game.Console.RenderFG = true
		game.Console.RenderBG = true
		game.Console.RenderVectors = false
		game.Console.ExtractVisuals(game.World)
	case TEMPERATURE:
		game.Console.RenderFG = false
		game.Console.RenderBG = true
		game.Console.RenderVectors = false
		game.Console.SetBG(common.ConvertArrayByteToRLColor(game.World.TempMap))
		game.Console.UpdateTexture()
	case HEIGHT:
		game.Console.RenderFG = false
		game.Console.RenderBG = true
		game.Console.RenderVectors = false
		game.Console.SetBG(common.ConvertArrayByteToRLColor(game.World.HeightMap))
		game.Console.UpdateTexture()
	case MOISTURE:
		game.Console.RenderFG = false
		game.Console.RenderBG = true
		game.Console.RenderVectors = false
		game.Console.SetBG(common.ConvertArrayByteToRLColor(game.World.MoistureMap))
		game.Console.UpdateTexture()

	case DISTANCE_TO_WATER:
		game.Console.RenderFG = false
		game.Console.RenderBG = true
		game.Console.RenderVectors = false
		game.Console.SetBG(common.ConvertArrayByteToRLColor(game.World.DistanceToWater))
		game.Console.UpdateTexture()

	case FLOW:
		game.Console.RenderFG = false
		game.Console.RenderBG = true
		game.Console.RenderVectors = true
		game.Console.SetBG(common.ConvertArrayByteToRLColor(game.World.HeightMap))
		game.Console.VectorMap = game.World.FlowDirection
		game.Console.UpdateTexture()
	}
}

func getInput(inputContext *InputContext) (bool, bool) {

	changed := false
	regen := false
	if rl.IsKeyPressed(rl.KeyOne) {
		changed = true
		inputContext.SetModeColor = true
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		changed = true
		inputContext.SetModeHeight = true
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		changed = true
		inputContext.SetModeTemperature = true
	}
	if rl.IsKeyPressed(rl.KeyFour) {
		changed = true
		inputContext.SetModeMoisture = true
	}
	if rl.IsKeyPressed(rl.KeyFive) {
		changed = true
		inputContext.SetModeDistanceToWater = true
	}
	if rl.IsKeyPressed(rl.KeySix) {
		changed = true
		inputContext.SetModeFlow = true
	}
	if rl.IsKeyDown(rl.KeyW) {
		inputContext.MoveUp = true
	}
	if rl.IsKeyDown(rl.KeyS) {
		inputContext.MoveDown = true
	}
	if rl.IsKeyDown(rl.KeyA) {
		inputContext.MoveLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) {
		inputContext.MoveRight = true
	}
	if rl.IsKeyDown(rl.KeyP) {
		changed = true
		regen = true
		inputContext.WorldSizeIncrease = true
	}
	if rl.IsKeyDown(rl.KeyO) {
		changed = true
		regen = true
		inputContext.WorldSizeDecrease = true
	}

	inputContext.MouseX = rl.GetMouseX()
	inputContext.MouseY = rl.GetMouseY()

	return changed, regen
}

func drawGui(data GuiData, buffer GuiResponseBuffer) GuiResponseBuffer {

	guiOffsetX := float32(8)
	guiOffsetY := float32(8)
	offsetIncrement := float32(18)
	offsetIncrementX := float32(4)
	changed := false

	// fps widget
	fps := rl.GetFPS()
	strFPS := fmt.Sprintf("%.2f", fps)

	rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, 0, 0),
		"FPS: "+strFPS,
		rl.RayWhite, rl.Gray, rl.Gray)

	guiOffsetY += offsetIncrement

	// mode widget
	mode_name := MODE_NAMES[MODE]
	rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, 0, 0),
		"View mode: "+mode_name,
		rl.RayWhite, rl.Gray, rl.Gray)
	guiOffsetY += offsetIncrement

	// Seed Widget
	rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, 45, 0),
		"Seed:", rl.RayWhite, rl.Gray, rl.Gray)

	guiOffsetX += 45
	oldSeed := *buffer.Seed
	*buffer.Seed = int64(rayGui.Spinner(rl.NewRectangle(guiOffsetX, guiOffsetY, 45, offsetIncrement), int(*buffer.Seed), 1, 300))
	if *buffer.Seed != oldSeed {
		changed = true
	}
	guiOffsetX -= 45
	guiOffsetY += offsetIncrement

	// TODO: World scale corruption bug
	// WorldscaleWidget
	rect := rl.NewRectangle(guiOffsetX, guiOffsetY, 60, 18)
	rayGui.LabelEx(rect, "WorldScale", rl.RayWhite, rl.Gray, rl.Gray)
	rect = rl.NewRectangle(guiOffsetX+rect.Width, guiOffsetY, 100, 18)

	oldScale := *buffer.WorldScale
	*buffer.WorldScale = float64(rayGui.SliderBar(rect, float32(*buffer.WorldScale), 0.0, 2.0))
	if *buffer.WorldScale != oldScale {
		changed = true
	}

	rayGui.Label(rect, fmt.Sprintf("%.3f", *buffer.WorldScale))
	guiOffsetY += offsetIncrement

	// NoiseWeights
	rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, 100, 18), "Noise Layer Weights:", rl.RayWhite, rl.Gray, rl.Gray)
	guiOffsetY += offsetIncrement
	guiOffsetX += offsetIncrementX
	for i := range buffer.NoiseWeights {
		labelSizeX := float32(65)
		rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, labelSizeX, 18), goNoise.NOISETYPE_STRING[*buffer.NoiseTypes[i]], rl.RayWhite, rl.Gray, rl.Gray)

		rect = rl.NewRectangle(guiOffsetX+labelSizeX, guiOffsetY, 100, 18)
		oldWeight := *buffer.NoiseWeights[i]
		*buffer.NoiseWeights[i] = rayGui.SliderBar(rect, *buffer.NoiseWeights[i], 0.0, 1.0)
		if *buffer.NoiseWeights[i] != oldWeight {
			changed = true
		}

		rayGui.Label(rect, fmt.Sprintf("%.3f", *buffer.NoiseWeights[i]))
		guiOffsetY += offsetIncrement
	}
	guiOffsetX -= offsetIncrementX

	// NoiseBias
	rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, 100, 18), "Noise Layer Bias:", rl.RayWhite, rl.Gray, rl.Gray)
	guiOffsetY += offsetIncrement
	guiOffsetX += offsetIncrementX
	for i := range buffer.NoiseBias {
		labelSizeX := float32(65)
		rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, labelSizeX, 18), goNoise.NOISETYPE_STRING[*buffer.NoiseTypes[i]], rl.RayWhite, rl.Gray, rl.Gray)

		rect = rl.NewRectangle(guiOffsetX+labelSizeX, guiOffsetY, 100, 18)
		oldWeight := *buffer.NoiseBias[i]
		*buffer.NoiseBias[i] = float64(rayGui.SliderBar(rect, float32(oldWeight), 0.0, 1.0))
		if *buffer.NoiseBias[i] != oldWeight {
			changed = true
		}

		rayGui.Label(rect, fmt.Sprintf("%.3f", *buffer.NoiseBias[i]))
		guiOffsetY += offsetIncrement
	}
	guiOffsetX -= offsetIncrementX

	// NoiseGain
	rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, 100, 18), "Noise Layer Gain:", rl.RayWhite, rl.Gray, rl.Gray)
	guiOffsetY += offsetIncrement
	guiOffsetX += offsetIncrementX
	for i := range buffer.NoiseGain {
		labelSizeX := float32(65)
		rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, labelSizeX, 18), goNoise.NOISETYPE_STRING[*buffer.NoiseTypes[i]], rl.RayWhite, rl.Gray, rl.Gray)

		rect = rl.NewRectangle(guiOffsetX+labelSizeX, guiOffsetY, 100, 18)
		oldWeight := *buffer.NoiseGain[i]
		*buffer.NoiseGain[i] = float64(rayGui.SliderBar(rect, float32(oldWeight), 0.0, 1.0))
		if *buffer.NoiseGain[i] != oldWeight {
			changed = true
		}
		rayGui.Label(rect, fmt.Sprintf("%.3f", *buffer.NoiseGain[i]))
		guiOffsetY += offsetIncrement
	}
	guiOffsetX -= offsetIncrementX
	buffer.Changed = changed

	return buffer
}

func main() {

	// Setup game and window
	game := GameContainer{}
	game.CameraSizeX = 1400
	game.CameraSizeY = 720
	rl.InitWindow(
		game.CameraSizeX,
		game.CameraSizeY,
		"rl thingie")

	println("[GAME] Init")
	game.Console = rlConsole.NewConsoleLayer(
		mapWidth, mapHeight,
		rlConsole.NewSpriteSheetFont(
			rl.LoadTexture("resources/bitmapfont/ascii_8x8.png"),
			fontWidth,
			fontHeight))

	//TODO: Make own style
	//rayGui.LoadGuiStyle("")

	//Add ability to load images as heightmap
	//heightmapImg := rl.LoadImage("resources/bitmapfont/ascii_8x8.png")

	// World setup
	continentMaskParams := goNoise.NoiseParams{
		ScaleX:      100.0,
		ScaleY:      70.0,
		Octaves:     3,
		Persistance: 0.5,
		Lacunarity:  2.0,
		Gain:        0.3,
		Bias:        0.45}

	mountainMaskParams := goNoise.NoiseParams{
		ScaleX:      50.0,
		ScaleY:      30.0,
		Octaves:     3,
		Persistance: 0.5,
		Lacunarity:  2.0,
		Gain:        0.3,
		Bias:        0.3}

	mountainNoiseParams := goNoise.NoiseParams{
		ScaleX:      15.0,
		ScaleY:      10.0,
		Octaves:     4,
		Persistance: 0.55,
		Lacunarity:  2.2,
		Gain:        0.35,
		Bias:        0.35}

	landNoiseParams := []goNoise.NoiseParams{
		goNoise.NoiseParams{
			Type:        goNoise.FBM,
			Weight:      0.5,
			ScaleX:      40.0,
			ScaleY:      20.0,
			Octaves:     5,
			Persistance: 0.6,
			Lacunarity:  2.5,
			Gain:        0.5,
			Bias:        0.5},
		goNoise.NoiseParams{
			Type:        goNoise.RIDGED,
			Weight:      0.5,
			ScaleX:      16.0,
			ScaleY:      8.0,
			Octaves:     4,
			Persistance: 0.5,
			Lacunarity:  2.0,
			Gain:        0.7,
			Bias:        0.4},
		goNoise.NoiseParams{
			Type:        goNoise.NEGATIVE_RIDGED,
			Weight:      0.9,
			ScaleX:      48.0,
			ScaleY:      38.0,
			Octaves:     5,
			Persistance: 0.5,
			Lacunarity:  2.0,
			Gain:        0.5,
			Bias:        0.5},
		goNoise.NoiseParams{
			Type:        goNoise.PERTURBED,
			Weight:      0.6,
			ScaleX:      25.0,
			ScaleY:      12.0,
			Octaves:     8,
			Persistance: 0.5,
			Lacunarity:  2.0,
			Gain:        0.5,
			Bias:        0.75}}

	worldParams := worldgen.WorldParams{
		Seed:          1,
		Width:         mapWidth,
		Height:        mapHeight,
		Scale:         float64(float32(1.0)),
		WaterHeight:   100,
		FWaterHeight:  float64(100) / 255.0,
		MountainLine:  160,
		FMountainLine: float64(160) / 255.0}

	game.World = worldgen.GenerateWorld(
		worldParams,
		continentMaskParams, mountainMaskParams,
		mountainNoiseParams, landNoiseParams)

	// GUI setup
	guiResponseBuffer := GuiResponseBuffer{
		Seed:         &worldParams.Seed,
		WorldScale:   &worldParams.Scale,
		NoiseTypes:   []*goNoise.NoiseType{},
		NoiseWeights: []*float32{},
		NoiseBias:    []*float64{},
		NoiseGain:    []*float64{}}

	for i := range landNoiseParams {
		guiResponseBuffer.NoiseTypes = append(
			guiResponseBuffer.NoiseTypes,
			&landNoiseParams[i].Type)
		guiResponseBuffer.NoiseWeights = append(
			guiResponseBuffer.NoiseWeights,
			&landNoiseParams[i].Weight)
		guiResponseBuffer.NoiseBias = append(
			guiResponseBuffer.NoiseBias,
			&landNoiseParams[i].Bias)
		guiResponseBuffer.NoiseGain = append(
			guiResponseBuffer.NoiseGain,
			&landNoiseParams[i].Gain)
	}
	// Camera
	camGlobalPos := rl.NewVector2(8, 8)
	camGlobalOffset := rl.NewVector2(0, 0)
	game.Camera = rl.NewCamera2D(camGlobalPos, camGlobalOffset, 0.0, 1.0)

	updateConsoleVisuals(&game)

	rl.SetTargetFPS(60)
	println("[GAME] Main loop")

	// MAIN GAME LOOP
	for !rl.WindowShouldClose() {

		// Input
		inputContext := newInputContext()
		changed, regen := getInput(&inputContext) // TODO: fix this nonsense leftover shit
		regen = guiResponseBuffer.Changed
		guiData := GuiData{inputContext.MouseX, inputContext.MouseY}

		// TODO: Tidy up this whole area
		// Logic
		if regen {

			game.World = worldgen.GenerateWorld(
				worldParams,
				continentMaskParams, mountainMaskParams,
				mountainNoiseParams, landNoiseParams)

			updateConsoleVisuals(&game)
		}

		//Get info for mouse pointer
		worldMousePosX := ((-8 + inputContext.MouseX + int32(game.Camera.Target.X)) / fontWidth)
		worldMousePosY := ((-8 + inputContext.MouseY + int32(game.Camera.Target.Y)) / fontHeight)
		var biomeStr, heightStr, tempStr, moistureStr string
		if worldMousePosX > 0 && worldMousePosX < mapWidth && worldMousePosY > 0 && worldMousePosY < mapHeight {
			biomeStr = game.World.GetBiomeAtPoint(int(worldMousePosX), int(worldMousePosY))
			heightStr = fmt.Sprintf("%.2f", float32(game.World.GetHeightAtPoint(int(worldMousePosX), int(worldMousePosY)))/255.0)
			moistureStr = fmt.Sprintf("%.2f", float32(game.World.GetMoistureAtPoint(int(worldMousePosX), int(worldMousePosY)))/255.0)
			tempStr = fmt.Sprintf("%.2f", float32(game.World.GetTemperatureAtPoint(int(worldMousePosX), int(worldMousePosY)))/255.0)
		}

		if changed {
			setConsoleMode(inputContext)
			updateConsoleVisuals(&game)
		}

		moveCamera(&game, inputContext)

		// Drawing world
		rl.BeginDrawing()
		rl.BeginMode2D(game.Camera)
		{
			rl.ClearBackground(rl.NewColor(0, 0, 0, 255))

			game.Console.Draw(
				int(game.Camera.Target.X)/fontWidth, int(game.Camera.Target.Y)/fontHeight,
				int(game.CameraSizeX)/fontWidth, int(game.CameraSizeY)/fontHeight)

			rl.DrawRectangleLines(
				worldMousePosX*fontHeight,
				worldMousePosY*fontHeight,
				fontWidth, fontHeight, rl.Red)

			rl.DrawText(biomeStr+" -- Height: "+heightStr+"/1.0"+" Temp: "+tempStr+"/1.0"+" Moisture: "+moistureStr+"/1.0",
				worldMousePosX*fontWidth-4, worldMousePosY*fontHeight-12, 8, rl.White)

		}
		rl.EndMode2D()

		// Drawing GUI
		guiResponseBuffer = drawGui(guiData, guiResponseBuffer)

		rl.EndDrawing()
	}
	// Unload all stuff here
	game.Console.Unload()

	rl.CloseWindow()
	println("[GAME] exiting")
}
