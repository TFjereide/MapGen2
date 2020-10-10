package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
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
