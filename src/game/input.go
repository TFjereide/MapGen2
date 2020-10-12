package game

import (
	"common"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

//  Actions performed through normal inputs
type InputContext struct {
	MouseX, MouseY     int32
	MoveLeft           bool
	MoveRight          bool
	MoveUp             bool
	MoveDown           bool
	WorldSizeIncrease  bool
	WorldSizeDecrease  bool
	SetModeColor       bool
	SetModeHeight      bool
	SetModeTemperature bool
	SetModeMoisture    bool
	SetModeWatersheds  bool
	SetModeFlow        bool
}

func NewInputContext() InputContext {
	return InputContext{
		MoveLeft:           false,
		MoveRight:          false,
		MoveUp:             false,
		MoveDown:           false,
		SetModeColor:       false,
		SetModeHeight:      false,
		SetModeTemperature: false,
		SetModeMoisture:    false,
		SetModeWatersheds:  false,
		SetModeFlow:        false}
}

func GetInput(inputContext *InputContext) (bool, bool) {

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
		inputContext.SetModeWatersheds = true
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

func SetConsoleMode(input InputContext) int {
	if input.SetModeColor {
		return COLOR
	} else if input.SetModeHeight {
		return HEIGHT
	} else if input.SetModeTemperature {
		return TEMPERATURE
	} else if input.SetModeMoisture {
		return MOISTURE
	} else if input.SetModeWatersheds {
		return WATERSHEDS
	} else if input.SetModeFlow {
		return FLOW
	}
	return COLOR
}

func UpdateConsoleVisuals(game *GameContainer, mode int) {
	switch mode {
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

	case WATERSHEDS:
		game.Console.RenderFG = false
		game.Console.RenderBG = true
		game.Console.RenderVectors = false

		//bytearray := make([]byte, len(game.World.DistanceToLand))
		//for i, v := range game.World.DistanceToLand {
		//	bytearray[i] = v * 10
		//}
		colbuffer := make([]rl.Color, len(game.World.Watershed))

		for i, v := range game.World.Watershed {
			col := game.Console.RandColors[v]
			col.A = uint8(((game.World.FlowStrength[i] * 0.5) + 0.5) * 255.0)

			colbuffer[i] = col
		}
		game.Console.SetBG(colbuffer)
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

func DrawMousePointer(game *GameContainer, inputContext InputContext, fontWidth, fontHeight, mapWidth, mapHeight int32) {
	//Get info for mouse pointer
	worldMousePosX := ((-8 + inputContext.MouseX + int32(game.Camera.Target.X)) / fontWidth)
	worldMousePosY := ((-8 + inputContext.MouseY + int32(game.Camera.Target.Y)) / fontHeight)
	if worldMousePosX < 0 || worldMousePosX > game.CameraSizeX {
		worldMousePosX = 0
	}
	if worldMousePosY < 0 || worldMousePosY > game.CameraSizeY {
		worldMousePosY = 0
	}

	var biomeStr, heightStr, tempStr, moistureStr, distStr string
	if worldMousePosX > 0 && worldMousePosX < mapWidth && worldMousePosY > 0 && worldMousePosY < mapHeight {
		biomeStr = game.World.GetBiomeAtPoint(int(worldMousePosX), int(worldMousePosY))
		heightStr = fmt.Sprintf("%.2f", float32(game.World.GetHeightAtPoint(int(worldMousePosX), int(worldMousePosY)))/255.0)
		moistureStr = fmt.Sprintf("%.2f", float32(game.World.GetMoistureAtPoint(int(worldMousePosX), int(worldMousePosY)))/255.0)
		tempStr = fmt.Sprintf("%.2f", float32(game.World.GetTemperatureAtPoint(int(worldMousePosX), int(worldMousePosY)))/255.0)
		distStr = fmt.Sprintf("%d", game.World.GetDistToLandAtPoint(int(worldMousePosX), int(worldMousePosY)))
	}

	rl.DrawText(biomeStr+" -- Height: "+heightStr+"/1.0"+" Temp: "+tempStr+"/1.0"+" Moisture: "+moistureStr+"/1.0 -d:"+distStr,
		worldMousePosX*fontWidth-4, worldMousePosY*fontHeight-12, 8, rl.White)

	rl.DrawRectangleLines(
		worldMousePosX*fontHeight,
		worldMousePosY*fontHeight,
		fontWidth, fontHeight, rl.Red)

}
