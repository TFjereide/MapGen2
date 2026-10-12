package game

import (
	"fmt"
	rayGui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"goNoise"
)

const (
	COLOR int = iota
	HEIGHT
	TEMPERATURE
	MOISTURE
	WATERSHEDS
	FLOW
)

var MODE_NAMES = [...]string{
	"color",
	"height",
	"temperature",
	"moisture",
	"watersheds",
	"flow"}

// Other stuff the Gui needs to draw
type GuiData struct {
	MouseX, MouseY int32
}

// Changes returned by RayGui's lil' widgets
type GuiResponseBuffer struct {
	Changed     bool
	Seed        *int64
	WorldScale  *float64
	NoiseNames  []string
	NoiseParams []*goNoise.NoiseParams
}

func DrawMinimap(game *GameContainer, fontWidth, fontHeight int32, scale float32, drawRect bool) {
	minimapScale := scale
	minimapOffset := rl.NewVector2(
		float32(game.CameraSizeX-int32(float32(game.Console.Texture.Width)*minimapScale)),
		float32(game.CameraSizeY-int32(float32(game.Console.Texture.Height)*minimapScale)))

	rl.DrawTextureEx(game.Console.Texture, minimapOffset, 0.0, float32(minimapScale), rl.White)

	if drawRect {
		rl.DrawRectangleLines(
			int32(minimapOffset.X+float32(game.Camera.Target.X/float32(fontWidth))*minimapScale),  ///float32(game.Console.Texture.Width))),
			int32(minimapOffset.Y+float32(game.Camera.Target.Y/float32(fontHeight))*minimapScale), ///float32(game.Console.Texture.Height))),
			int32(float32(game.CameraSizeX/(fontWidth))*minimapScale),
			int32(float32(game.CameraSizeY/(fontHeight))*minimapScale),
			rl.Red)
	}
}
func DrawGui(buffer *GuiResponseBuffer, mode int) {

	guiOffsetX := float32(8)
	guiOffsetY := float32(8)
	offsetIncrement := float32(18)
	offsetIncrementX := float32(8)
	changed := false

	elementsHeight := float32(0.0)
	for range buffer.NoiseParams {
		elementsHeight += (offsetIncrement * 4)
	}
	rl.DrawRectangle(4, 4, 180, int32(elementsHeight+60.0), rl.Color{125, 125, 125, 125})

	// fps widget
	fps := rl.GetFPS()
	strFPS := fmt.Sprintf("%.2f", fps)

	rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, 0, 0),
		"FPS: "+strFPS,
		rl.RayWhite, rl.Gray, rl.Gray)

	guiOffsetY += offsetIncrement

	// mode widget
	mode_name := MODE_NAMES[mode]
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

	guiOffsetY += offsetIncrement

	for i, _ := range buffer.NoiseParams {

		param := buffer.NoiseParams[i]
		str := buffer.NoiseNames[i]

		//Draw Type
		lStr := str + goNoise.NOISETYPE_STRING[param.Type]
		rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, 100, 18), lStr, rl.RayWhite, rl.Gray, rl.Gray)
		guiOffsetY += offsetIncrement
		guiOffsetX += offsetIncrementX

		//Draw weight
		labelSizeX := float32(40)
		if param.Weight != -1.0 {
			rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, labelSizeX, 18), "Weight", rl.RayWhite, rl.Blank, rl.Blank)
			rect = rl.NewRectangle(labelSizeX+guiOffsetX, guiOffsetY, 100, 18)
			oldWeight := param.Weight
			param.Weight = rayGui.SliderBar(rect, param.Weight, 0.0, 1.0)
			if param.Weight != oldWeight {
				changed = true
			}
			rayGui.Label(rect, fmt.Sprintf("%.3f", param.Weight))
			guiOffsetY += offsetIncrement
		}

		//Draw Bias
		rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, labelSizeX, 18), "Bias", rl.RayWhite, rl.Blank, rl.Blank)
		rect = rl.NewRectangle(labelSizeX+guiOffsetX, guiOffsetY, 100, 18)
		oldBias := param.Bias
		param.Bias = float64(rayGui.SliderBar(rect, float32(param.Bias), 0.0, 1.0))
		if param.Bias != oldBias {
			changed = true
		}
		rayGui.Label(rect, fmt.Sprintf("%.3f", param.Bias))
		guiOffsetY += offsetIncrement

		//Draw Gain
		rayGui.LabelEx(rl.NewRectangle(guiOffsetX, guiOffsetY, labelSizeX, 18), "Gain", rl.RayWhite, rl.Blank, rl.Blank)
		rect = rl.NewRectangle(labelSizeX+guiOffsetX, guiOffsetY, 100, 18)
		oldGain := param.Gain
		param.Gain = float64(rayGui.SliderBar(rect, float32(param.Gain), 0.0, 1.0))
		if param.Gain != oldGain {
			changed = true
		}
		rayGui.Label(rect, fmt.Sprintf("%.3f", param.Gain))
		guiOffsetY += offsetIncrement
		guiOffsetX -= offsetIncrementX

	}
	guiOffsetX -= offsetIncrementX

	guiOffsetX -= offsetIncrementX
	buffer.Changed = changed

}
func drawNoiseParams(params *goNoise.NoiseParams) bool {
	changed := false
	return changed
}
