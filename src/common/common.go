package common

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func ConvertArrayFloat64ToByte(f []float64) []byte {
	bytemap := make([]byte, len(f))
	for i := 0; i < len(f); i++ {
		bytemap[i] = byte(f[i] * 255)
	}
	return bytemap
}

func ConvertArrayByteToFloat64(b []byte) []float64 {
	floatmap := make([]float64, len(b))
	for i := 0; i < len(b); i++ {
		floatmap[i] = float64(b[i]) / 255
	}
	return floatmap
}

// Returns slice of grayscale/brightness raylib colors
func ConvertArrayByteToRLColor(b []byte) []rl.Color {
	colSlice := make([]rl.Color, len(b))
	for i := 0; i < len(b); i++ {
		elevation := b[i]
		colVal := rl.NewColor(elevation, elevation, elevation, 255)
		colSlice[i] = colVal
	}
	return colSlice

}

func BToI(b bool) int {
	if b {
		return 1
	}
	return 0
}

func BlurColorSlice(colors []rl.Color, width, height int, blurAmount float32) []rl.Color {

	newColors := make([]rl.Color, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			idx := GetIdx(x, y, width)
			currCol := colors[idx]

			for dy := -1; dy < 2; dy++ {
				for dx := -1; dx < 2; dx++ {
					if dx != 0 && dy != 0 {

						if x+dx < 0 || x+dx >= width || y+dy < 0 || y+dy >= height {
							continue
						}
						_idx := GetIdx(x+dx, y+dy, width)
						mixCol := colors[_idx]
						currCol = ColorMixWeighted(mixCol, currCol, blurAmount)

					}
				}
			}
			newColors[idx] = currCol
		}
	}
	return newColors
}

func colorMix(a, b rl.Color) rl.Color {

	col := rl.NewColor(0, 0, 0, 0)

	col.R = (a.R + b.R) / 2
	col.G = (a.G + b.G) / 2
	col.B = (a.B + b.B) / 2
	col.A = (a.A + b.A) / 2

	return col

}

// Color multiply operation. Disregards alpha
func ColorMultiply(a, b rl.Color) rl.Color {

	//from byte to normalized float values
	fr := (float32(a.R) / 255.0) * (float32(b.R) / 255.0) * 255
	fg := (float32(a.G) / 255.0) * (float32(b.G) / 255.0) * 255
	fb := (float32(a.B) / 255.0) * (float32(b.B) / 255.0) * 255

	return rl.NewColor(byte(fr), byte(fg), byte(fb), 255)
}

// fades one color into another by wight
func ColorMixWeighted(fc, sc rl.Color, weight float32) rl.Color {

	col := rl.NewColor(0, 0, 0, 0)

	col.R = MixByteByWeight(fc.R, sc.R, weight)
	col.G = MixByteByWeight(fc.G, sc.G, weight)
	col.B = MixByteByWeight(fc.B, sc.B, weight)
	col.A = MixByteByWeight(fc.A, sc.A, weight)
	return col
}

func ColorAddWeighted(fc, sc rl.Color, weight float32) rl.Color {
	col := rl.NewColor(0, 0, 0, 0)
	col.R = AddByteByWeight(fc.R, sc.R, weight)
	col.G = AddByteByWeight(fc.G, sc.G, weight)
	col.B = AddByteByWeight(fc.B, sc.B, weight)
	col.A = AddByteByWeight(fc.A, sc.A, weight)
	return col
}
func ColorFloat32Multiply(c rl.Color, brightness float32) rl.Color {
	c.R = byte(ClampF32(float32(c.R)*brightness, 0.0, 255.0))
	c.G = byte(ClampF32(float32(c.G)*brightness, 0.0, 255.0))
	c.B = byte(ClampF32(float32(c.B)*brightness, 0.0, 255.0))
	return c
}
func ColorByteMultiply(c rl.Color, brightness byte) rl.Color {
	bBrightness := float32(brightness) / 255.0
	c.R = byte(ClampF32(float32(c.R)*bBrightness, 0.0, 255.0))
	c.G = byte(ClampF32(float32(c.G)*bBrightness, 0.0, 255.0))
	c.B = byte(ClampF32(float32(c.B)*bBrightness, 0.0, 255.0))
	return c
}

func MixByteByWeight(a, b byte, weight float32) byte {
	result := float32(a) * weight
	result += float32(b)
	result /= (1.0 + weight)

	return byte(result)
}
func AddByteByWeight(a, b byte, weight float32) byte {
	result := float32(a) * weight
	result += float32(b)
	if result > 255.0 {
		result = 255.0
	}
	return byte(result)
}

func ByteToGrayColor(v byte) rl.Color {
	return rl.NewColor(v, v, v, 255)
}
func Float64ToGrayColor(v float64) rl.Color {
	return rl.NewColor(uint8(v*255), uint8(v*255), uint8(v*255), 255)
}

func GetIdx(x, y, width int) int {
	return ((y * width) + x)
}
func ClampF32(val, bottom, top float32) float32 {
	if val > top {
		return top
	} else if val < bottom {
		return bottom
	}
	return val
}
func ClampF64(val, bottom, top float64) float64 {
	if val > top {
		return top
	} else if val < bottom {
		return bottom
	}
	return val
}

// Will curve the value towards a value in a smooth fashion
func GetBias(time, bias float64) float64 {
	return (time / ((((1.0 / bias) - 2.0) * (1.0 - time)) + 1.0))
}

// Will curve value above and below 0.5 towards another
func GetGain(time, gain float64) float64 {
	if time < 0.5 {
		return GetBias(time*2.0, gain) / 2.0
	} else {
		return GetBias(time*2.0-1.0, 1.0-gain)/2.0 + 0.5
	}
}
