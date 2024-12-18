package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Kilemonn/flow/flow/config"
	"github.com/Kilemonn/flow/flow/serial"
)

const (
	APPLICATION_VERSION = "0.1.0"
	MENU_OPTION_HELP    = "help"

	MENU_OPTION_SERIAL    = "serial"
	MENU_OPTION_SERIAL_LS = "serialls"

	MENU_OPTION_CONFIG_APPLY = "config-apply"
)

func main() {
	var (
		comFlag        string
		baudFlag       uint
		parityFlag     string
		dataSizeFlag   uint
		twoStopBitFlag bool

		configFilePath string
	)

	flag.StringVar(&comFlag, "com", "", "com serial name")
	flag.UintVar(&baudFlag, "baud", 0, "baud rate")
	flag.StringVar(&parityFlag, "parity", "", "parity bit")
	flag.UintVar(&dataSizeFlag, "data-size", 0, "data size")
	flag.BoolVar(&twoStopBitFlag, "two-stop-bits", false, "two stop bit")

	flag.StringVar(&configFilePath, "f", "", "configuration file path")
	flag.Parse()

	if len(os.Args) <= 1 || len(flag.Args()) == 0 {
		printHelp()
		return
	}

	switch flag.Args()[0] {
	case MENU_OPTION_HELP:
		printHelp()
	case MENU_OPTION_SERIAL:
		serial.StartSerial(comFlag, baudFlag, parityFlag, dataSizeFlag, twoStopBitFlag)
	case MENU_OPTION_SERIAL_LS:
		serial.SerialList()
	case MENU_OPTION_CONFIG_APPLY:
		config.ApplyConfigurationFromFile(configFilePath)
		fmt.Println("Closed cleanly...")
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Printf("flow - cli v%s.\n", APPLICATION_VERSION)
	fmt.Printf("%s -f <file configuration path> - Create and apply the connection forwarding between reader and writers defined in the config file.\n", MENU_OPTION_CONFIG_APPLY)
	fmt.Printf("%s -com <COM0 or /dev/tty/USB0> -baud <baud rate> -parity <Even / Odd> -data-size <default is 8> -two-stop-bits <true is 2, false is 1 (default)> - Create a serial connection with device.\n", MENU_OPTION_SERIAL)
	fmt.Printf("%s - List connected serial devices.\n", MENU_OPTION_SERIAL_LS)
}
