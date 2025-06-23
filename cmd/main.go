package main

import (
	"log"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/tren03/chip8/internal/chip8"
)

const (
	scale        = 10
	screenWidth  = 64 * scale
	screenHeight = 32 * scale
)

func main() {
	// tests.TestGen()

	c := chip8.NewChip8()
	if err := c.LoadROM("roms/3-test-opcode.ch8"); err != nil {
		log.Fatal(err)
	}

	rl.InitWindow(int32(screenWidth), int32(screenHeight), "CHIP-8 Emulator")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		c.Cycle()

		if c.Draw {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			for y := 0; y < 32; y++ {
				for x := 0; x < 64; x++ {
					if c.Display[y*64+x] == 1 {
						rl.DrawRectangle(int32(x*scale), int32(y*scale), scale, scale, rl.White)
					}
				}
			}

			rl.EndDrawing()
			c.Draw = false
		}

		time.Sleep(time.Second / 600) // Optional: cap emulation speed
	}
}