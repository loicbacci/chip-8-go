package main

import (
	"fmt"
	"os"

	emu "github.com/baccigal/chip-8/emu"
	ui "github.com/baccigal/chip-8/ui"
)

func main() {
	var filename string
	if len(os.Args) < 2 {
		filename = "roms/IBM Logo.ch8"
	} else {
		filename = os.Args[1]
	}

	fmt.Println("Loading file:", filename)

	chip := emu.NewConsole()
	chip.LoadRom(filename)

	ui.Run(chip)
}
