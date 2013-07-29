package mbus

import (
	"encoding/json"
	"strings"
	"time"
	"fmt"
)

// TODO: Load into the frame
type Frame struct {
	Value    string
	Hexified string
	Time     time.Time
	ID       int
	Key      string
}

// Our socket handler for new sniffed data
func (fr Frame) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ID":                      fr.ID,
		"Value":                   fr.Value,
		"Length":                  fr.Length(),
		"Control":                 fr.Control(),
		"Manufacturer":            fr.Manufacturer(),
		"ManufacturerString":      fr.ManufacturerString(),
		"Address":                 fr.Address(),
		"Identification":          fr.Identification(),
		"Version":                 fr.Version(),
		"Hexified":                fr.Hexified,
		"DeviceType":              fr.DeviceTypeString(),
		"ControlInformationField": fr.ControlInformationField(),
		"AccessNumber":            fr.AccessNumber(),
		"StatusField":             fr.StatusField(),
		"Configuration":           fr.Configuration(),
		"IV":                      fr.IV(),
		"Time":                    fr.Time,
	})
}

// Frame creates and initializes a new frame
func NewFrame(fr string) (*Frame, error) {
	frame := &Frame{Value: fmt.Sprintf("%s", strings.ToUpper(fmt.Sprintf("%s",fr)))}

	return frame, nil
}

// TODO: Check whether a frame is valid
func (fr *Frame) Valid() bool {
	return false
}

// Type A / Type B
func (fr *Frame) Format() string {
	// Check whether field 11 is the CI
	// If field 11 a CI field it's frame format B
	ciFieldFormatB := fr.getHexValue(11)
	if ciFieldFormatB == "51" || ciFieldFormatB == "71" || ciFieldFormatB == "72" || ciFieldFormatB == "78" || ciFieldFormatB == "7A" || ciFieldFormatB == "81" {
		return "B"
	}

	// Quick'n dirty: Assume A for everything else
	return "A"
}
