package chip8

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	MemorySize   = 4096
	ProgramStart = 0x200
)

type Chip8 struct {
	Memory  [MemorySize]uint8
	V       [16]uint8
	I       uint16
	PC      uint16
	Display [64 * 32]uint8
	Draw    bool
}

func NewChip8() *Chip8 {
	c := &Chip8{}
	c.PC = ProgramStart
	for i := 0; i < len(fontset); i++ {
		fmt.Println(0x50+i)
		c.Memory[0x50+i] = fontset[i]
	}
	return c
}

func (c *Chip8) LoadROM(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for i := 0; i < len(data); i++ {
		c.Memory[ProgramStart+i] = data[i]
	}
	return nil
}

func (c *Chip8) Cycle() {
	slog.Info("PC value in cycle", "PC", c.PC)
	opcode := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
	c.PC += 2

	switch opcode & 0xF000 {
	case 0x0000:
		if opcode == 0x00E0 {
			// 00E0: clear screen
			for i := range c.Display {
				c.Display[i] = 0
			}
			c.Draw = true
		}

	case 0x1000:
		// 1NNN: jump to address NNN
		addr := opcode & 0x0FFF
		c.PC = addr

	case 0x6000:
		// 6XNN: set VX = NN
		x := (opcode & 0x0F00) >> 8
		nn := uint8(opcode & 0x00FF)
		c.V[x] = nn

	case 0x7000:
		// 7XNN: VX += NN
		x := (opcode & 0x0F00) >> 8
		nn := uint8(opcode & 0x00FF)
		c.V[x] += nn

	case 0xA000:
		// ANNN: set I = NNN
		addr := opcode & 0x0FFF
		c.I = addr

	case 0xD000:
		// DXYN: draw sprite

		// example of maksing and right bit shift operation
		//   1101 0001 0010 0101
		// & 0000 0000 1111 0000
  		// -------------------
  		// - 0000 0000 0010 0000
		// then shift by 4 
		//   0000 0000 0010 0000 >> 4 = 0000 0000 0000 0010 = 2


		x := c.V[(opcode&0x0F00)>>8]
		y := c.V[(opcode&0x00F0)>>4]
		height := opcode & 0x000F
		
		slog.Default().Info("Draw sprite", "x", x, "y", y, "height", height)

		c.V[0xF] = 0 // VF = 0 (collision flag)
		for row := uint16(0); row < height; row++ {
			sprite := c.Memory[c.I+row]
			for col := uint16(0); col < 8; col++ {
				if (sprite & (0x80 >> col)) != 0 {
					// Wrap if needed
					px := (int(x) + int(col)) % 64
					py := (int(y) + int(row)) % 32
					i := py*64 + px

					if c.Display[i] == 1 {
						c.V[0xF] = 1
					}
					c.Display[i] ^= 1
				}
			}
		}
		c.Draw = true
	default:
		fmt.Printf("Unknown opcode: 0x%04X\n", opcode)
	}
}

func (c *Chip8) PrintDisplay() {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if c.Display[y*64+x] == 1 {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
