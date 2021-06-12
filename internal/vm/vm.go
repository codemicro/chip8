package vm

import "time"

type memory [4*1024]byte

type uiDriver interface {
	PublishNewDisplay([32][64]bool)
	GetPressedKeys() []uint8
}

type Chip8 struct {
	ui uiDriver

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

func NewChip8(ui uiDriver) *Chip8 {
	// TODO: Load ROM here
	// TODO: Font

	c := &Chip8{
		ui: ui,
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

	for {
		var disp [32][64]bool
		disp[5][7] = true
		c.ui.PublishNewDisplay(disp)
		time.Sleep(time.Second)

		var disp2 [32][64]bool
		disp2[5][12] = true
		c.ui.PublishNewDisplay(disp2)
		time.Sleep(time.Second)
	}
}
