package vm

import (
	"math/rand"
	"time"
)

type memory [4*1024]byte

type uiDriver interface {
	PublishNewDisplay([32][64]bool)
	GetPressedKeys() []uint8
}

type Chip8 struct {
	ui uiDriver
	clockSpeedHertz int

	// Main memory
	memory [4*1024]byte

	// Registers
	programCounter uint16
	indexRegister uint16
	stack Stack
	delay uint8
	sound uint8

	// General purpose registers
	v0, v1, v2, v3, v4, v5, v6, v7, v8, v9, va, vb, vc, vd, ve, vf [16]byte
}

func NewChip8(ui uiDriver, clockSpeedHertz int) *Chip8 {
	// TODO: Load ROM here
	// TODO: Font

	c := &Chip8{
		ui: ui,
		clockSpeedHertz: clockSpeedHertz,

		programCounter: 0x200,
	}

	return c
}

// descendingLoop decrements *v by 1 every 1/60 of a second
func descendingLoop(v *uint8) {
	delay := time.Second / 60
	go func() {
		if *v == 0 {
			*v = 255
		} else {
			*v -= 1
		}
		time.Sleep(delay)
	}()
}

func (c *Chip8) Run() {
	go descendingLoop(&c.delay)
	go descendingLoop(&c.sound)

	ticker := time.NewTicker(time.Second / time.Duration(c.clockSpeedHertz))
	defer ticker.Stop()
	done := make(chan bool)

	n := 0

MAINLOOP:
	for {
		select {
		case <-done:
			break MAINLOOP
		case <-ticker.C:
			if n == 30 {
				var disp [32][64]bool
				disp[rand.Intn(31)][rand.Intn(63)] = true
				c.ui.PublishNewDisplay(disp)

				n = 0
			} else {
				n += 1
			}
		}
	}
}
