package ui

import (
	"fmt"
	"os"

	emu "github.com/baccigal/chip-8/emu"
	sdl "github.com/veandco/go-sdl2/sdl"
)

const PIXEL_SIZE = 20

// Run launches the GUI and the emulator
func Run(cons *emu.Console) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("chip-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		emu.BUFFER_WIDTH*PIXEL_SIZE, emu.BUFFER_HEIGHT*PIXEL_SIZE, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	canvas, _ := sdl.CreateRenderer(window, -1, 0)
	defer canvas.Destroy()

	for {
		// Draw
		if cons.Cycle() {
			canvas.SetDrawColor(255, 0, 0, 255)
			canvas.Clear()

			for y := uint8(0); y < emu.BUFFER_HEIGHT; y++ {
				for x := uint8(0); x < emu.BUFFER_WIDTH; x++ {
					if cons.GetBit(x, y) == 1 {
						canvas.SetDrawColor(255, 255, 255, 255)
					} else {
						canvas.SetDrawColor(0, 0, 0, 0)
					}
					canvas.FillRect(&sdl.Rect{
						X: int32(x) * PIXEL_SIZE,
						Y: int32(y) * PIXEL_SIZE,
						H: PIXEL_SIZE,
						W: PIXEL_SIZE,
					})
				}
			}

			canvas.Present()
		}

		// Poll
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch et := e.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)

			case *sdl.KeyboardEvent:
				handleKey(et, cons)
			}
		}

		// Delay the clock
		sdl.Delay(1000 / 60)
	}
}

func handleKey(et *sdl.KeyboardEvent, cons *emu.Console) {
	var up bool
	if et.Type == sdl.KEYUP {
		up = false
	} else if et.Type == sdl.KEYDOWN {
		up = true
	} else {
		return
	}

	var index int
	switch et.Keysym.Sym {
	case sdl.K_1:
		index = 0x1
	case sdl.K_2:
		index = 0x2
	case sdl.K_3:
		index = 0x3
	case sdl.K_4:
		index = 0xC

	case sdl.K_q:
		index = 0x4
	case sdl.K_w:
		index = 0x5
	case sdl.K_e:
		index = 0x6
	case sdl.K_r:
		index = 0xD

	case sdl.K_a:
		index = 0x7
	case sdl.K_s:
		index = 0x8
	case sdl.K_d:
		index = 0x9
	case sdl.K_f:
		index = 0xE

	case sdl.K_y:
		index = 0xA
	case sdl.K_x:
		index = 0x0
	case sdl.K_c:
		index = 0xB
	case sdl.K_v:
		index = 0xF
	}

	cons.SetKey(index, up)
	fmt.Printf("Key %x is %v\n", index, up)
}
