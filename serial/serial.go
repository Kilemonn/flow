package serial

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Kilemonn/flow/stdio"
	goSerial "go.bug.st/serial"
)

const (
	POLLING_DELAY_MS = "100"
)

func parseParity(parity string) goSerial.Parity {
	parityVal := goSerial.NoParity
	if strings.ToLower(parity) == "even" {
		parityVal = goSerial.EvenParity
	} else if strings.ToLower(parity) == "odd" {
		parityVal = goSerial.OddParity
	}

	return parityVal
}

func parseStopBits(stopBits bool) goSerial.StopBits {
	stopBitsVal := goSerial.OneStopBit
	if stopBits {
		stopBitsVal = goSerial.TwoStopBits
	}

	return stopBitsVal
}

func parseSerialSettings(baud uint, parity string, dataLen uint, stopBits bool) (goSerial.Mode, error) {
	if dataLen == 0 {
		return goSerial.Mode{}, errors.New("data length must be set and greater than 0")
	} else if strings.ToLower(parity) != "even" && strings.ToLower(parity) != "odd" && len(parity) != 0 {
		return goSerial.Mode{}, errors.New("parity must be 'even' or 'odd'")
	} else if baud == 0 {
		return goSerial.Mode{}, errors.New("baud rate must be set and greater than 0")
	}

	mode := goSerial.Mode{
		BaudRate: int(baud),
		Parity:   parseParity(parity),
		DataBits: int(dataLen),
		StopBits: parseStopBits(stopBits),
	}

	return mode, nil
}

func OpenSerialConnection(com string, mode goSerial.Mode) (goSerial.Port, error) {
	return goSerial.Open(com, &mode)
}

func GetSerialPorts() []string {
	ports, err := goSerial.GetPortsList()
	if err != nil {
		fmt.Printf("Failed to retrieve ports list. Error: [%s]", err.Error())
		return []string{}
	}

	if len(ports) == 0 {
		fmt.Printf("No serial ports connected.")
		return []string{}
	}

	return ports
}

func printSerialPorts() {
	for _, port := range GetSerialPorts() {
		fmt.Printf("%v\n", port)
	}
}

func rxPrintThread(port goSerial.Port, writer io.Writer, newLine bool) error {
	buffer := make([]byte, 100)
	counter := 0
	if newLine {
		for {
			_, err := port.Read(buffer)
			if err != nil {
				return err
			}
			counter = counter + 1
			// fmt.Printf("Refresh%v = Received %v\n", counter, string(buffer))
			writer.Write(buffer)
		}
	} else {
		for {
			_, err := port.Read(buffer)
			if err != nil {
				return err
			}
			counter = counter + 1
			// fmt.Printf("Refresh%v = Received %v\n", counter, string(buffer))
			writer.Write(buffer)
		}
	}

}

func SerialList() {
	printSerialPorts()
}

func StartSerial(com string, baud uint, parity string, dataLen uint, stopBits bool) {
	mode, err := parseSerialSettings(baud, parity, dataLen, stopBits)
	if err != nil {
		fmt.Printf("Failed to parse serial settings. Error: [%s]", err.Error())
		return
	}

	port, err := OpenSerialConnection(com, mode)
	if err != nil {
		fmt.Printf("Failed to open serial connection. Error: [%s]", err.Error())
		return
	}

	reader, _ := stdio.CreateStdInReader()
	bytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Failed to read in all data from provided input stream. Error: [%s]", err.Error())
		return
	}

	// this needs to be fixed just testing things out
	writer, _ := stdio.CreateStdOutWriter()
	go rxPrintThread(port, writer, false)
	time.Sleep(1 * time.Second)

	// "ACDEF\n\r"
	n, err := port.Write(bytes)
	if err != nil {
		fmt.Printf("Failed to write bytes to port [%s]. Error [%s].\n", com, err.Error())
		return
	}
	fmt.Printf("Sent %d bytes\n", n)
	time.Sleep(2 * time.Second)
}
