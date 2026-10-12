package goNoise

import (
	snoise "github.com/ojrac/opensimplex-go"
	"math"
	"sync"
)

// DO NOT USE not ready yet boi
// Returns noise by multirow chunks
func GenFBMChunks(width int, height int, context NoiseParams, scale float64) *[]float64 {

	noise := snoise.NewNormalized(context.Seed)
	noiseMap := make([]float64, width*height)

	numchunk := 8
	chunksize := height / numchunk

	var chanBuffer sync.WaitGroup
	chunks := make(chan []float64, numchunk)

	for chnk := 0; chnk < numchunk; {
		chanBuffer.Add(1)
		go getFBMchunk(
			&chanBuffer, chunks,
			&noise,
			chnk*chunksize, width, chunksize,
			context, scale)
	}

	chanBuffer.Wait()
	close(chunks)

	// needs optimizing
	{
		for item := range chunks {
			yOffset := int(item[0])
			for i := 1; i < chunksize+1; i++ {
				x := (i - 1) % width
				y := yOffset + ((i - 1) / width)
				noiseMap[(x-1)+(y*width)] = item[x]
			}
		}
	}
	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0
	for i := 0; i < width*height; i++ {
		noiseHeight := noiseMap[i]
		if noiseHeight > maxNoiseHeight {
			maxNoiseHeight = noiseHeight
		} else if noiseHeight < minNoiseHeight {
			minNoiseHeight = noiseHeight
		}
	}

	for i := 0; i < len(noiseMap); i++ {
		noiseMap[i] /= maxNoiseHeight
		noiseMap[i] = GetBias(noiseMap[i], context.Bias)
		noiseMap[i] = GetGain(noiseMap[i], context.Gain)
	}

	return &noiseMap
}

func GenFBMParalell(width int, height int, context NoiseParams, scale float64) *[]float64 {

	noise := snoise.NewNormalized(context.Seed)
	noiseMap := make([]float64, width*height)

	var chanBuffer sync.WaitGroup
	rows := make(chan []float64, height)

	for y := 0; y < height; y++ {
		chanBuffer.Add(1)
		go getFBMRow(
			&chanBuffer, rows,
			&noise,
			y, width,
			context, scale)
	}

	chanBuffer.Wait()
	close(rows)

	// needs optimizing
	{
		for item := range rows {
			y := int(item[0])
			for x := 1; x < width+1; x++ {
				noiseMap[(x-1)+(y*width)] = item[x]
			}
		}
	}
	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0
	for i := 0; i < width*height; i++ {
		noiseHeight := noiseMap[i]
		if noiseHeight > maxNoiseHeight {
			maxNoiseHeight = noiseHeight
		} else if noiseHeight < minNoiseHeight {
			minNoiseHeight = noiseHeight
		}
	}

	for i := 0; i < len(noiseMap); i++ {
		noiseMap[i] /= maxNoiseHeight
		noiseMap[i] = GetBias(noiseMap[i], context.Bias)
		noiseMap[i] = GetGain(noiseMap[i], context.Gain)
	}

	return &noiseMap
}

func GenRidgedFBMParalell(width int, height int, context NoiseParams, scale float64, inverted bool, perturbAmount float64) *[]float64 {
	noise := snoise.NewNormalized(context.Seed)
	noiseMap := make([]float64, width*height)

	var chanBuffer sync.WaitGroup
	rows := make(chan []float64, height)

	for y := 0; y < height; y++ {
		chanBuffer.Add(1)
		go getFBMRow(
			&chanBuffer, rows,
			&noise,
			y, width, context, scale)
	}

	chanBuffer.Wait()
	close(rows)

	{
		for item := range rows {
			y := int(item[0])
			for x := 1; x < width+1; x++ {
				noiseMap[(x-1)+(y*width)] = item[x]
			}
		}
	}
	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0
	for i := 0; i < width*height; i++ {
		noiseHeight := noiseMap[i]
		if noiseHeight > maxNoiseHeight {
			maxNoiseHeight = noiseHeight
		} else if noiseHeight < minNoiseHeight {
			minNoiseHeight = noiseHeight
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

func GenPerturbedFBMParalell(width int, height int, context NoiseParams, scale float64, perturbAmount float64) *[]float64 {

	noise := snoise.NewNormalized(context.Seed)
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

	var chanBuffer sync.WaitGroup
	rows := make(chan []float64, height)

	for y := 0; y < height; y++ {
		chanBuffer.Add(1)
		go getFBMRowPerturbed(
			&chanBuffer, rows,
			&noise,
			y, width,
			context, perturbContext, scale, perturbAmount)
	}

	chanBuffer.Wait()
	close(rows)

	{
		for item := range rows {
			y := int(item[0])
			for x := 1; x < width+1; x++ {
				noiseMap[(x-1)+(y*width)] = item[x]
			}
		}
	}
	// Make elevationmap with floats
	maxNoiseHeight := 1.0
	minNoiseHeight := 0.0
	for i := 0; i < width*height; i++ {
		noiseHeight := noiseMap[i]
		if noiseHeight > maxNoiseHeight {
			maxNoiseHeight = noiseHeight
		} else if noiseHeight < minNoiseHeight {
			minNoiseHeight = noiseHeight
		}
	}

	for i := 0; i < len(noiseMap); i++ {
		noiseMap[i] /= maxNoiseHeight
		noiseMap[i] = GetBias(noiseMap[i], context.Bias)
		noiseMap[i] = GetGain(noiseMap[i], context.Gain)
	}

	return &noiseMap
}

// Small jobs helper functions
func getFBMRow(buffer *sync.WaitGroup, c chan []float64, pnoise *snoise.Noise, y, width int, context NoiseParams, scale float64) {
	defer buffer.Done()
	noise := *pnoise

	row := make([]float64, width+1)
	row[0] = float64(y)
	for x := 1; x < width+1; x++ {
		row[x] = getFBM(&noise, x, y, context, scale)
	}
	c <- row
}
func getFBMchunk(buffer *sync.WaitGroup, c chan []float64, pnoise *snoise.Noise, rowOffset, width, chunkSize int, context NoiseParams, scale float64) {
	defer buffer.Done()
	noise := *pnoise

	chunk := make([]float64, chunkSize+1)
	chunk[0] = float64(rowOffset)

	for i := 1; i < chunkSize+1; i++ {

		x := (i - 1) % width
		y := rowOffset + ((i - 1) / width)

		chunk[i] = getFBM(&noise, x, y, context, scale)
	}
	c <- chunk
}

func getFBMRowPerturbed(buffer *sync.WaitGroup, c chan []float64, pnoise *snoise.Noise,
	y, width int, context, perturbContext NoiseParams, scale, perturbAmount float64) {

	defer buffer.Done()
	noise := *pnoise

	row := make([]float64, width+1)
	row[0] = float64(y)
	for x := 1; x < width+1; x++ {
		//offsetx := getFBM(&noise, x-1, y, perturbContext, scale*perturbAmount)
		///offsety := getFBM(&noise, x-1, y, perturbContext, scale*perturbAmount)

		//row[x] = getFloat64FBM(&noise, float64(x)+offsetx, float64(y)+offsety, context, scale)
		perturbation := getFloat64FBM(&noise, float64(x), float64(y), context, scale)
		amplitude := 1.0
		frequency := 1.0
		noiseHeight := 0.0
		// get new value
		for i := 0; i < context.Octaves; i++ {
			sampleX := float64(x) / (context.ScaleX * scale) * (frequency * perturbation)
			sampleY := float64(y) / (context.ScaleY * scale) * (frequency * perturbation)

			nValue := noise.Eval2(sampleX, sampleY)
			noiseHeight += nValue * amplitude
			amplitude *= context.Persistance
			frequency *= context.Lacunarity
		}
		row[x] = noiseHeight
	}
	c <- row
}
