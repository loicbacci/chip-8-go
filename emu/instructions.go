package emu

import "math/rand"

// fetch Fetches the next instruction and moves the pc.
func (cons *Console) fetch() uint16 {
	b0 := cons.memory[cons.pc]
	b1 := cons.memory[cons.pc+1]
	instruction := (uint16(b0) << 8) | uint16(b1)

	cons.pc += 2
	return instruction
}

// decodeExecute decodes the instruction and executes it.
// It returns true if the screen needs to be drawn.
func (cons *Console) decodeExecute(instr uint16) bool {
	// Compute instructions arguments
	x := (instr & uint16(0x0F00)) >> 8
	y := (instr & uint16(0x00F0)) >> 4
	n := uint8(instr & uint16(0x000F))
	nn := uint8(instr & uint16(0x00FF))
	nnn := instr & (0x0FFF)

	needToDraw := false

	if instr == uint16(0x00E0) {
		// clear screen
		cons.clearBuffer()
	}
	if instr == uint16(0x00EE) {
		// return from subroutine
		cons.pc = cons.pop()
	}

	switch (instr & uint16(0xF000)) >> 12 {
	case uint16(0x1):
		// jump
		cons.pc = nnn

	case uint16(0x2):
		// call
		cons.push(cons.pc)
		cons.pc = nnn

	case uint16(0x3):
		// skip equal
		if cons.v[x] == nn {
			cons.pc += 2
		}

	case uint16(0x4):
		// skip not equal
		if cons.v[x] != nn {
			cons.pc += 2
		}

	case uint16(0x5):
		// skip equal register
		if cons.v[x] == cons.v[y] {
			cons.pc += 2
		}

	case uint16(0x6):
		// set register vx
		cons.v[x] = nn

	case uint16(0x8):
		switch n {
		case uint8(0x0):
			// set
			cons.v[x] = cons.v[y]

		case uint8(0x1):
			// OR
			cons.v[x] |= cons.v[y]

		case uint8(0x2):
			// AND
			cons.v[x] &= cons.v[y]

		case uint8(0x3):
			// XOR
			cons.v[x] ^= cons.v[y]

		case uint8(0x4):
			// add carry
			add := uint16(cons.v[x]) + uint16(cons.v[y])
			cons.v[x] = uint8(add)

			if add > 255 {
				cons.v[15] = 1
			} else {
				cons.v[15] = 0
			}

		case uint8(0x5):
			// sub (x - y)
			if cons.v[x] > cons.v[y] {
				cons.v[15] = 1
			} else {
				cons.v[15] = 0
			}
			cons.v[x] -= cons.v[y]

		case uint8(0x6):
			// shift right
			if cons.config.cosmac {
				cons.v[x] = cons.v[y]
			}
			cons.v[15] = cons.v[x] & 1
			cons.v[x] >>= 1

		case uint8(0x7):
			// sub (y - x)
			if cons.v[y] > cons.v[x] {
				cons.v[15] = 1
			} else {
				cons.v[15] = 0
			}
			cons.v[x] = cons.v[y] - cons.v[x]

		case uint8(0xE):
			// shift left
			if cons.config.cosmac {
				cons.v[x] = cons.v[y]
			}
			cons.v[15] = cons.v[x] & 1
			cons.v[x] <<= 1
		}

	case uint16(0x7):
		// add value to register vx
		cons.v[x] += nn

	case uint16(0x9):
		// skip not equal register
		if cons.v[x] != cons.v[y] {
			cons.pc += 2
		}

	case uint16(0xA):
		// set index register I
		cons.i = nnn

	case uint16(0xB):
		if cons.config.cosmac {
			cons.pc = nnn + uint16(cons.v[0])
		} else {
			cons.pc = nnn + uint16(cons.v[x])
		}

	case uint16(0xC):
		rndByte := uint8(rand.Int31n(256))
		cons.v[x] = rndByte & nn

	case uint16(0xD):
		// display/draw
		xCoord := cons.v[x] % BUFFER_WIDTH
		yCoord := cons.v[y] % BUFFER_HEIGHT

		cons.v[15] = 0

		for i := uint8(0); i < n && yCoord+i < BUFFER_HEIGHT; i++ {
			bt := cons.memory[cons.i+uint16(i)]

			for j := uint8(0); j < 8 && xCoord+j < BUFFER_WIDTH; j++ {
				bit := (bt >> (7 - j)) & 1

				if bit == 1 {
					cons.xorBit(xCoord+j, yCoord+i, 1)

					if cons.GetBit(xCoord+j, yCoord+i) == 1 {
						cons.v[15] = 1
					}
				}
			}
		}

		needToDraw = true

	case uint16(0xE):
		switch nn {
		case uint8(0x9E):
			// skip if key pressed
			// TODO

		case uint8(0xA1):
			// skip if key not pressed
			// TODO 
		}

	case uint16(0xF):
		switch nn {
		case uint8(0x07):
			// vx to delay timer
			cons.v[x] = cons.delayTimer

		case uint8(0x15):
			// delay timer to vx
			cons.delayTimer = cons.v[x]

		case uint8(0x18):
			// sound timer to vx
			cons.soundTimer = cons.v[x]

		case uint8(0x1E):
			// add to index
			cons.i += uint16(cons.v[x])

			if !cons.config.cosmac && cons.i > uint16(0x1000) {
				cons.v[15] = 1
			}

		case uint8(0x0A):
			// get key
			// TODO

		case uint8(0x29):
			// font character
			cons.i = getCharAddr(cons.v[x])

		case uint8(0x33):
			// BCD
			nbr := cons.v[x]
			hundreds := (nbr / 100) % 10
			tens := (nbr / 10) % 10
			ones := nbr % 10
			cons.memory[cons.i] = hundreds
			cons.memory[cons.i + 1] = tens
			cons.memory[cons.i + 2] = ones

		case uint8(0x55):
			// store
			// TODO

		case uint8(0x65):
			// load
			// TODO
		}
	}

	return needToDraw
}
