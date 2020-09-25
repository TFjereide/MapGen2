package goNoise

import (
	snoise "github.com/ojrac/opensimplex-go"
	"math"
)

type NoiseType uint8

const (
	FBM NoiseType = iota
	RIDGED
	NEGATIVE_RIDGED
	PERTURBED
)

var NOISETYPE_STRING = [...]string{
	"FBM",
	"+RIDGED",
	"-RIDGED",
	"PERTURBED"}

type NoiseParams struct {
	Type        NoiseType
	Weight      float32
	ScaleX      float64
	ScaleY      float64
	Scale       float64
	ScaleRatio  float64
	Octaves     int
	Persistance float64
	Lacunarity  float64
	Gain        float64
	Bias        float64
	Seed        int64
}

// Returns more interesting, varied noise
func GetLayeredNoise(width, height int, parameters []NoiseParams, scale float64) *[]float64 {

	numLayers := len(parameters)
	noiseLayers := make([][]float64, numLayers)

	// Generate Noises
	for i := 0; i < numLayers; i++ {
		// Offset the seed for more variety
		if parameters[i].Type == FBM {
			noiseLayers[i] = *GenFBM(width, height, parameters[i], scale)
		} else if parameters[i].Type == RIDGED {
			noiseLayers[i] = *GenRidgedFBM(width, height, parameters[i], scale, false)
		} else if parameters[i].Type == NEGATIVE_RIDGED {
			noiseLayers[i] = *GenRidgedFBM(width, height, parameters[i], scale, true)
		} else if parameters[i].Type == PERTURBED {
			noiseLayers[i] = *GenPerturbedFBM(width, height, parameters[i], scale, 0.1)
		}
	}

	floatVals := make([]float64, width*height)
	weightSum := 0.0
	for i := 0; i < len(parameters); i++ {
		weightSum += float64(parameters[i].Weight)
	}
	for i := 0; i < width*height; i++ {
		for j := 0; j < numLayers; j++ {
			floatVals[i] += (noiseLayers[j][i] * float64(parameters[j].Weight))
		}
		floatVals[i] /= weightSum
	}

	return &floatVals
}

func getNoise(noise snoise.Noise, x, y int, scaleX, scaleY float64, scale float64) float64 {
	return noise.Eval2(float64(x)*scaleX, float64(y)*scaleY)
}

func GenFBM(width int, height int, context NoiseParams, scale float64) *[]float64 {

	// Turns out this is super expensive. Increases time by 30+ ms per noise thing
	//layers := make([]snoise.Noise, context.Octaves)
	//for i := 0; i < context.Octaves; i++ {
	//	layers[i] = snoise.NewNormalized(int64(i) + context.Seed)
	//}

	noise := snoise.NewNormalized(context.Seed)
	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0

	noiseMap := make([]float64, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			amplitude := 1.0
			frequency := 1.0
			noiseHeight := 0.0

			// get new value
			for i := 0; i < context.Octaves; i++ {
				sampleX := float64(x) / (context.ScaleX * scale) * frequency
				sampleY := float64(y) / (context.ScaleY * scale) * frequency

				nValue := noise.Eval2(sampleX, sampleY)
				noiseHeight += nValue * amplitude
				amplitude *= context.Persistance
				frequency *= context.Lacunarity

			}

			if noiseHeight > maxNoiseHeight {
				maxNoiseHeight = noiseHeight
			} else if noiseHeight < minNoiseHeight {
				minNoiseHeight = noiseHeight
			}

			noiseMap[y*width+x] = noiseHeight

		}
	}

	for i := 0; i < len(noiseMap); i++ {
		noiseMap[i] /= maxNoiseHeight
		noiseMap[i] = GetBias(noiseMap[i], context.Bias)
		noiseMap[i] = GetGain(noiseMap[i], context.Gain)
	}

	return &noiseMap
}
func GenRidgedFBM(width int, height int, context NoiseParams, scale float64, inverted bool) *[]float64 {

	//layers := make([]snoise.Noise, context.Octaves)
	//for i := 0; i < context.Octaves; i++ {
	//	layers[i] = snoise.NewNormalized(int64(i) + context.Seed)
	//}

	noise := snoise.NewNormalized(context.Seed)

	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0

	noiseMap := make([]float64, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			amplitude := 1.0
			frequency := 1.0
			noiseHeight := 0.0

			// get new value
			for i := 0; i < context.Octaves; i++ {
				sampleX := float64(x) / (context.ScaleX * scale) * frequency
				sampleY := float64(y) / (context.ScaleY * scale) * frequency

				nValue := noise.Eval2(sampleX, sampleY)
				noiseHeight += nValue * amplitude
				amplitude *= context.Persistance
				frequency *= context.Lacunarity

			}

			if noiseHeight > maxNoiseHeight {
				maxNoiseHeight = noiseHeight
			} else if noiseHeight < minNoiseHeight {
				minNoiseHeight = noiseHeight
			}

			noiseMap[y*width+x] = noiseHeight

		}
	}
	for i := 0; i < height*width; i++ {
		noiseMap[i] /= maxNoiseHeight
		noiseMap[i] = math.Abs((noiseMap[i] * 2.0) - 1.0)
		if inverted {
			noiseMap[i] = 1.0 - noiseMap[i]
		}
		noiseMap[i] = GetBias(noiseMap[i], context.Bias)
		noiseMap[i] = GetGain(noiseMap[i], context.Gain)

	}
	return &noiseMap
}

func GenPerturbedFBM(width int, height int, context NoiseParams, scale float64, perturbAmount float64) *[]float64 {

	//layers := make([]snoise.Noise, context.Octaves)
	//for i := 0; i < context.Octaves; i++ {
	//	layers[i] = snoise.NewNormalized(int64(i) + context.Seed)
	//}

	noise := snoise.NewNormalized(context.Seed)

	//perturbAmount := 0.4
	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0

	noiseMap := make([]float64, width*height)
	perturbContext := NoiseParams{
		ScaleX:      context.ScaleX,
		ScaleY:      context.ScaleY,
		Octaves:     3,
		Persistance: 0.5,
		Lacunarity:  2.0,
		Gain:        0.5,
		Bias:        0.5,
		Seed:        1}

	perturbMap := *GenFBM(width, height, perturbContext, 1.0)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			amplitude := 1.0
			frequency := 1.0
			noiseHeight := 0.0

			perturbation := perturbMap[getIdx(x, y, width)] * perturbAmount

			// get new value
			for i := 0; i < context.Octaves; i++ {
				sampleX := float64(x+200) / (context.ScaleX * scale) * (frequency * perturbation) //perturbAmount * fbm(layers, float64(-x), float64(-y), context.ScaleX, context.ScaleY, context.Octaves, context.Persistance, context.Lacunarity))
				sampleY := float64(y+200) / (context.ScaleY * scale) * (frequency * perturbation) //perturbAmount * fbm(layers, float64(-x), float64(-y), context.ScaleX, context.ScaleY, context.Octaves, context.Persistance, context.Lacunarity))

				nValue := noise.Eval2(sampleX, sampleY)
				noiseHeight += nValue * amplitude
				amplitude *= context.Persistance
				frequency *= context.Lacunarity

			}

			if noiseHeight > maxNoiseHeight {
				maxNoiseHeight = noiseHeight
			} else if noiseHeight < minNoiseHeight {
				minNoiseHeight = noiseHeight
			}

			noiseMap[y*width+x] = noiseHeight

		}
	}

	for i := 0; i < len(noiseMap); i++ {
		noiseMap[i] /= maxNoiseHeight
		noiseMap[i] = GetBias(noiseMap[i], context.Bias)
		noiseMap[i] = GetGain(noiseMap[i], context.Gain)
	}

	return &noiseMap
}

func fbm(layers []snoise.Noise, x, y, scaleX, scaleY float64, octaves int, persistance, lacunarity float64) float64 {
	amplitude := 1.0
	frequency := 1.0
	noiseHeight := 0.0

	// get new value
	for i := 0; i < octaves; i++ {
		sampleX := x / scaleX * frequency
		sampleY := y / scaleY * frequency

		nValue := layers[i].Eval2(sampleX, sampleY)
		noiseHeight += nValue * amplitude
		amplitude *= persistance
		frequency *= lacunarity

	}
	return noiseHeight
}

func ridgedfbm(layers []snoise.Noise, x, y, scaleX, scaleY float64, octaves int, persistance, lacunarity float64) float64 {
	amplitude := 1.0
	frequency := 1.0
	noiseHeight := 0.0

	// get new value
	for i := 0; i < octaves; i++ {
		sampleX := x / scaleX * frequency
		sampleY := y / scaleY * frequency

		nValue := layers[i].Eval2(sampleX, sampleY)*2.0 - 1.0

		noiseHeight += nValue * amplitude
		amplitude *= persistance
		frequency *= lacunarity

	}
	return 1.0 - math.Abs(noiseHeight)
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

func getIdx(x, y, width int) int {
	return ((y * width) + x)
}
