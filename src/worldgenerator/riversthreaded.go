package worldgenerator

// import (
// 	"common"
// 	"fmt"
// 	"time"
// )

// func generateWatershedsThreaded(world *World, parameters WorldParams, pheightMap, pMoistureMap *[]float64) []float64 {
// 	start := time.Now().UnixNano() / 1000000

// 	// Djikstra stuff
// 	//riverNum := 10
// 	shoreline := make([]int, 0)

// 	// Find shoreline positions
// 	for i := range world.DistanceToWater {
// 		if world.DistanceToLand[i] == 1 {
// 			shoreline = append(shoreline, i)
// 		}
// 	}

// 	width := parameters.Width
// 	height := parameters.Height
// 	size := height * width
// 	heightMap := *pheightMap
// 	moistureMap := *pMoistureMap
// 	distanceToWater := world.DistanceToWater
// 	fWaterHeight := parameters.FWaterHeight

// 	// Data used by all watersheds
// 	totalTakenTiles := make([]int, size)
// 	givesWaterTo := make([]int, size)    // a list of tiles
// 	getsWaterFrom := make([][]int, size) // what tiles will run down to this tile
// 	flowStrength := make([]float64, size)
// 	maxFlow := 0.0
// 	maxCost := 0.0

// 	numWatersheds := 300
// 	finishedWatersheds := 0
// 	interval := len(shoreline) / numWatersheds

// 	// Data local to one watershed
// 	frontier := make([]common.PriorityQueue, numWatersheds)
// 	costSoFar := make([]map[int]float64, numWatersheds)
// 	watershedSize := make([]int, numWatersheds)
// 	watershedEnds := make([][]int, numWatersheds)
// 	donzo := make([]bool, numWatersheds)

// 	//init
// 	for i := range givesWaterTo {
// 		givesWaterTo[i] = -1
// 	}

// 	for ws := 0; ws < numWatersheds; ws++ {
// 		watershedSize[ws] = 1
// 		frontier[ws] = common.NewPriorityQueue()
// 		costSoFar[ws] = make(map[int]float64)

// 		wsGlobalIdx := shoreline[ws*interval]
// 		frontier[ws].Insert(wsGlobalIdx, 0.01)
// 		totalTakenTiles[wsGlobalIdx] = wsGlobalIdx // Claim Mouth of watershed
// 		givesWaterTo[wsGlobalIdx] = 0
// 	}

// 	// Each loop represents one tick of watershed growth
// 	for finishedWatersheds < numWatersheds {
// 		// Each watershed
// 		for ws := 0; ws < numWatersheds; ws++ {
// 		}

// 	}

// 	for i := range givesWaterTo {
// 		if givesWaterTo[i] == -1 {
// 			givesWaterTo[i] = 0
// 		}
// 	}
// 	//flow through riversheds from all ends
// 	for watershedidx, ends := range watershedEnds {
// 		for i := range ends {

// 			current := givesWaterTo[watershedEnds[watershedidx][i]]
// 			accFlow := 0.0
// 			perTileFlow := 0.01

// 			for current != 0 {
// 				accFlow += perTileFlow
// 				flowStrength[current] += (accFlow + (moistureMap[current] * 0.6))

// 				if flowStrength[current] > maxFlow {
// 					maxFlow = flowStrength[current]
// 				}
// 				current = givesWaterTo[current]
// 			}
// 		}
// 	}

// 	// Figure out the flows of watersheds
// 	for i, _ := range flowStrength {
// 		flowStrength[i] /= maxFlow
// 		if flowStrength[i] > 0.001 {

// 			moistureMap[i] = (moistureMap[i] + (flowStrength[i] * 0.5))

// 			for _, n := range getIdxNeighbours(i, width) {
// 				if n >= size || n < 0 {
// 					continue
// 				}
// 				moistureMap[n] = (moistureMap[n] + (flowStrength[i] * 0.25))
// 				//if heightMap[n] > parameters.FWaterHeight+0.1 {
// 				///	heightMap[n] *= 1.0 - (flowStrength[i] * 0.25)
// 				//}
// 			}
// 		}
// 	}

// 	fmt.Println(fmt.Sprintf("Max cost: %.3f", maxCost))

// 	world.FlowStrength = flowStrength
// 	world.FlowsFrom = getsWaterFrom
// 	world.FlowsTo = givesWaterTo

// 	t := time.Now().UnixNano() / 1000000
// 	fmt.Println("Watershed msec: ", t-start)
// 	return moistureMap
// }

// func expandWatershed( ws int,
// 	heightMap*, distanceToWater []byte,
// 	watershedSize
// 	frontier,
// 	costSoFar,
// 	totalTakenTiles,
// 	donzo
// 	)
// 	{
// 		if donzo {
// 			continue
// 		}
// 		if frontier.Len() == 0 {
// 			donzo = true
// 			finishedWatersheds++
// 			continue
// 		}

// 		// Get Top priority tile
// 		topPrio, _ := frontier.Pop()
// 		currTile := topPrio.(int)

// 		if totalTakenTiles[currTile] != 0 && totalTakenTiles[currTile] != watershedIdx {
// 			continue
// 		}
// 		//Assign current tile as taken
// 		totalTakenTiles[currTile] = watershedIdx
// 		watershedSize[ws]++

// 		//Check all the neighbours
// 		currentIdxNeighbors := getIdxNeighbours(currTile, width)
// 		rejected := 0
// 		for _, next := range currentIdxNeighbors {

// 			newCost := estimateCost(currTile, next,
// 				&heightMap, &distanceToWater, &costSoFar[ws])

// 			oldCost, alreadyEstimated := costSoFar[ws][next]

// 			if !alreadyEstimated || newCost < oldCost {
// 				//check if already there
// 				alreadyThere := givesWaterTo[next]
// 				if alreadyThere != -1 || totalTakenTiles[next] != 0 || heightMap[next] > 0.7 || heightMap[next] < fWaterHeight {
// 					rejected++
// 					continue
// 				}

// 				if newCost > maxCost {
// 					maxCost = newCost
// 				}
// 				costSoFar[ws][next] = newCost
// 				frontier[ws].Insert(next, newCost)
// 				givesWaterTo[next] = currTile
// 				getsWaterFrom[currTile] = append(getsWaterFrom[currTile], next)
// 			}
// 		}
// 		// No suitable growth into neighbours
// 		if rejected >= 3 {
// 			watershedEnds = append(watershedEnds[ws], currTile)
// 		}
// 	}
// }
//
