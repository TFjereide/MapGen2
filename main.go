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
const mapWidth = 1024 //512
const mapHeight = 512 //256
const SCREEN_WIDTH = 1920
const SCREEN_HEIGHT = 1000

func main() {
	rand.Seed(1)
	rl.InitWindow(
		SCREEN_WIDTH,
		SCREEN_HEIGHT,
		"rl thingie")

	//Add ability to load images as heightmap
	//heightmapImg := rl.LoadImage("resources/bitmapfont/EuropeMenasia.jpg")

	println("[GAME] Init")

	camGlobalPos := rl.NewVector2(8, 8)
	camGlobalOffset := rl.NewVector2(0, 0)
	MODE := g.COLOR
	game := g.GameContainer{
		Mode:        g.COLOR,
		CameraSizeX: SCREEN_WIDTH,
		CameraSizeY: SCREEN_HEIGHT,
		Camera:      rl.NewCamera2D(camGlobalPos, camGlobalOffset, 0.0, 1.0),
		Console: rlConsole.NewConsoleLayer(
			mapWidth, mapHeight,
			rlConsole.NewSpriteSheetFont(
				rl.LoadTexture("resources/bitmapfont/ascii_8x8.png"),
				fontWidth,
				fontHeight)),
		InputActions: g.NewInputContext(),
		WorldParams:  *initWorldParams(),
	}

	game.GuiResponse = *initGuiResponseBuffer(&game.WorldParams)
	g.GetInput(&game.InputActions)
	game.World = worldgen.GenerateWorld(game.WorldParams)
	gameUpdate(&game)
	g.UpdateConsoleVisuals(&game, MODE)

	rl.SetTargetFPS(120)
	println("[GAME] Main loop")

	// MAIN GAME LOOP
	for !rl.WindowShouldClose() {

		// Input
		game.InputActions = g.NewInputContext()
		g.GetInput(&game.InputActions)

		gameUpdate(&game)

		// Drawing world
		rl.BeginDrawing()
		rl.BeginMode2D(game.Camera)
		{
			rl.ClearBackground(rl.NewColor(0, 0, 0, 255))
			game.Console.Draw(
				int(game.Camera.Target.X)/fontWidth, int(game.Camera.Target.Y)/fontHeight,
				int(game.CameraSizeX)/fontWidth, int(game.CameraSizeY)/fontHeight)
			g.DrawMousePointer(&game, fontWidth, fontHeight, mapWidth, mapHeight)

		}
		rl.EndMode2D()
		// Drawing GUI
		g.DrawGui(&game.GuiResponse, game.Mode)
		g.DrawMinimap(&game, fontWidth, fontHeight, 0.5, true)

		rl.EndDrawing()
	}
	// Unload all stuff here
	game.Console.Unload()

	rl.CloseWindow()
	println("[GAME] exiting")
}

func gameUpdate(game *g.GameContainer) {

	if game.GuiResponse.Changed {
		game.World = worldgen.GenerateWorld(game.WorldParams)
		g.UpdateConsoleVisuals(game, game.Mode)
		println("GENERATE")

	}
	if game.InputActions.ModeChanged {
		game.Mode = g.SetConsoleMode(game.InputActions)
		g.UpdateConsoleVisuals(game, game.Mode)
	}
	g.MoveCamera(game, game.InputActions, fontWidth, fontHeight, mapWidth, mapHeight)

}

func initWorldParams() *worldgen.WorldParams {

	// World setup
	worldParams := worldgen.WorldParams{
		Seed:          53,
		Width:         mapWidth,
		Height:        mapHeight,
		Scale:         float64(float32(2.0)),
		WaterHeight:   100,
		FWaterHeight:  float64(100) / 255.0,
		MountainLine:  160,
		FMountainLine: float64(160) / 255.0,

		ContinentMaskParams: goNoise.NoiseParams{
			Type:        goNoise.FBM,
			Weight:      -1.0,
			ScaleX:      150.0,
			ScaleY:      90.0,
			Octaves:     4,
			Persistance: 0.5,
			Lacunarity:  2.0,
			Gain:        0.315,
			Bias:        0.315},
		MountainMaskParams: goNoise.NoiseParams{
			Type:        goNoise.FBM,
			Weight:      -1.0,
			ScaleX:      60.0,
			ScaleY:      40.0,
			Octaves:     7,
			Persistance: 0.5,
			Lacunarity:  2.0,
			Gain:        0.20,
			Bias:        0.35},

		MountainNoiseParams: []goNoise.NoiseParams{
			goNoise.NoiseParams{
				Type:        goNoise.NEGATIVE_RIDGED,
				Weight:      0.7,
				ScaleX:      15.0,
				ScaleY:      10.0,
				Octaves:     5,
				Persistance: 0.54,
				Lacunarity:  2.0,
				Gain:        0.48,
				Bias:        0.35},
			goNoise.NoiseParams{
				Type:        goNoise.RIDGED,
				Weight:      0.65,
				ScaleX:      25.0,
				ScaleY:      12.5,
				Octaves:     5,
				Persistance: 0.50,
				Lacunarity:  2.0,
				Gain:        0.33,
				Bias:        0.70},
			goNoise.NoiseParams{
				Type:        goNoise.PERTURBED,
				Weight:      0.4,
				ScaleX:      60.0,
				ScaleY:      30.0,
				Octaves:     6,
				Persistance: 0.5,
				Lacunarity:  2.0,
				Gain:        0.5,
				Bias:        0.133}},

		LandNoiseParams: []goNoise.NoiseParams{
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
				Weight:      0.8,
				ScaleX:      16.0,
				ScaleY:      8.0,
				Octaves:     4,
				Persistance: 0.5,
				Lacunarity:  2.0,
				Gain:        0.5,
				Bias:        0.5},
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
				ScaleX:      60.0,
				ScaleY:      30.0,
				Octaves:     6,
				Persistance: 0.5,
				Lacunarity:  2.0,
				Gain:        0.5,
				Bias:        0.5}},
	}

	return &worldParams
}

func initGuiResponseBuffer(worldParams *worldgen.WorldParams) *g.GuiResponseBuffer {

	// GUI setup
	guiResponseBuffer := g.GuiResponseBuffer{
		Seed:        &worldParams.Seed,
		WorldScale:  &worldParams.Scale,
		NoiseNames:  make([]string, 0),
		NoiseParams: make([]*goNoise.NoiseParams, 0)}

	for i, _ := range worldParams.LandNoiseParams {
		guiResponseBuffer.NoiseNames = append(guiResponseBuffer.NoiseNames, "Land ")
		guiResponseBuffer.NoiseParams = append(guiResponseBuffer.NoiseParams, &worldParams.LandNoiseParams[i])
	}

	for i, _ := range worldParams.MountainNoiseParams {
		guiResponseBuffer.NoiseNames = append(guiResponseBuffer.NoiseNames, "Mt. ")
		guiResponseBuffer.NoiseParams = append(guiResponseBuffer.NoiseParams, &worldParams.MountainNoiseParams[i])
	}

	guiResponseBuffer.NoiseNames = append(guiResponseBuffer.NoiseNames, "Mt.Mask ")
	guiResponseBuffer.NoiseParams = append(guiResponseBuffer.NoiseParams, &worldParams.MountainMaskParams)

	guiResponseBuffer.NoiseNames = append(guiResponseBuffer.NoiseNames, "cont.Mask ")
	guiResponseBuffer.NoiseParams = append(guiResponseBuffer.NoiseParams, &worldParams.ContinentMaskParams)

	return &guiResponseBuffer
}
