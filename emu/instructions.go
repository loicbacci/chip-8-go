package emu

import "math/rand"

func (chip *Chip) decodeExecute(instr uint16) bool {
	x := (instr & uint16(0x0F00)) >> 8
	y := (instr & uint16(0x00F0)) >> 4
	n := uint8(instr & uint16(0x000F))
	nn := uint8(instr & uint16(0x00FF))
	nnn := instr & (0x0FFF)

	needToDraw := false

	if instr == uint16(0x00E0) {
		// clear screen
		chip.clearScreen()
	}
	if instr == uint16(0x00EE) {
		// return
		chip.PC = chip.pop()
	}

	switch (instr & uint16(0xF000)) >> 12 {
	case uint16(0x1):
		// jump
		chip.PC = nnn

	case uint16(0x2):
		// call
		chip.push(chip.PC)
		chip.PC = nnn

	case uint16(0x3):
		// skip equal
		if chip.V[x] == nn {
			chip.PC += 2
		}

	case uint16(0x4):
		// skip not equal
		if chip.V[x] != nn {
			chip.PC += 2
		}

	case uint16(0x5):
		// skip equal register
		if chip.V[x] == chip.V[y] {
			chip.PC += 2
		}

	case uint16(0x6):
		// set register vx
		chip.V[x] = nn

	case uint16(0x8):
		switch n {
		case uint8(0x0):
			// set
			chip.V[x] = chip.V[y]

		case uint8(0x1):
			// OR
			chip.V[x] |= chip.V[y]

		case uint8(0x2):
			// AND
			chip.V[x] &= chip.V[y]

		case uint8(0x3):
			// XOR
			chip.V[x] ^= chip.V[y]

		case uint8(0x4):
			// add carry
			add := uint16(chip.V[x]) + uint16(chip.V[y])
			chip.V[x] = uint8(add)

			if add > 255 {
				chip.V[15] = 1
			} else {
				chip.V[15] = 0
			}

		case uint8(0x5):
			// sub (x - y)
			if chip.V[x] > chip.V[y] {
				chip.V[15] = 1
			} else {
				chip.V[15] = 0
			}
			chip.V[x] -= chip.V[y]

		case uint8(0x6):
			// shift right
			if chip.Config.cosmac {
				chip.V[x] = chip.V[y]
			}
			chip.V[15] = chip.V[x] & 1
			chip.V[x] >>= 1

		case uint8(0x7):
			// sub (y - x)
			if chip.V[y] > chip.V[x] {
				chip.V[15] = 1
			} else {
				chip.V[15] = 0
			}
			chip.V[x] = chip.V[y] - chip.V[x]

		case uint8(0xE):
			// shift left
			if chip.Config.cosmac {
				chip.V[x] = chip.V[y]
			}
			chip.V[15] = chip.V[x] & 1
			chip.V[x] <<= 1
		}

	case uint16(0x7):
		// add value to register vx
		chip.V[x] += nn

	case uint16(0x9):
		// skip not equal register
		if chip.V[x] != chip.V[y] {
			chip.PC += 2
		}

	case uint16(0xA):
		// set index register I
		chip.I = nnn

	case uint16(0xB):
		if chip.Config.cosmac {
			chip.PC = nnn + uint16(chip.V[0])
		} else {
			chip.PC = nnn + uint16(chip.V[x])
		}

	case uint16(0xC):
		rndByte := uint8(rand.Int31n(256))
		chip.V[x] = rndByte & nn

	case uint16(0xD):
		// display/draw
		xCoord := chip.V[x] % WIDTH
		yCoord := chip.V[y] % HEIGHT

		chip.V[15] = 0

		for i := uint8(0); i < n && yCoord+i < HEIGHT; i++ {
			bt := chip.Memory[chip.I+uint16(i)]

			for j := uint8(0); j < 8 && xCoord+j < WIDTH; j++ {
				bit := (bt >> (7 - j)) & 1

				if bit == 1 {
					chip.setBit(xCoord+j, yCoord+i, 1)

					if chip.GetBit(xCoord+j, yCoord+i) == 1 {
						chip.V[15] = 1
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
			chip.V[x] = chip.DelayTimer

		case uint8(0x15):
			// delay timer to vx
			chip.DelayTimer = chip.V[x]

		case uint8(0x18):
			// sound timer to vx
			chip.SoundTimer = chip.V[x]

		case uint8(0x1E):
			// add to index
			chip.I += uint16(chip.V[x])

			if !chip.Config.cosmac && chip.I > uint16(0x1000) {
				chip.V[15] = 1
			}

		case uint8(0x0A):
			// get key
			// TODO

		case uint8(0x29):
			// font character
			chip.I = GetCharAddr(chip.V[x])

		case uint8(0x33):
			// BCD
			nbr := chip.V[x]
			hundreds := (nbr / 100) % 10
			tens := (nbr / 10) % 10
			ones := nbr % 10
			chip.Memory[chip.I] = hundreds
			chip.Memory[chip.I + 1] = tens
			chip.Memory[chip.I + 2] = ones

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
