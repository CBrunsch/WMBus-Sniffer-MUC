package mbus

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// Converts a specified value to string
func (fr *Frame) hexToString(first int, last ...int) string {
	decodedString, _ := hex.DecodeString(fr.getHexValue(first))
	return fmt.Sprint(decodedString[0])
}

// Converts a specified value to string
func (fr *Frame) hexToInt(first int, last ...int) int {
	value, _ := strconv.Atoi(fr.hexToString(first))
	return value
}

// Get the hex value at the specified location or the specified range (optional)
func (fr *Frame) getHexValue(first int, last ...int) string {
	if len(last) == 1 {
		return fr.Value[(first*2)-2 : (last[0] * 2)]
	}
	return fr.Value[(first*2)-2 : (first * 2)]
}

// Get the hex value at the specified location
func getHexValue(data string, first int, last ...int) string {
	if len(last) == 1 {
		return data[(first*2)-2 : (last[0] * 2)]
	}
	return data[(first*2)-2 : (first * 2)]
}
