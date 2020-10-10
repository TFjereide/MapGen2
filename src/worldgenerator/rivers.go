package worldgenerator

import (
	"common"
	"fmt"
	"math/rand"
	"time"
)

func generateWatersheds(world *World, parameters WorldParams, pheightMap, pMoistureMap *[]float64) []float64 {
	start := time.Now().UnixNano() / 1000000

	// Djikstra stuff
	//riverNum := 10
	shoreline := make([]int, 0)

	// Find shoreline positions
	for i := range world.DistanceToWater {
		if world.DistanceToLand[i] == 1 {
			shoreline = append(shoreline, i)
		}
	}

	rand.Shuffle(len(shoreline), func(i, j int) {
		shoreline[i], shoreline[j] = shoreline[j], shoreline[i]
	})

	width := parameters.Width
	height := parameters.Height
	size := height * width
	heightMap := *pheightMap
	moistureMap := *pMoistureMap
	distanceToWater := world.DistanceToWater
	fWaterHeight := parameters.FWaterHeight

	// Data used by all watersheds
	totalTakenTiles := make([]int, size)
	givesWaterTo := make([]int, size)    // a list of tiles
	getsWaterFrom := make([][]int, size) // what tiles will run down to this tile
	flowStrength := make([]float64, size)
	maxFlow := 0.0
	maxCost := 0.0
	maxPossibleCost := 10000.0

	numWatersheds := 600
	finishedWatersheds := 0
	interval := len(shoreline) / numWatersheds

	// Data local to one watershed
	frontier := make([]common.PriorityQueue, numWatersheds)
	costSoFar := make([]map[int]float64, numWatersheds)
	watershedSize := make([]int, numWatersheds)
	watershedEnds := make([][]int, numWatersheds)
	donzo := make([]bool, numWatersheds)

	//init
	for i := range givesWaterTo {
		givesWaterTo[i] = -1
	}

	for ws := 0; ws < numWatersheds; ws++ {
		watershedSize[ws] = 1
		frontier[ws] = common.NewPriorityQueue()
		costSoFar[ws] = make(map[int]float64)

		wsGlobalIdx := shoreline[ws*interval]
		frontier[ws].Insert(wsGlobalIdx, 0.01)
		totalTakenTiles[wsGlobalIdx] = wsGlobalIdx // Claim Mouth of watershed
		givesWaterTo[wsGlobalIdx] = 0
	}

	// Each loop represents one tick of watershed growth
	for finishedWatersheds < numWatersheds {
		// Each watershed
		for ws := 0; ws < numWatersheds; ws++ {

			watershedIdx := shoreline[ws*interval]

			if donzo[ws] {
				continue
			}
			if frontier[ws].Len() == 0 {
				donzo[ws] = true
				finishedWatersheds++
				continue
			}

			// Get Top priority tile
			topPrio, _ := frontier[ws].Pop()
			currTile := topPrio.(int)

			if totalTakenTiles[currTile] != 0 && totalTakenTiles[currTile] != watershedIdx {
				continue
			}
			//Assign current tile as taken
			totalTakenTiles[currTile] = watershedIdx
			watershedSize[ws]++

			//Check all the neighbours
			currentIdxNeighbors := getIdxNeighbours(currTile, width)
			rejected := 0
			for _, next := range currentIdxNeighbors {

				newCost := estimateCost(currTile, next,
					&heightMap, &distanceToWater, &costSoFar[ws])

				oldCost, alreadyEstimated := costSoFar[ws][next]

				if !alreadyEstimated || newCost < oldCost {
					//check if already there
					alreadyThere := givesWaterTo[next]
					if alreadyThere != -1 || totalTakenTiles[next] != 0 || heightMap[next] > 0.7 ||
						heightMap[next] < fWaterHeight || newCost > maxPossibleCost {
						rejected++
						continue
					}

					if newCost > maxCost {
						maxCost = newCost
					}
					costSoFar[ws][next] = newCost
					frontier[ws].Insert(next, newCost)
					givesWaterTo[next] = currTile
					getsWaterFrom[currTile] = append(getsWaterFrom[currTile], next)
				}
			}
			// No suitable growth into neighbours
			if rejected >= 3 {
				watershedEnds[ws] = append(watershedEnds[ws], currTile)
			}
		}

	}
	for i := range givesWaterTo {
		if givesWaterTo[i] == -1 {
			givesWaterTo[i] = 0
		}
	}

	//flow through riversheds from all ends
	for watershedidx, ends := range watershedEnds {
		for i := range ends {

			current := givesWaterTo[watershedEnds[watershedidx][i]]
			accFlow := 0.0
			perTileFlow := 0.01

			for current != 0 {
				accFlow += perTileFlow
				flowStrength[current] += (accFlow + (moistureMap[current] * 0.6))

				if flowStrength[current] > maxFlow {
					maxFlow = flowStrength[current]
				}
				current = givesWaterTo[current]
			}
		}
	}

	// Figure out the flows of watersheds
	for i, _ := range flowStrength {
		flowStrength[i] /= maxFlow
		if flowStrength[i] > 0.001 {

			moistureMap[i] = (moistureMap[i] + (flowStrength[i] * 0.5))

			for _, n := range getIdxNeighbours(i, width) {
				if n >= size || n < 0 {
					continue
				}
				moistureMap[n] = (moistureMap[n] + (flowStrength[i] * 0.25))
				if heightMap[n] > parameters.FWaterHeight+0.1 {
					heightMap[n] *= 1.0 - (flowStrength[i] * 0.25)
				}
			}
		}
	}

	fmt.Println(fmt.Sprintf("Max cost: %.3f", maxCost))

	world.FlowStrength = flowStrength
	world.FlowsFrom = getsWaterFrom
	world.FlowsTo = givesWaterTo

	t := time.Now().UnixNano() / 1000000
	fmt.Println("Watershed msec: ", t-start)
	return moistureMap

}

func generateWatershedsSimple(world *World, parameters WorldParams, pheightMap, pMoistureMap *[]float64) []float64 {
	start := time.Now().UnixNano() / 1000000

	shoreline := make([]int, 0)
	// Find shoreline positions
	for i := range world.DistanceToWater {
		if world.DistanceToLand[i] == 1 {
			shoreline = append(shoreline, i)
		}
	}

	width := parameters.Width
	height := parameters.Height
	size := height * width
	heightMap := *pheightMap
	moistureMap := *pMoistureMap
	distanceToWater := world.DistanceToWater
	fWaterHeight := parameters.FWaterHeight

	numWatersheds := 200

	// Create watershed tree.  Each tile can only belong to one tree
	totalTakenTiles := make([]int, size)
	givesWaterTo := make([]int, size)    // a list of tiles
	getsWaterFrom := make([][]int, size) // what tiles will run down to this tile
	flowStrength := make([]float64, size)

	maxFlow := 0.0
	maxCost := 0.0

	interval := len(shoreline) / numWatersheds
	for i := 0; i < len(shoreline); i += interval {

		watershedIdx := shoreline[i]

		if totalTakenTiles[watershedIdx] != 0 {
			continue
		}

		totalTakenTiles[watershedIdx] = watershedIdx

		//Start a Djikstra search from shoreline
		frontier := common.NewPriorityQueue()
		frontier.Insert(watershedIdx, 0.0)
		cameFrom := make(map[int]int)
		cameFrom[watershedIdx] = 0
		costSoFar := make(map[int]float64)
		costSoFar[watershedIdx] = 0

		watershedEnds := make([]int, 0)

		//Max size
		watershedSize := 1
		// Create graph from shorelines
		for frontier.Len() != 0 {

			topPrio, _ := frontier.Pop()
			currTile := topPrio.(int)

			//Assign current tile as taken by a watershed
			totalTakenTiles[currTile] = watershedIdx

			watershedSize++

			//Check all the neighbours
			currentIdxNeighbors := getIdxNeighbours(currTile, width)
			rejected := 0
			for _, next := range currentIdxNeighbors {

				newCost := estimateCost(currTile, next, &heightMap, &distanceToWater, &costSoFar)

				_, alreadyEstimated := costSoFar[next]

				if !alreadyEstimated || newCost < costSoFar[next] {
					//check if already there
					_, alreadyThere := cameFrom[next]
					if alreadyThere || totalTakenTiles[next] != 0 || heightMap[next] > 0.7 || heightMap[next] < fWaterHeight || watershedSize >= 1200 || newCost > 10.0 {
						rejected++
						continue
					}
					costSoFar[next] = newCost
					if newCost > maxCost {
						maxCost = newCost
					}
					priority := newCost
					frontier.Insert(next, priority)
					cameFrom[next] = currTile
					getsWaterFrom[currTile] = append(getsWaterFrom[currTile], next)
				}
			}
			if rejected >= 3 {
				watershedEnds = append(watershedEnds, currTile)
			}

		}

		//find flows to
		for key, val := range cameFrom {
			givesWaterTo[key] = val
		}

		//flow through riversheds from all ends
		for i, _ := range watershedEnds {

			current := givesWaterTo[watershedEnds[i]]
			accFlow := 0.0
			perTileFlow := 0.01

			for current != 0 {
				accFlow += perTileFlow
				flowStrength[current] += (accFlow + (moistureMap[current] * 0.6))

				if flowStrength[current] > maxFlow {
					maxFlow = flowStrength[current]
				}
				current = givesWaterTo[current]
			}
		}

	}
	for i, _ := range flowStrength {
		flowStrength[i] /= maxFlow
		if flowStrength[i] > 0.001 {

			moistureMap[i] = (moistureMap[i] + (flowStrength[i] * 0.5))

			for _, n := range getIdxNeighbours(i, width) {
				if n >= size || n < 0 {
					continue
				}
				moistureMap[n] = (moistureMap[n] + (flowStrength[i] * 0.25))
				//if heightMap[n] > parameters.FWaterHeight+0.1 {
				///	heightMap[n] *= 1.0 - (flowStrength[i] * 0.25)
				//}
			}
		}
	}

	fmt.Println(fmt.Sprintf("Max cost: %.3f", maxCost))

	world.FlowStrength = flowStrength
	world.FlowsFrom = getsWaterFrom
	world.FlowsTo = givesWaterTo

	t := time.Now().UnixNano() / 1000000
	fmt.Println("Watershed msec: ", t-start)
	return moistureMap
}

func estimateCost(current, next int, heightMap *[]float64, distanceToWater *[]byte, costSoFar *map[int]float64) float64 {
	_heightMap := *heightMap
	_distanceToWater := *distanceToWater
	_costSoFar := *costSoFar

	elevationCost := (_heightMap[current] - _heightMap[next]) // if next is higher the cost is negative
	if elevationCost > 0.0 {
		elevationCost *= 100.0

	}
	newCost := _costSoFar[current] + elevationCost
	if _distanceToWater[next] < 2 {
		newCost += 1.0
	}

	return newCost
}
