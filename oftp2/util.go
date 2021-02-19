package oftp2

import (
	"fmt"
	"strconv"
)

func isBool(input string) bool {
	return input == "Y" || input == "N"
}

func fillUpString(in string, desiredSize int) (string, error) {
	if len(in) > desiredSize {
		return in, fmt.Errorf("exceeded capacity: %s (%d)", in, desiredSize)
	}
	return fmt.Sprintf("%"+strconv.Itoa(desiredSize)+"s", in), nil
}

func fillUpInt(in int, desiredSize int) (string, error) {
	if result := strconv.Itoa(in); len(result) > desiredSize {
		return result, fmt.Errorf("exceeded capacity: %d (%d)", in, desiredSize)
	}
	return fmt.Sprintf("%0"+strconv.Itoa(desiredSize)+"d", in), nil
}

func boolToString(input bool) string {
	if input {
		return "Y"
	}
	return "N"
}
