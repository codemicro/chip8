// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/emulator/vm/execute.go

package vm

import (
	"encoding/binary"
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// clearScreen - 00E0
func (c *Chip8) clearScreen() {
	c.disp = [32][64]bool{}
	c.ui.PublishNewDisplay(c.disp)
}

// subroutineReturn - 00EE
func (c *Chip8) subroutineReturn() {
	c.pc = c.stack.Pop()
}

// jump - 1NNN
func (c *Chip8) jump() {
	c.pc = c.getAddressFromCIR()
}

// subroutineCall - 2NNN
func (c *Chip8) subroutineCall() {
	nnn := c.getAddressFromCIR()
	c.stack.Push(c.pc)
	c.pc = nnn
}

// skipEqRegConst - 3XNN skip one if register equal to constant
func (c *Chip8) skipEqRegConst() {
	nn := c.get8bitConstant()
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	if nn == *vx {
		c.pc += 2
	}
}

// skipNotEqRegConst - 4XNN skip one if register not equal to constant
func (c *Chip8) skipNotEqRegConst() {
	nn := c.get8bitConstant()
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	if nn != *vx {
		c.pc += 2
	}
}

// skipEqRegReg - 5XY0 skip one if registers equal
func (c *Chip8) skipEqRegReg() {
	x := c.getRegisterPointer(c.cir[0] & 0x0F)
	y := c.getRegisterPointer(c.cir[1] >> 4)
	if *x == *y {
		c.pc += 2
	}
}

// setRegisterToConstant - 6XNN set VX to NN
func (c *Chip8) setRegisterToConstant() {
	nn := c.get8bitConstant()
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	*vx = nn
}

// addConstantToRegister - 7XNN add NN to VX without setting carry flag
func (c *Chip8) addConstantToRegister() {
	nn := c.get8bitConstant()
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	*vx += nn
}

// setRegisterToRegister - 8XY0 set VX to VY
func (c *Chip8) setRegisterToRegister() {
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vy := c.getRegisterPointer(c.cir[1] >> 4)
	*vx = *vy
}

// setRegisterToLogicalOr - 8XY1 set VX to logical OR of VX and VY
func (c *Chip8) setRegisterToLogicalOr(){
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vy := c.getRegisterPointer(c.cir[1] >> 4)
	*vx = *vx | *vy
}

// setRegisterToLogicalAnd - 8XY2 set VX to logical AND of VX and VY
func (c *Chip8) setRegisterToLogicalAnd(){
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vy := c.getRegisterPointer(c.cir[1] >> 4)
	*vx = *vx & *vy
}

// setRegisterToLogicalXor - 8XY3 set VX to logical XOR of VX and VY
func (c *Chip8) setRegisterToLogicalXor(){
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vy := c.getRegisterPointer(c.cir[1] >> 4)
	*vx = *vx ^ *vy
}

// setRegisterToSum - 8XY4 set VX to the sum of VX and VY then set the carry flag as appropriate
func (c *Chip8) setRegisterToSum() {
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vy := c.getRegisterPointer(c.cir[1] >> 4)
	vf := &c.vf

	res := int(*vx) + int(*vy)

	*vx += *vy

	if res > 255 {
		*vf = 0x01
	} else {
		*vf = 0x00
	}
}

// setRegisterToDifferenceA - 8XY5 set VX to VX - VY then set the carry flag as appropriate
func (c *Chip8) setRegisterToDifferenceA() {
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vy := c.getRegisterPointer(c.cir[1] >> 4)
	vf := &c.vf

	setCarry := *vx > *vy

	*vx -= *vy

	if setCarry {
		*vf = 0x01
	} else {
		*vf = 0x00
	}
}

// shiftRight - 8XY6 set VX to VY (if CopyRegistersOnShift), shift the value of VX 1 bit right and set VF to the bit
// shifted out
func (c *Chip8) shiftRight() {
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vf := &c.vf

	if c.CopyRegistersOnShift {
		vy := c.getRegisterPointer(c.cir[1] >> 4)
		*vx = *vy
	}

	shiftedBit := *vx & 0x01

	*vx = *vx >> 1
	*vf = shiftedBit
}

// shiftLeft - 8XYE set VX to VY (if CopyRegistersOnShift), shift the value of VX 1 bit left and set VF to the bit
// shifted out
func (c *Chip8) shiftLeft() {
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vf := &c.vf

	if c.CopyRegistersOnShift {
		vy := c.getRegisterPointer(c.cir[1] >> 4)
		*vx = *vy
	}

	shiftedBit := (*vx & 0x80) >> 7

	*vx = *vx << 1
	*vf = shiftedBit
}

// setRegisterToDifferenceB - 8XY7 set VX to VY - VX then set the carry flag as appropriate
func (c *Chip8) setRegisterToDifferenceB() {
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	vy := c.getRegisterPointer(c.cir[1] >> 4)
	vf := &c.vf

	setCarry := *vx < *vy

	*vx = *vy - *vx

	if setCarry {
		*vf = 0x01
	} else {
		*vf = 0x00
	}
}

// skipNotEqRegReg - 9XY0 skip one if registers not equal
func (c *Chip8) skipNotEqRegReg() {
	regX := c.getRegisterPointer(c.cir[0] & 0x0F)
	regY := c.getRegisterPointer(c.cir[1] >> 4)
	if *regX != *regY {
		c.pc += 2
	}
}

// setIndexRegister - ANNN set index register to NNN
func (c *Chip8) setIndexRegister() {
	c.ir = c.getAddressFromCIR()
}

// jumpWithOffset - BNNN set PC to NNN + V0 - if VariableOffsetRegister, BXNN set PC to XNN + VX
func (c *Chip8) jumpWithOffset() {
	nnn := c.getAddressFromCIR()
	offset := c.v0

	if c.VariableOffsetRegister {
		offset = *c.getRegisterPointer(c.cir[0] & 0x0F)
	}

	c.pc = nnn + uint16(offset)
}

// random - CXNN generate a random byte, AND it with NN and store in VX
func (c *Chip8) random() {
	rnd := make([]byte, 4)
	binary.BigEndian.PutUint32(rnd, random.Uint32())

	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	nn := c.get8bitConstant()

	*vx = rnd[0] & nn
}

// display - DXYN draw an N pixel tall sprite from the memory location in the index register at the coordinate of the
// values in (VX, VY)
func (c *Chip8) display() {
	spriteHeight := c.get4BitConstant()

	startingXCoord := int(*c.getRegisterPointer(c.cir[0] & 0x0F) % 64)
	startingYCoord := int(*c.getRegisterPointer(c.cir[1] >> 4) % 32)

	vf := c.getRegisterPointer(0x0F)
	*vf = 0x00

	for y := 0; y < int(spriteHeight); y += 1 {

		if startingYCoord+y >= 32 { // if we'll be trying to draw out of bounds
			continue
		}

		rowData := c.memory[int(c.ir)+y]
		for x := 0; x < 8; x += 1 {

			if startingXCoord+x >= 64 {
				continue
			}

			pixelData := rowData & 0x80 // get most significant bit - ie, 10000000

			if pixelData == 0x80 { // if most significant bit is set
				currentValue := c.disp[startingYCoord+y][startingXCoord+x]
				c.disp[startingYCoord+y][startingXCoord+x] = !currentValue
				if currentValue {
					*vf = 0x01
				}
			}

			rowData = rowData << 1
		}
	}

	c.ui.PublishNewDisplay(c.disp)
}

// skipIfKey - EX9E skip one if key with the value stored in VX is pressed
func (c *Chip8) skipIfKey() {
	vxn := *c.getRegisterPointer(c.cir[0] & 0x0F)
	pressedKeys := c.ui.GetPressedKeys()

	for _, key := range pressedKeys {
		if key == vxn {
			c.pc += 2
			break
		}
	}
}

// skipIfNotKey - EXA1 skip one if key with the value stored in VX is not pressed
func (c *Chip8) skipIfNotKey() {
	vxn := *c.getRegisterPointer(c.cir[0] & 0x0F)
	pressedKeys := c.ui.GetPressedKeys()

	for _, key := range pressedKeys {
		if key == vxn {
			return
		}
	}
	c.pc += 2
}

// getDelayTimer - FX07 set value of VX to the current value of the delay timer
func (c *Chip8) getDelayTimer() {
	*c.getRegisterPointer(c.cir[0] & 0x0F) = c.delay
}

// setDelayTimer - FX15 set delay timer to value of VX
func (c *Chip8) setDelayTimer() {
	c.delay = *c.getRegisterPointer(c.cir[0] & 0x0F)
}

// setSoundTimer - FX18 set sound timer to the value of VX
func (c *Chip8) setSoundTimer() {
	c.sound = *c.getRegisterPointer(c.cir[0] & 0x0F)
}

// getPressedKey - FX0A blocks until a key is pressed. Stores that key's value in VX then continues.
func (c *Chip8) getPressedKey() {
	// TODO: On the original COSMAC VIP, the key was only registered when it was pressed and then released.

	pressedKeys := c.ui.GetPressedKeys()
	if len(pressedKeys) == 0 {
		// block
		c.pc -= 2
		return
	}

	*c.getRegisterPointer(c.cir[0] & 0x0F) = pressedKeys[0]
}

// addToIndexRegister - FX1E adds the value of VX to the index register and set VF accordingly if the index register
// "overflows" above 0x0FFF. Setting VF does not occur if DisableSetFlagOnIrOverflow is true.
func (c *Chip8) addToIndexRegister() {
	c.ir += uint16(*c.getRegisterPointer(c.cir[0] & 0x0F))

	if !c.DisableSetFlagOnIrOverflow {
		if c.ir > 0x0FFF {
			c.vf = 0x01
		} else {
			c.vf = 0x00
		}
	}
}

// getFontCharacter - FX29 set the index register to the address of the hex character in VX
func (c *Chip8) getFontCharacter() {
	c.ir = getFontCharacterLocation(*c.getRegisterPointer(c.cir[0] & 0x0F))
}

// convertToDecimal - FX33 take the value of VX, converts it to a denary number and the put each individual digit in the
// memory location specified by the index register + the digit number.
// Eg 0x9C -> 156 -> memory[ic] = 1, memory[ic+1] = 5, memory[ic+2] = 6
func (c *Chip8) convertToDecimal() {
	vxn := *c.getRegisterPointer(c.cir[0] & 0x0F)

	x := vxn % 10
	y := ((vxn - x) / 10) % 10
	z := (vxn - x - y*10) / 100

	c.memory[c.ir] = z
	c.memory[c.ir+1] = y
	c.memory[c.ir+2] = x
}

// storeMemory - FX55 store the value of each general purpose register from V0 to VX inclusive in consecutive memory
// addresses starting from the current value of the index register. If IncrementIndexRegisterOnLoadSave is true, the
// index register will be incremented as a result of this process. Else, a temporary variable will be used.
func (c *Chip8) storeMemory() {
	x := c.cir[0] & 0x0F
	for i := byte(0x00); i <= x; i += 1 {
		c.memory[c.ir + uint16(i)] = *c.getRegisterPointer(i)
	}
	if c.IncrementIndexRegisterOnLoadSave {
		c.ir += uint16(x)
	}
}

// loadMemory - FX65 loads the value of each general purpose register from V0 to VX inclusive from consecutive memory
// addresses starting from the current value of the index register. If IncrementIndexRegisterOnLoadSave is true, the
// index register will be incremented as a result of this process. Else, a temporary variable will be used.
func (c *Chip8) loadMemory() {
	x := c.cir[0] & 0x0F
	for i := byte(0x00); i <= x; i += 1 {
		*c.getRegisterPointer(i) = c.memory[c.ir + uint16(i)]
	}
	if c.IncrementIndexRegisterOnLoadSave {
		c.ir += uint16(x)
	}
}