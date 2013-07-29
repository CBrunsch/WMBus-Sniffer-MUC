package main

import (
	"./mbus"
	"fmt"
	"github.com/tarm/goserial"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Send data to the serial sending device, this function takes care
// of the binary encoding of the passed data.
// The data should be in format FF FF FF FF
func sendData(data string) {

	// Add it directly to the DB in case the sender is not working correctly
	// as the firmware is somewhat error prone
	if *DemoMode {
		if strings.HasPrefix(data, "FF 00") {
			newData := strings.Replace(data, " ", "", -1)
			frame, _ := mbus.NewFrame(fmt.Sprintf("%s", newData[4:]))
			go addFrameToDB(frame)
		}
	}

	// Initialize the sender config
	senderC := new(serial.Config)
	senderC.Name = *sendingTTY
	senderC.Baud = 9600

	// Initialize the sender port
	var err error
	senderPort, err := serial.OpenPort(senderC)
	if err != nil {
		log.Fatal(err)
	}

	strippedData := strings.Split(data, " ")
	var checkSum uint64
	hexLiteral := []byte{}
	for index := range strippedData {
		actualByte, _ := strconv.ParseUint(strippedData[index], 16, 10)
		hexLiteral = append(hexLiteral, byte(actualByte))
		checkSum = checkSum ^ actualByte
	}

	// Append checksum
	hexLiteral = append(hexLiteral, byte(checkSum))
	senderPort.Write(hexLiteral)
	senderPort.Close()

	time.Sleep(time.Millisecond * 250)
}

// Initializes the sending device (Reset + T2 mode)
func initializeSender() {
	//	Reset the device
	sendData("FF 05 00")
	time.Sleep(time.Second * 3)

	// Set device mode to T2
	sendData("FF 09 03 46 01 07")
	time.Sleep(time.Second * 3)

	// Set device mode to T2
	sendData("FF 09 03 46 01 07")
	time.Sleep(time.Second * 3)
}

// Sanitizes a frame
func sanitizeFrame(frame string) string {
	if len(frame) > 0 {
		// Remove double spaces
		frame = strings.Replace(frame, "  ", " ", -1)

		splittedFrame := strings.Split(frame, "")
		containsSpace := strings.Contains(frame, " ")
		if !containsSpace {
			var endString string
			var i int
			for i = 0; i < len(splittedFrame); {
				endString = endString + splittedFrame[i] + splittedFrame[i+1] + " "
				i = i + 2
			}
			frame = endString[:len(endString)-1]
		}
	}
	return strings.ToUpper(frame)
}

// Handler for sending frames
func sendFrameHandler(w http.ResponseWriter, r *http.Request) {
	praemble := "FF 00" // Command + praemble
	frame := sanitizeFrame(r.FormValue("frame"))
	frame2 := sanitizeFrame(r.FormValue("frame2"))
	numberString := r.FormValue("number")
	var numberInt int
	if numberString != "" {
		numberInt, _ = strconv.Atoi(numberString)
	} else {
		numberInt = 1
	}

	for i := 0; i < numberInt; i++ {
		sendData(praemble + " " + frame)
		if frame2 != "" {
			sendData(praemble + " " + frame2)
		}
	}
}

// Handler for sending data partially encrypted
func sendPartiallyHandler(w http.ResponseWriter, r *http.Request) {
	praemble := "FF 00" // Command + praemble
	frame := sanitizeFrame(r.FormValue("frame"))
	appendedData := sanitizeFrame(r.FormValue("appendedData"))
	frameWithoutLength := frame[2:]
	newLength := strings.ToUpper(fmt.Sprintf("%#x", (len(frameWithoutLength)+len(appendedData)+1)/3))
	sendData(praemble + " " + newLength[2:] + frameWithoutLength + " " + appendedData)
}
