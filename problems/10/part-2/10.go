package problems_10_2

// #cgo LDFLAGS: -lgeos_c
// #include <geos_c.h>

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	common_functions "aoc.2023/lib/common/functions"
	common_types "aoc.2023/lib/common/types"
)

// Map all variables as Enums
var ANIMAL_FLAG = []byte("S")[0]

const (
	DIRECTION_RIGHT_ENUM  = 0
	DIRECTION_BOTTOM_ENUM = 1
	DIRECTION_LEFT_ENUM   = 2
	DIRECTION_UP_ENUM     = 3
	DIRECTION_UNSET_ENUM  = -1
)

var (
	VERTICAL_PIPE_ENUM   = []byte("|")[0]
	HORIZONTAL_PIPE_ENUM = []byte("-")[0]
	NORTH_EAST_PIPE_ENUM = []byte("L")[0]
	NORTH_WEST_PIPE_ENUM = []byte("J")[0]
	SOUTH_WEST_PIPE_ENUM = []byte("7")[0]
	SOUTH_EAST_PIPE_ENUM = []byte("F")[0]
	GROUND_PIPE_ENUM     = []byte(".")[0]
)

func SolveChallenge(problemId string) string {
	// Process the input
	inputFilePath := fmt.Sprintf("problems/%s/input.txt", problemId)
	scanner := common_functions.CreateInputScanner(inputFilePath)
	defer scanner.File.Close()

	var answer float64 = 0

	pipeNetworkMap := getPipeNetworkMap(*scanner)

	polygonArea := getPolygonArea(pipeNetworkMap)
	b := float64(len(pipeNetworkMap))

	answer = (polygonArea - (b / 2)) + 1

	return strconv.FormatFloat(answer, 'f', -1, 64)
}

func getPolygonArea(polygon [][2]int) float64 {
	var (
		xProd int = 0
		yProd int = 0
	)

	for i := 1; i < len(polygon); i++ {
		xProd += (polygon[i][0] * polygon[i-1][1])
		yProd += (polygon[i][1] * polygon[i-1][0])
	}

	var absSum = math.Abs(float64(xProd) - float64(yProd))

	return math.Floor(absSum / 2)
}

func getPipeNetworkMap(scanner common_types.FileInputScanner) [][2]int {
	// Get the map and position of the animal to start trip
	animalPos, pipeNetwork := processPipeNetworkInput(scanner)

	// At the beginning the animal has only 4 directions to travel...
	var pipeNetworkMap [][2]int

	pipeNetworkMap = append(pipeNetworkMap, animalPos)

	var animalPositions = [4]int{
		DIRECTION_RIGHT_ENUM,
		DIRECTION_BOTTOM_ENUM,
		DIRECTION_LEFT_ENUM,
		DIRECTION_UP_ENUM,
	}

	for _, animalDir := range animalPositions {

		var (
			pathResult      = animalPos
			directionResult = animalDir
			minPosIndex     = 0
		)

		// Infinite loop that breaks when we can't travel in the pipe network
		for {
			// Recalculate the new step
			pathResult, directionResult = computePath(pathResult, pipeNetwork, directionResult)

			// If we reached a limit in the map (out of boundaries limit)
			// Or there isn't a valid connected pipe
			// Or there is a ground
			// Or if there is the animal again! (infinite loop)
			if directionResult == DIRECTION_UNSET_ENUM {

				// We make sure if it's an infinite path,
				// then is necessary divide by 2 to get the greatest path
				posTmp := pipeNetwork[pathResult[0]][pathResult[1]]
				if posTmp == ANIMAL_FLAG {

					fmt.Println(pipeNetworkMap[minPosIndex])
					pipeNetworkMapTmp := pipeNetworkMap[minPosIndex:]
					pipeNetworkMapTmp = append(pipeNetworkMapTmp, pipeNetworkMap[0:minPosIndex]...)

					return pipeNetworkMapTmp
				}

				break
			}

			pipeNetworkMap = append(pipeNetworkMap, [2]int{pathResult[0], pathResult[1]})

			if pathResult[0] <= pipeNetworkMap[minPosIndex][0] && pathResult[1] < pipeNetworkMap[minPosIndex][1] {
				minPosIndex = len(pipeNetworkMap) - 1
			}
		}
	}

	panic("The pipe network is NOT a loop!")
}

func processPipeNetworkInput(scanner common_types.FileInputScanner) ([2]int, []string) {
	var (
		pipeNetwork   []string
		animalPos     [2]int
		animalFlagStr = string(ANIMAL_FLAG)
	)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		pipeNetwork = append(pipeNetwork, line)

		// Detect animal position
		animalColPos := strings.Index(line, animalFlagStr)
		if animalColPos > -1 {
			animalPos = [2]int{i, animalColPos}
		}
	}

	return animalPos, pipeNetwork
}

func computePath(path [2]int, pipeNetwork []string, pipeDirection int) ([2]int, int) {
	// Direction is to know how to move through the map (columns and rows)
	switch pipeDirection {
	case DIRECTION_RIGHT_ENUM:
		path[1]++

		// Out of boundaries
		if path[1] >= len(pipeNetwork[1]) {
			return path, DIRECTION_UNSET_ENUM
		}

		newDirection := validateRightDirection(pipeNetwork[path[0]][path[1]])
		return path, newDirection
	case DIRECTION_BOTTOM_ENUM:
		path[0]++

		if path[0] >= len(pipeNetwork) {
			return path, DIRECTION_UNSET_ENUM
		}

		newDirection := validateBottomDirection(pipeNetwork[path[0]][path[1]])
		return path, newDirection
	case DIRECTION_LEFT_ENUM:
		path[1]--

		if path[1] < 0 {
			return path, DIRECTION_UNSET_ENUM
		}

		newDirection := validateLeftDirection(pipeNetwork[path[0]][path[1]])
		return path, newDirection
	case DIRECTION_UP_ENUM:
		path[0]--

		if path[0] < 0 {
			return path, DIRECTION_UNSET_ENUM
		}

		newDirection := validateUpDirection(pipeNetwork[path[0]][path[1]])
		return path, newDirection
	default:
		panic("Bad pipe direction")
	}
}

// These 4 methods are to calculate the new direction according to the pipe
func validateRightDirection(pipe byte) int {
	switch pipe {
	case HORIZONTAL_PIPE_ENUM:
		return DIRECTION_RIGHT_ENUM

	case NORTH_WEST_PIPE_ENUM:
		return DIRECTION_UP_ENUM

	case SOUTH_WEST_PIPE_ENUM:
		return DIRECTION_BOTTOM_ENUM

	default:
		return DIRECTION_UNSET_ENUM
	}
}

func validateBottomDirection(pipe byte) int {
	switch pipe {
	case VERTICAL_PIPE_ENUM:
		return DIRECTION_BOTTOM_ENUM

	case NORTH_EAST_PIPE_ENUM:
		return DIRECTION_RIGHT_ENUM

	case NORTH_WEST_PIPE_ENUM:
		return DIRECTION_LEFT_ENUM

	default:
		return DIRECTION_UNSET_ENUM
	}
}

func validateLeftDirection(pipe byte) int {
	switch pipe {
	case HORIZONTAL_PIPE_ENUM:
		return DIRECTION_LEFT_ENUM

	case NORTH_EAST_PIPE_ENUM:
		return DIRECTION_UP_ENUM

	case SOUTH_EAST_PIPE_ENUM:
		return DIRECTION_BOTTOM_ENUM

	default:
		return DIRECTION_UNSET_ENUM
	}
}

func validateUpDirection(pipe byte) int {
	switch pipe {
	case VERTICAL_PIPE_ENUM:
		return DIRECTION_UP_ENUM

	case SOUTH_EAST_PIPE_ENUM:
		return DIRECTION_RIGHT_ENUM

	case SOUTH_WEST_PIPE_ENUM:
		return DIRECTION_LEFT_ENUM

	default:
		return DIRECTION_UNSET_ENUM
	}
}
