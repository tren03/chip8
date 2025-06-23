package tests

import (
	"log/slog"
	"os"
)

func TestGen() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Our test instructions (MSB first)
	rom := []byte{
		0x00, 0xE0, // Clear screen
		0x60, 0x08, // Set V0 = 8 (X position)
		0x61, 0x10, // Set V1 = 16 (Y position)
		0xA0, 0x82, // I = 0x82 (sprite for A)
		0xD0, 0x15, // Draw at (V0, V1), height 5
	}

	err := os.WriteFile("roms/test.ch8", rom, 0644)
	if err != nil {
		logger.Error("Failed to write ROM", "error", err)
		return
	}
	logger.Info("created test files")
}
