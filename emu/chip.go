package emu

import (
	"fmt"
	"os"
)

type Chip struct {
	Memory  [4096]byte
	Display [32]uint64

	PC uint16
	I  uint16
	// TODO STACK
	Stack      []uint16
	DelayTimer uint8
	SoundTimer uint8

	V [16]uint8

	Config *Config
}

const WIDTH = 64
const HEIGHT = 32

func (chip *Chip) LoadRom(filename string) {
	file, _ := os.ReadFile(filename)

	for i, bt := range file {
		chip.Memory[i+int(0x200)] = bt
	}
}

func (chip *Chip) fetch() uint16 {
	b0 := chip.Memory[chip.PC]
	b1 := chip.Memory[chip.PC+1]
	instruction := (uint16(b0) << 8) | uint16(b1)

	chip.PC += 2
	return instruction
}

func (chip *Chip) push(elem uint16) {
	chip.Stack = append(chip.Stack, elem)
}

func (chip *Chip) pop() uint16 {
	top := len(chip.Stack) - 1
	elem := chip.Stack[top]
	chip.Stack = chip.Stack[:top]
	return elem
}

func (chip *Chip) GetBit(x, y uint8) uint8 {
	return uint8(chip.Display[y]>>(WIDTH-1-x)) & 1
}

func (chip *Chip) setBit(x, y, bit uint8) {
	o := uint64(bit) << (WIDTH - 1 - x)
	chip.Display[y] ^= o
}

func (chip *Chip) clearScreen() {
	chip.Display = [32]uint64{}
}

func (chip *Chip) showScreen() {
	fmt.Println("DISPLAY")
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			bit := chip.GetBit(uint8(x), uint8(y))
			var c rune
			if bit == 0 {
				c = ' '
			} else {
				c = '\u2588'
			}
			fmt.Printf("%c ", c)
		}
		fmt.Print("\n")
	}
}

func (chip *Chip) Cycle() bool {
	return chip.decodeExecute(chip.fetch())
}



func NewChip() *Chip {
	memory := [4096]byte{}
	display := [32]uint64{}
	stack := make([]uint16, 0)
	v := [16]uint8{}
	config := NewConfig(false)

	chip := Chip{
		memory,
		display,
		uint16(0x200),
		0,
		stack,
		0,
		0,
		v,
		config,
	}

	chip.PutFont()

	return &chip
}
