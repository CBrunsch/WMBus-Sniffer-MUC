package mbus

import (
	"strconv"
)

// Parser for the first block 
// The length of the frame in byte (exclusive length byte)
func (fr *Frame) Length() int {
	return fr.hexToInt(1)
}

// Control type
func (fr *Frame) Control() int {
	control, _ := strconv.Atoi(fr.getHexValue(2))
	return control
}

// Get the vendor of a specific device
func (fr *Frame) Manufacturer() string {
	return fr.getHexValue(4) + fr.getHexValue(3)
}

// TODO: Get the manufacturer as string
func (fr *Frame) ManufacturerString() string {
	manString := fr.getHexValue(3) + fr.getHexValue(4)
	switch manString {
	case "A205":
		manString = "AMB"
	case "2D2C":
		manString = "KAM"
	case "7916":
		manString = "ESY"
	case "9226":
		manString = "ITG"
	}
	return manString
}

// Get the adress of the device
func (fr *Frame) Address() string {
	return fr.getHexValue(5, 10)
}

// Get the identification of a device
func (fr *Frame) Identification() string {
	return fr.getHexValue(8) + " " + fr.getHexValue(7) + " " + fr.getHexValue(6) + " " + fr.getHexValue(5)
}

// Get the version of a specific device
func (fr *Frame) Version() int {
	version, _ := strconv.Atoi(fr.getHexValue(9))
	return version
}

// Get the device type of a specific device
func (fr *Frame) DeviceType() string {
	return fr.getHexValue(10)
}

// Get the device type of a device as string
func (fr *Frame) DeviceTypeString() string {
	var deviceTypes = map[string]string{
		"00": "Other",
		"01": "Oil",
		"02": "Electricity",
		"03": "Gas",
		"04": "Heat",
		"05": "Steam",
		"06": "Warm Water",
		"07": "Water",
		"08": "Heat Cost Allocator",
		"09": "Compressed Air",
		"0A": "Cooling load meter",
		"0B": "Cooling load meter",
		"0C": "Heat (Volume measured at flow temperature)",
		"0D": "Heat /Cooling load meter",
		"0E": "Bus / System component",
		"0F": "Unknown medium",
		"15": "Hot Water",
		"16": "Cold water",
		"17": "Dual register (hot/cold) Water meter",
		"18": "Pressure",
		"19": "A/D Converter",
		"37": "Radio converter",
	}
	deviceType := deviceTypes[fr.DeviceType()]
	if deviceType == "" { // Fallback in case no device type has been recognized
		return fr.DeviceType()
	}
	return deviceType
}
