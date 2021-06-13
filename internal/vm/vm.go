package vm

import (
	"fmt"
	"time"
)

type memory [4 * 1024]byte

type uiDriver interface {
	PublishNewDisplay([32][64]bool)
	GetPressedKeys() []uint8
	StartTone()
	StopTone()
}

type Chip8 struct {
	Debug bool

	// CopyRegistersOnShift affects `8XY6` and `8XYE`. If true, the value of VY will be copied into VX before a shift
	// occurs.
	CopyRegistersOnShift bool
	// VariableOffsetRegister affects `BNNN`. If true, `BNNN` will act as `BXNN` and will jump to NNN + VX. Else, `BNNN`
	// will jump to NNN + V0.
	VariableOffsetRegister bool
	// DisableSetFlagOnIrOverflow affects `FX1E`. If true, `FX1E` will not set VF. Else, it will set VF accordingly if
	// the index register "overflows" above 0x0FFF.
	DisableSetFlagOnIrOverflow bool
	// IncrementIndexRegisterOnLoadSave affects `FX55` and `FX65`. If true, the index register will be incremented when
	// loading or saving registers to/from memory. Else, a temporary value will be indexed instead, and the index
	// register will not be changed.
	IncrementIndexRegisterOnLoadSave bool

	ui              uiDriver
	clockSpeedHertz int
	disp            [32][64]bool

	// Main memory
	memory memory

	// Registers
	cir   [2]byte // current instruction register
	pc    uint16  // program counter
	ir    uint16  // index register
	stack Stack   // call stack
	delay uint8
	sound uint8

	// General purpose registers
	v0, v1, v2, v3, v4, v5, v6, v7, v8, v9, va, vb, vc, vd, ve, vf byte
}

func NewChip8(rom []byte, ui uiDriver, clockSpeedHertz int) *Chip8 {

	c := &Chip8{
		ui:              ui,
		clockSpeedHertz: clockSpeedHertz,

		pc: 0x200,
	}

	// load ROM
	for i, b := range rom {
		c.memory[int(c.pc) + i] = b
	}

	loadFont(&c.memory)

	// original COSMAC VIP settings
	//c.CopyRegistersOnShift = true
	//c.VariableOffsetRegister = false
	//c.DisableSetFlagOnIrOverflow = true
	//c.IncrementIndexRegisterOnLoadSave = true

	// Super Chip settings
	//c.CopyRegistersOnShift = false
	//c.VariableOffsetRegister = true
	//c.DisableSetFlagOnIrOverflow = false
	//c.IncrementIndexRegisterOnLoadSave = false

	c.CopyRegistersOnShift = true
	c.VariableOffsetRegister = false
	c.DisableSetFlagOnIrOverflow = false
	c.IncrementIndexRegisterOnLoadSave = false

	return c
}

// decrement decrements *v by 1 if it's not zero
func decrement(v *uint8) {
	if *v != 0 {
		*v -= 1
	}
}

func (c *Chip8) Run() {

	programTicker := time.NewTicker(time.Second / time.Duration(c.clockSpeedHertz))
	defer programTicker.Stop()

	decrementTicker := time.NewTicker(time.Second / 60)
	defer decrementTicker.Stop()

	done := make(chan bool)

MAINLOOP:
	for {
		select {
		case <-done:
			break MAINLOOP
		case <-decrementTicker.C:
			decrement(&c.delay)
			decrement(&c.sound)

			if c.sound == 0 {
				c.ui.StopTone()
			} else {
				c.ui.StartTone()
			}

			// TODO: Make noise!
		case <-programTicker.C:

			// FETCH
			c.cir[0] = c.memory[c.pc]
			c.cir[1] = c.memory[c.pc+1]
			c.pc += 2

			if c.Debug {
				fmt.Printf("DEBUG: pc:0x%04x cir:0x%04x\n", c.pc - 2, c.cir)
				fmt.Print("       ")
				for i := 0; i < 16; i += 1 {
					fmt.Printf("v%x:%02x ", i, *c.getRegisterPointer(byte(i)))
				}
				fmt.Printf("\n       ir:%04x\n", c.ir)
			}

			// DECODE + EXECUTE
			switch c.cir[0] & 0xF0 { // first four bits
			case 0x00:

				switch c.cir {
				case [2]byte{0x00, 0xE0}:
					// 00E0 - clear screen
					c.clearScreen()
				case [2]byte{0x00, 0xEE}:
					// 00EE - subroutine return
					c.subroutineReturn()
				default:
					panic(fmt.Errorf("UNHANDLED at %x: %x\n", c.pc, c.cir))
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
				c.setRegisterToConstant()

			case 0x70:
				// 7XNN - add NN to VX without setting carry flag
				c.addConstantToRegister()

			case 0x80:
				// Logical and arithmetic instructions
				switch c.cir[1] & 0x0F { // switch on the last four bits - eg 00000000 00001111
				case 0x00:
					// 8XY0 - set VX to VY
					c.setRegisterToRegister()
				case 0x01:
					// 8XY1 - set VX to logical OR of VX and VY
					c.setRegisterToLogicalOr()
				case 0x02:
					// 8XY2 - set VX to logical AND of VX and VY
					c.setRegisterToLogicalAnd()
				case 0x03:
					// 8XY3 - set VX to logical XOR of VX and VY
					c.setRegisterToLogicalXor()
				case 0x04:
					// 8XY4 - set VX to the sum of VX and VY then set the carry flag as appropriate
					c.setRegisterToSum()
				case 0x05:
					// 8XY5 - set VX to VX - VY then set the carry flag as appropriate
					c.setRegisterToDifferenceA()
				case 0x06:
					// 8XY6 - set VX to VY (if CopyRegistersOnShift), shift the value of VX 1 bit right and set VF to
					// the bit shifted out
					c.shiftRight()
				case 0x07:
					// 8XY7 - set VX to VY - VX then set the carry flag as appropriate
					c.setRegisterToDifferenceB()
				case 0x0E:
					// 8XYE - set VX to VY (if CopyRegistersOnShift), shift the value of VX 1 bit left and set VF to the
					// bit shifted out
					c.shiftLeft()
				default:
					panic(fmt.Errorf("UNHANDLED at %x: %x\n", c.pc, c.cir))
				}

			case 0x90:
				// 9XYN - skip one if registers not equal
				c.skipNotEqRegReg()

			case 0xA0:
				// ANNN - set index register to constant
				c.setIndexRegister()

			case 0xB0:
				// BNNN - set PC to NNN + V0 - if VariableOffsetRegister, BXNN - set PC to XNN + VX
				c.jumpWithOffset()

			case 0xC0:
				// CXNN - generate a random byte, AND it with NN and store in VX
				c.random()

			case 0xD0:
				// DXYN - display
				c.display()

			case 0xE0:
				// Input
				switch c.cir[1] {
				case 0x9E:
					// EX9E - skip one if key with the value stored in VX is pressed
					c.skipIfKey()
				case 0xA1:
					// EXA1 - skip one if key with the value stored in VX is not pressed
					c.skipIfNotKey()
				default:
					panic(fmt.Errorf("UNHANDLED at %x: %x\n", c.pc, c.cir))
				}

			case 0xF0:
				switch c.cir[1] {
				case 0x07:
					// FX07 - set value of VX to the current value of the delay timer
					c.getDelayTimer()
				case 0x15:
					// FX15 - set delay timer to value of VX
					c.setDelayTimer()
				case 0x18:
					// FX18 - set sound timer to the value of VX
					c.setSoundTimer()
				case 0x0A:
					// FX0A - blocks until a key is pressed. Stores that key's value in VX then continues.
					c.getPressedKey()
				case 0x1E:
					// FX1E - adds the value of VX to the index register and set VF accordingly if the index register
					// "overflows" above 0x0FFF
					c.addToIndexRegister()
				case 0x29:
					// FX29 - set the index register to the address of the hex character in VX
					c.getFontCharacter()
				case 0x33:
					// FX33 - take the value of VX, converts it to a denary number and the put each individual digit in
					// the memory location specified by the index register + the digit number
					c.convertToDecimal()
				case 0x55:
					// FX55 - store the value of each general purpose register from V0 to VX inclusive in consecutive
					// memory addresses starting from the current value of the index register
					c.storeMemory()
				case 0x65:
					// FX65 - loads the value of each general purpose register from V0 to VX inclusive from consecutive
					// memory addresses starting from the current value of the index register
					c.loadMemory()
				default:
					panic(fmt.Errorf("UNHANDLED at %x: %x\n", c.pc, c.cir))
				}

			default:
				panic(fmt.Errorf("UNHANDLED at %x: %x\n", c.pc, c.cir))
			}

		}
	}
}
