package goNoise

import (
	snoise "github.com/ojrac/opensimplex-go"
	"math"
)

func GenFBMJob(c chan *[]float64, width int, height int, context NoiseParams, scale float64) {

	//numRoutines := 8
	//heightInterval := height / numRoutines
	defer wg.Done()

	noise := snoise.NewNormalized(context.Seed)
	noiseMap := make([]float64, width*height)

	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			noiseHeight := getFBM(noise, x, y, context, scale)
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

	c <- &noiseMap
}

func GenRidgedFBMJob(c chan *[]float64, width int, height int, context NoiseParams, scale float64, inverted bool) {
	wg.Done()
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
	c <- &noiseMap
}

func GenPerturbedFBMJob(c chan *[]float64, width int, height int, context NoiseParams, scale float64, perturbAmount float64) {
	wg.Done()
	noise := snoise.NewNormalized(context.Seed)

	//perturbAmount := 0.4
	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0

	noiseMap := make([]float64, width*height)
	perturbContext := NoiseParams{
		ScaleX:      context.ScaleX,
		ScaleY:      context.ScaleY,
		Octaves:     1,
		Persistance: 0.5,
		Lacunarity:  2.0,
		Gain:        0.5,
		Bias:        0.5,
		Seed:        1}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			amplitude := 1.0
			frequency := 1.0
			noiseHeight := 0.0

			perturbation := getFBM(noise, x, y, perturbContext, scale) * perturbAmount

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

	c <- &noiseMap
}
