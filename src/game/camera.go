package game

import "common"

func MoveCamera(game *GameContainer, input InputContext, fontWidth, fontHeight, mapWidth, mapHeight int32) {

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
