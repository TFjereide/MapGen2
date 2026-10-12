package goNoise

import snoise "github.com/ojrac/opensimplex-go"

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

func getIdx(x, y, width int) int {
	return ((y * width) + x)
}

// Helper Func for a single point
func getFBM(pnoise *snoise.Noise, x, y int, context NoiseParams, scale float64) float64 {
	amplitude := 1.0
	frequency := 1.0
	noiseHeight := 0.0
	noise := *pnoise

	// get new value
	for i := 0; i < context.Octaves; i++ {
		sampleX := float64(x) / (context.ScaleX * scale) * frequency
		sampleY := float64(y) / (context.ScaleY * scale) * frequency

		nValue := noise.Eval2(sampleX, sampleY)
		noiseHeight += nValue * amplitude
		amplitude *= context.Persistance
		frequency *= context.Lacunarity
	}

	return noiseHeight
}

// Helper Func for a single point
func getFloat64FBM(pnoise *snoise.Noise, x, y float64, context NoiseParams, scale float64) float64 {
	amplitude := 1.0
	frequency := 1.0
	noiseHeight := 0.0
	noise := *pnoise

	// get new value
	for i := 0; i < context.Octaves; i++ {
		sampleX := float64(x) / (context.ScaleX * scale) * frequency
		sampleY := float64(y) / (context.ScaleY * scale) * frequency

		nValue := noise.Eval2(sampleX, sampleY)
		noiseHeight += nValue * amplitude
		amplitude *= context.Persistance
		frequency *= context.Lacunarity
	}

	return noiseHeight
}
