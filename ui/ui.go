package ui

import (
	"os"

	sdl "github.com/veandco/go-sdl2/sdl"
	emu "github.com/baccigal/chip-8/emu"
)

func Run(chip *emu.Chip) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("chip-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		emu.WIDTH*10, emu.HEIGHT*10, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	canvas, _ := sdl.CreateRenderer(window, -1, 0)
	defer canvas.Destroy()

	for {
		// Draw
		if chip.Cycle() {
			canvas.SetDrawColor(255, 0, 0, 255)
			canvas.Clear()

			for y := uint8(0); y < emu.HEIGHT; y++ {
				for x := uint8(0); x < emu.WIDTH; x++ {
					if chip.GetBit(x, y) == 1 {
						canvas.SetDrawColor(255, 255, 255, 255)
					} else {
						canvas.SetDrawColor(0, 0, 0, 0)
					}
					canvas.FillRect(&sdl.Rect{
						X: int32(x) * 10,
						Y: int32(y) * 10,
						H: 10,
						W: 10,
					})
				}
			}

			canvas.Present()
		}

		// Poll
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			}
		}

		// Delay the clock
		sdl.Delay(1000 / 60)
	}
}
