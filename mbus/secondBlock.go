package mbus

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// Get the start of the second block
func (fr *Frame) startSecondBlock() int {
	if fr.Format() == "B" {
		return 11
	}

	// Assume block A
	return 13
}

// Parser for the second block
func (fr *Frame) ControlInformationField() string {
	if fr.Format() == "B" {
		return fr.getHexValue(11)
	}
	return fr.getHexValue(13)
}

const NoHeader = 0
const shortHeader = 1
const longHeader = 2

// Get information which data header is sent
// returns the const defined above
func (fr *Frame) DataHeader() int {
	switch fr.ControlInformationField() {
	case "78":
		return NoHeader
	case "7A":
		return shortHeader
	case "72":
		return longHeader
	}

	// Hack: Assume that everything else is invalid
	return 0
}

// TODO: Support for encrypted devices

// Get the access number
func (fr *Frame) AccessNumber() string {
	switch fr.DataHeader() {
	case NoHeader:
		return ""
	case longHeader:
		return fr.getHexValue(fr.startSecondBlock() + 9)
	}

	// Hack:Assume short header
	return fr.getHexValue(fr.startSecondBlock() + 1)
}

// Get the status field
func (fr *Frame) StatusField() string {
	switch fr.DataHeader() {
	case NoHeader:
		return ""
	case shortHeader:
		return fr.getHexValue(fr.startSecondBlock() + 2)
	}

	// Hack: Assume long header
	return fr.getHexValue(fr.startSecondBlock() + 10)
}

// Get the configuration field
func (fr *Frame) Configuration() string {
	switch fr.DataHeader() {
	case NoHeader:
		return ""
	case shortHeader:
		return fr.getHexValue(fr.startSecondBlock()+3, fr.startSecondBlock()+4)
	}

	// Hack: Assume long header
	return fr.getHexValue(fr.startSecondBlock()+11, fr.startSecondBlock()+12)
}

// Get the configuration length
func (fr *Frame) ConfigurationLength() int {
	decodedString, _ := hex.DecodeString("30")
	value, _ := strconv.Atoi(fmt.Sprint(decodedString[0]))
	return value
}

// Starting point of the encrypted value
func (fr *Frame) ValueStart() int {
	switch fr.DataHeader() {
	case NoHeader:
		return 0
	case shortHeader:
		return fr.startSecondBlock() + 5
	}

	// Hack: Assume long header
	return fr.startSecondBlock() + 12
}

// Additional Info
func (fr *Frame) IV() string {
	return fr.getHexValue(3) + fr.getHexValue(4) + fr.Address() + fr.AccessNumber() + fr.AccessNumber() + fr.AccessNumber() + fr.AccessNumber() + fr.AccessNumber() + fr.AccessNumber() + fr.AccessNumber() + fr.AccessNumber()
}
