package oftp2

import (
	"fmt"
	"strconv"
)

func isBool(input string) bool {
	return input == "Y" || input == "N"
}

func fillUpString(in string, desiredSize int) (string, error) {
	lenIn := len(in)
	if lenIn > desiredSize {
		return in, fmt.Errorf("exceeded capacity: %s (%d)", in, desiredSize)
	}
	for i := lenIn; i < desiredSize; i++ {
		in = " " + in
	}
	return in, nil
}

func fillUpInt(in int, desiredSize int) (string, error) {
	result := strconv.Itoa(in)
	lenIn := len(result)
	if lenIn > desiredSize {
		return result, fmt.Errorf("exceeded capacity: %d (%d)", in, desiredSize)
	}
	for i := lenIn; i < desiredSize; i++ {
		result = "0" + result
	}
	return result, nil
}

func boolToString(input bool) string {
	if input {
		return "Y"
	}
	return "N"
}
