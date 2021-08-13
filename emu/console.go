package emu

import (
	"fmt"
	"os"
)

type Console struct {
	memory  [4096]byte
	buffer [32]uint64

	pc uint16
	i  uint16
	stack      []uint16
	delayTimer uint8	// TODO implement
	soundTimer uint8	// TODO implement

	v [16]uint8

	config *Config
}

const BUFFER_WIDTH = 64
const BUFFER_HEIGHT = 32

// LoadRom loads a rom into the console.
func (cons *Console) LoadRom(filename string) {
	file, _ := os.ReadFile(filename)

	for i, bt := range file {
		cons.memory[i+int(0x200)] = bt
	}
}

// Pushes a value into the stack.
func (cons *Console) push(value uint16) {
	cons.stack = append(cons.stack, value)
}

// Pops a value from the stack.
func (cons *Console) pop() uint16 {
	top := len(cons.stack) - 1
	elem := cons.stack[top]
	cons.stack = cons.stack[:top]
	return elem
}

// GetBit returns the bit positioned at (x, y).
func (cons *Console) GetBit(x, y uint8) uint8 {
	return uint8(cons.buffer[y]>>(BUFFER_WIDTH-1-x)) & 1
}

// Xors the bit at (x, y) with the argument.
func (chip *Console) xorBit(x, y, bit uint8) {
	o := uint64(bit) << (BUFFER_WIDTH - 1 - x)
	chip.buffer[y] ^= o
}

// Clears the screen.
func (cons *Console) clearBuffer() {
	cons.buffer = [32]uint64{}
}

// Prints the buffer to stdout.
// Used for debugging purposes.
func (cons *Console) printBuffer() {
	fmt.Println("DISPLAY")
	for y := 0; y < BUFFER_HEIGHT; y++ {
		for x := 0; x < BUFFER_WIDTH; x++ {
			bit := cons.GetBit(uint8(x), uint8(y))
			var c rune
			if bit == 0 {
				c = ' '
			} else {
				// full block
				c = '\u2588'
			}
			fmt.Printf("%c ", c)
		}
		fmt.Print("\n")
	}
}

// Cycles runs the next instruction.
// It returns true if the screen needs to be drawn.
func (chip *Console) Cycle() bool {
	return chip.decodeExecute(chip.fetch())
}

// NewConsole creates a new console.
func NewConsole() *Console {
	memory := [4096]byte{}
	display := [32]uint64{}
	stack := make([]uint16, 0)
	v := [16]uint8{}
	config := NewConfig(false)

	console := &Console{
		memory,
		display,
		// put the PC to the beginning of the ROM
		uint16(0x200),
		0,
		stack,
		0,
		0,
		v,
		config,
	}

	putFont(console)

	return console
}
