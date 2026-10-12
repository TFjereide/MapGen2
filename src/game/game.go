package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	rlConsole "rlConsole"
	worldgen "worldgenerator"
)

//Contains the entire game state, for easy access
type GameContainer struct {
	World       *worldgen.World
	WorldParams worldgen.WorldParams

	//camera variables
	Camera      rl.Camera2D
	CameraSizeX int32
	CameraSizeY int32

	//Console
	Mode    int
	Console *rlConsole.ConsoleLayer

	//Gui
	GuiResponse GuiResponseBuffer
	//Input
	InputActions InputContext
}
