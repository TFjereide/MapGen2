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
func DrawGui(data GuiData, buffer GuiResponseBuffer, mode int) GuiResponseBuffer {

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
