package game

import "common"

func MoveCamera(game *GameContainer, input InputContext, fontWidth, fontHeight, mapWidth, mapHeight int32) {

	moveDirX := common.BToI(input.MoveRight) - common.BToI(input.MoveLeft)
	moveDirY := common.BToI(input.MoveDown) - common.BToI(input.MoveUp)

	game.Camera.Target.X += float32(moveDirX) * 20.0
	game.Camera.Target.Y += float32(moveDirY) * 20.0

	game.Camera.Target.X = common.ClampF32(game.Camera.Target.X, 0, float32(mapWidth*fontWidth-game.CameraSizeX))
	game.Camera.Target.Y = common.ClampF32(game.Camera.Target.Y, 0, float32(mapHeight*fontHeight-game.CameraSizeY))

}
