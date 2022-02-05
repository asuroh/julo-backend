package str

import (
	"fmt"
	"strconv"
)

// Contains ...
func Contains(slices []string, comparizon string) bool {
	for _, a := range slices {
		if a == comparizon {
			return true
		}
	}

	return false
}

// StringToBool ...
func StringToBool(data string) bool {
	res, err := strconv.ParseBool(data)
	if err != nil {
		res = false
	}

	return res
}

// StringToInt ...
func StringToInt(data string) int {
	res, err := strconv.Atoi(data)
	if err != nil {
		res = 0
	}

	return res
}

// StringToFloat ...
func StringToFloat(data string) float64 {
	res, err := strconv.ParseFloat(data, 64)
	if err != nil {
		res = 0
	}

	return res
}

// Float64ToString ...
func Float64ToString(data float64) string {
	return fmt.Sprintf("%g", data)
}

// IntToString ...
func IntToString(data int) string {
	res := strconv.Itoa(data)
	return res
}

// DefaultData ...
func DefaultData(data, defaultData string) string {
	if data == "" {
		return defaultData
	}

	return data
}

// ShowString ...
func ShowString(isShow bool, data string) string {
	if isShow {
		return data
	}

	return ""
}
