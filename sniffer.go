package main

import (
	"./mbus"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ziutek/serial"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Send data to the serial sending device, this function takes care
// of the binary encoding of the passed data.
// The data should be in format FF FF FF FF
func sendDataToSniffer(data string) {
	// Initialize the sender port
	var err error
	senderPort, err := serial.Open(*sniffingTTY)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
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

	time.Sleep(time.Millisecond * 100)
}

// Reads continous the data from the connected serial device
func readData() {
	s, err := serial.Open(*sniffingTTY)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	snifferActive = true

	var buf []byte

	//	Reset the device
	sendDataToSniffer("FF 05 00")

	for {
		newData := make([]byte, 1000000)
		n, err := s.Read(newData)
		if n > 0 {
			// Sprintf is used here to prevent interpreting the bytes (e.g. 0x13 as newline character)
			buf = append(buf, fmt.Sprintf("%s", hex.EncodeToString(newData[:n]))...)
			if err != nil {
				log.Fatal(err)
			}

			if bytes.Count(buf, []byte("ff03")) > 0 && snifferActive {
				frames := strings.Split(fmt.Sprintf("%s", buf), "ff03")

				for i := range frames {
					// Only proceed if we received the length byte as we will only read until there
					if len(frames[i]) > 2 {
						decodedString, err := hex.DecodeString(fmt.Sprintf("%s", frames[i][:2]))

						if err == nil {
							// The sniffer appends +1 to the length byte, therefore we have to decrement the value
							lengthByte, err := strconv.Atoi(fmt.Sprintf("%d", decodedString[0]))
							if err == nil {
								lengthByte = lengthByte - 1
								newLength := fmt.Sprintf("%#x", lengthByte)

								// Only read if the length is correct
								if len(frames[i]) >= (lengthByte*2)+2 && lengthByte > 10 {
									frame, _ := mbus.NewFrame(fmt.Sprintf("%s", strings.ToUpper(newLength[2:]+frames[i][2:(lengthByte*2)+2])))
									// HACK: Only sniff encrypted frames
									if frame.Configuration() != "0000" && frame.Configuration() != "" {
										go addFrameToDB(frame)
									}

									buf = []byte(strings.Join(frames[i+1:], ""))
								}
							}
						}

					}
				}

				// Just drop everything if the sniffer is not active
				if !snifferActive {
					buf = []byte("")
				}
			}
		}
	}
}

// Contains the current sniffer state (active/stopped)
var snifferActive bool

// Disable or enable the sniffer
func statusSnifferHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("status") == "stop" {
		snifferActive = false
	} else {
		snifferActive = true
	}
}
