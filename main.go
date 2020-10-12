package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	//imgui "github.com/inkyblackness/imgui-go"
	g "game"
	"goNoise"
	"math/rand"

	rlConsole "rlConsole"
	worldgen "worldgenerator"
)

const fontWidth = 8
const fontHeight = 8
const mapWidth = 512  //512
const mapHeight = 256 //256

func main() {

	rand.Seed(1)

	// Setup game and window
	MODE := g.COLOR
	game := g.GameContainer{}
	game.CameraSizeX = 1920
	game.CameraSizeY = 1000
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
		ScaleX:      150.0,
		ScaleY:      90.0,
		Octaves:     3,
		Persistance: 0.50,
		Lacunarity:  2.0,
		Gain:        0.25,
		Bias:        0.45}

	mountainMaskParams := goNoise.NoiseParams{
		ScaleX:      60.0,
		ScaleY:      40.0,
		Octaves:     7,
		Persistance: 0.5,
		Lacunarity:  2.0,
		Gain:        0.20,
		Bias:        0.35}

	mountainNoiseParams := []goNoise.NoiseParams{
		goNoise.NoiseParams{
			Type:        goNoise.NEGATIVE_RIDGED,
			Weight:      1.0,
			ScaleX:      15.0,
			ScaleY:      10.0,
			Octaves:     5,
			Persistance: 0.54,
			Lacunarity:  2.0,
			Gain:        0.35,
			Bias:        0.35},
		goNoise.NoiseParams{
			Type:        goNoise.RIDGED,
			Weight:      0.7,
			ScaleX:      25.0,
			ScaleY:      12.5,
			Octaves:     5,
			Persistance: 0.50,
			Lacunarity:  2.0,
			Gain:        0.2,
			Bias:        0.6}}

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
		Seed:          34,
		Width:         mapWidth,
		Height:        mapHeight,
		Scale:         float64(float32(2.0)),
		WaterHeight:   100,
		FWaterHeight:  float64(100) / 255.0,
		MountainLine:  160,
		FMountainLine: float64(160) / 255.0}

	game.World = worldgen.GenerateWorld(
		worldParams,
		continentMaskParams, mountainMaskParams,
		mountainNoiseParams, landNoiseParams)

	// GUI setup
	guiResponseBuffer := g.GuiResponseBuffer{
		Seed:         &worldParams.Seed,
		WorldScale:   &worldParams.Scale,
		NoiseTypes:   []*goNoise.NoiseType{},
		NoiseWeights: []*float32{},
		NoiseBias:    []*float64{},
		NoiseGain:    []*float64{}}

	for i, _ := range landNoiseParams {
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

	g.UpdateConsoleVisuals(&game, MODE)

	// Camera
	camGlobalPos := rl.NewVector2(8, 8)
	camGlobalOffset := rl.NewVector2(0, 0)
	game.Camera = rl.NewCamera2D(camGlobalPos, camGlobalOffset, 0.0, 1.0)

	rl.SetTargetFPS(60)
	println("[GAME] Main loop")

	// MAIN GAME LOOP
	for !rl.WindowShouldClose() {

		// Input
		inputContext := g.NewInputContext()
		changed, regen := g.GetInput(&inputContext) // TODO: fix this nonsense leftover shit
		regen = guiResponseBuffer.Changed
		guiData := g.GuiData{inputContext.MouseX, inputContext.MouseY}

		// TODO: Tidy up this whole area
		// Logic
		if regen {

			game.World = worldgen.GenerateWorld(
				worldParams,
				continentMaskParams, mountainMaskParams,
				mountainNoiseParams, landNoiseParams)

			g.UpdateConsoleVisuals(&game, MODE)
		}

		if changed {
			MODE = g.SetConsoleMode(inputContext)
			g.UpdateConsoleVisuals(&game, MODE)
		}
		g.MoveCamera(&game, inputContext, fontWidth, fontHeight, mapWidth, mapHeight)

		// Drawing world
		rl.BeginDrawing()
		rl.BeginMode2D(game.Camera)
		{
			rl.ClearBackground(rl.NewColor(0, 0, 0, 255))

			game.Console.Draw(
				int(game.Camera.Target.X)/fontWidth, int(game.Camera.Target.Y)/fontHeight,
				int(game.CameraSizeX)/fontWidth, int(game.CameraSizeY)/fontHeight)

			g.DrawMousePointer(&game, inputContext, fontWidth, fontHeight, mapWidth, mapHeight)

		}
		rl.EndMode2D()

		// Drawing GUI
		guiResponseBuffer = g.DrawGui(guiData, guiResponseBuffer, MODE)
		g.DrawMinimap(&game, fontWidth, fontHeight, 0.5, true)
		// Minimap

		rl.EndDrawing()
	}
	// Unload all stuff here
	game.Console.Unload()

	rl.CloseWindow()
	println("[GAME] exiting")
}
