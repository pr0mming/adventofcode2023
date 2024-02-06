package problems_9_2

import (
	"fmt"
	"strconv"
	"strings"

	common_functions "aoc.2023/lib/common/functions"
)

func SolveChallenge(problemId string) string {
	// Process the input
	inputFilePath := fmt.Sprintf("problems/%s/input.txt", problemId)
	scanner := common_functions.CreateInputScanner(inputFilePath)
	defer scanner.File.Close()

	answer := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Extract records (array) from the input
		historyRecord := common_functions.GetIntegersArr(strings.Fields(line), false)

		answer += computeHistoryRecord(historyRecord)
	}

	return strconv.Itoa(answer)
}

func computeHistoryRecord(historyRecord []int) int {
	var zeros int                            // Flag to check if we have all the values of the history record are zeros
	var offsetIndex = len(historyRecord) - 1 // Each iteration the history records is len - 1, we use it to control the new record and access the last item
	var sequence []int                       // Keep all the first numbers for each new history record

	for zeros < len(historyRecord) {
		// New empty copy of the history record
		var newHistoryRecord = make([]int, offsetIndex)

		// App up last item
		sequence = append(sequence, historyRecord[0])

		// Fill the array with differences
		for i := 1; i < len(historyRecord); i++ {
			diff := historyRecord[i] - historyRecord[i-1]
			newHistoryRecord[i-1] = diff

			if diff == 0 {
				zeros++
			}
		}

		historyRecord = newHistoryRecord
		offsetIndex-- // The next array will be fewer one
	}

	var sum = 0

	for i := len(sequence) - 1; i >= 0; i-- {
		sum = sequence[i] - sum
	}

	return sum
}
