package vm

import (
	"fmt"
	"time"
)

type memory [4 * 1024]byte

type uiDriver interface {
	PublishNewDisplay([32][64]bool)
	GetPressedKeys() []uint8
}

type Chip8 struct {
	Debug bool

	ui              uiDriver
	clockSpeedHertz int
	disp            [32][64]bool

	// Main memory
	memory [4 * 1024]byte

	// Registers
	cir   [2]byte // current instruction register
	pc    uint16  // program counter
	ir    uint16  // index register
	stack Stack
	delay uint8
	sound uint8

	// General purpose registers
	v0, v1, v2, v3, v4, v5, v6, v7, v8, v9, va, vb, vc, vd, ve, vf byte
}

func NewChip8(rom []byte, ui uiDriver, clockSpeedHertz int) *Chip8 {
	// TODO: Load ROM here
	// TODO: Font

	c := &Chip8{
		ui:              ui,
		clockSpeedHertz: clockSpeedHertz,

		pc: 0x200,
	}

	for i, b := range rom {
		c.memory[int(c.pc) + i] = b
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

MAINLOOP:
	for {
		select {
		case <-done:
			break MAINLOOP
		case <-ticker.C:

			// FETCH
			c.cir[0] = c.memory[c.pc]
			c.cir[1] = c.memory[c.pc+1]
			c.pc += 2

			if c.Debug {
				fmt.Printf("DEBUG: 0x%04x\n", c.cir)
			}

			// DECODE
			switch c.cir[0] & 0xF0 { // first four bytes
			case 0x00:

				switch c.cir {
				case [2]byte{0x00, 0xE0}:
					// 00E0 - clear screen
					c.clearScreen()
				case [2]byte{0x00, 0xEE}:
					// 00EE - subroutine return
					c.subroutineReturn()
				default:
					fmt.Printf("UNHANDLED at %x: %x\n", c.pc, c.cir)
				}

			case 0x10:
				// 1NNN - jump
				c.jump()

			case 0x20:
				// 2NNN - subroutine call
				c.subroutineCall()

			case 0x30:
				// 3XNN - skip one if register equal to constant
				c.skipEqRegConst()

			case 0x40:
				// 4XNN - skip one if register not equal to constant
				c.skipNotEqRegConst()

			case 0x50:
				// 5XY0 - skip one if registers equal
				c.skipEqRegReg()

			case 0x60:
				// 6XNN - set VX to NN
				c.setRegister()

			case 0x70:
				// 7XNN - add NN to VX
				c.addToRegister()

			case 0x90:
				// 9XYN - skip one if registers not equal
				c.skipNotEqRegReg()

			case 0xA0:
				// ANNN - set index register to constant
				c.setIndexRegister()

			case 0xD0:
				// DXYN - display
				c.display()

			default:
				fmt.Printf("UNHANDLED at %x: %x\n", c.pc, c.cir)
			}

		}
	}
}
