package mbus

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"strings"
)

type DIFs struct {
	DIF      string
	VIF      string
	Value    string
	Original string
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Get the plain text out of a frame
func (fr *Frame) PlainText() []DIFs {
	if fr.Key != "" {
		// Check whether the encrypted value is too long, if this is the case it is 
		// a plaintext extension and we have to parse this data
		if (len(fr.Value[fr.ValueStart()*2:]) / 2) > fr.ConfigurationLength() {
			return getDIFs(fr.Value[fr.ValueStart()*2+fr.ConfigurationLength()*2:], true)
		}

		//if(fr.ValueStart()+fr.ConfigurationLength())*2
		key, err := hex.DecodeString(fr.Key)
		if err == nil {
			ciphertext, err := hex.DecodeString(fr.Value[fr.ValueStart()*2:])
			if err == nil {
				block, err := aes.NewCipher(key)
				if err == nil {
					iv, err := hex.DecodeString(fr.IV())
					if err == nil {
						mode := cipher.NewCBCDecrypter(block, iv)
						mode.CryptBlocks(ciphertext, ciphertext)
						plaintext := fmt.Sprintf("%s\n", hex.EncodeToString(ciphertext))
						values := getDIFs(plaintext, false)
						return values
					}
				}
			}
		}
	}
	return nil
}

func getDIFs(originalString string, extension bool) []DIFs {
	var values []DIFs

	// Format the string properly
	originalString = strings.ToUpper(originalString)
	originalString = strings.TrimSpace(originalString)
	originalString = strings.TrimLeft(originalString, "2F")
	if strings.HasSuffix(originalString, "2F") {
		originalString = originalString[:len(originalString)-4]
	} else {
		originalString = originalString[:len(originalString)]
	}

	outString := strings.Split(originalString, "0E13")
	if len(outString) > 1 {
		outString = strings.Split(outString[1], "0DFD")
	}

	for i := range outString {
		var dif string
		var vif string
		var original string
		var value string
		switch i {
		case 0:
			if len(outString) > 1 {
				dif = "0E"
				vif = "13"
				value = outString[i][3:4] + outString[i][0:1] + outString[i][1:2]
				original = "0E 13 " + fmt.Sprintf("% X", outString[i])
			} else {
				dif = "0D"
				vif = "FD"
				ascii, _ := hex.DecodeString(fmt.Sprintf("%s", outString[i][8:len(outString[i])]))
				asciiOutput := fmt.Sprintf("%s", ascii)
				value = reverse(asciiOutput)

				// Decode string
				hexbyte, _ := hex.DecodeString(fmt.Sprintf("%s", outString[i][:len(outString[i])]))
				original = "0D FD " + fmt.Sprintf("% X", hexbyte)
			}
		case 1:
			dif = "0D"
			vif = "FD"
			ascii, _ := hex.DecodeString(fmt.Sprintf("%s", outString[i][:len(outString[i])]))
			asciiOutput := fmt.Sprintf("%s", ascii)
			value = reverse(asciiOutput)
			
			// Decode string
			hexbyte, _ := hex.DecodeString(fmt.Sprintf("%s", outString[i][:len(outString[i])]))
			original = "0D FD " + fmt.Sprintf("% X", hexbyte)
		}

		singleDIF := DIFs{DIF: dif, VIF: vif, Value: value, Original: original}

		values = append(values, singleDIF)
	}
	return values
}
