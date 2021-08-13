package emu

const FONT_OFFSET = int(0x50)

func (chip *Chip) PutFont() {
	i := FONT_OFFSET
	for _, digit := range font {
		for _, bt := range digit {
			chip.Memory[i] = bt
		}
	}
}

func GetCharAddr(ch uint8) uint16 {
	return uint16(FONT_OFFSET) + uint16(ch) * 5
}

var font = [16][5]byte{
	zero,
	one,
	two,
	three,
	four,
	five,
	six,
	seven,
	eight,
	nine,
	aFont,
	bFont,
	cFont,
	dFont,
	eFont,
	fFont,
}

var zero = [5]byte{
	byte(0xF0),
	byte(0x90),
	byte(0x90),
	byte(0x90),
	byte(0xF0),
}

var one = [5]byte{
	byte(0x20),
	byte(0x60),
	byte(0x20),
	byte(0x20),
	byte(0x70),
}

var two = [5]byte{
	byte(0xF0),
	byte(0x10),
	byte(0xF0),
	byte(0x80),
	byte(0xF0),
}

var three = [5]byte{
	byte(0xF0),
	byte(0x10),
	byte(0xF0),
	byte(0x10),
	byte(0xF0),
}

var four = [5]byte{
	byte(0x90),
	byte(0x90),
	byte(0xF0),
	byte(0x10),
	byte(0x10),
}

var five = [5]byte{
	byte(0xF0),
	byte(0x80),
	byte(0xF0),
	byte(0x10),
	byte(0xF0),
}

var six = [5]byte{
	byte(0xF0),
	byte(0x80),
	byte(0xF0),
	byte(0x90),
	byte(0xF0),
}

var seven = [5]byte{
	byte(0xF0),
	byte(0x10),
	byte(0x20),
	byte(0x40),
	byte(0x40),
}

var eight = [5]byte{
	byte(0xF0),
	byte(0x90),
	byte(0xF0),
	byte(0x90),
	byte(0xF0),
}

var nine = [5]byte{
	byte(0xF0),
	byte(0x90),
	byte(0xF0),
	byte(0x10),
	byte(0xF0),
}

var aFont = [5]byte{
	byte(0xF0),
	byte(0x90),
	byte(0xF0),
	byte(0x90),
	byte(0x90),
}

var bFont = [5]byte{
	byte(0xE0),
	byte(0x90),
	byte(0xE0),
	byte(0x90),
	byte(0xE0),
}

var cFont = [5]byte{
	byte(0xF0),
	byte(0x80),
	byte(0x80),
	byte(0x80),
	byte(0xF0),
}

var dFont = [5]byte{
	byte(0xE0),
	byte(0x90),
	byte(0x90),
	byte(0x90),
	byte(0xE0),
}

var eFont = [5]byte{
	byte(0xF0),
	byte(0x80),
	byte(0xF0),
	byte(0x80),
	byte(0xF0),
}

var fFont = [5]byte{
	byte(0xF0),
	byte(0x80),
	byte(0xF0),
	byte(0x80),
	byte(0x80),
}