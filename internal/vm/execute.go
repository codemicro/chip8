package vm

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

// setRegister - 6XNN set VX to NN
func (c *Chip8) setRegister() {
	nn := c.get8bitConstant()
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	*vx = nn
}

// addToRegister - 7XNN add NN to VX
func (c *Chip8) addToRegister() {
	nn := c.get8bitConstant()
	vx := c.getRegisterPointer(c.cir[0] & 0x0F)
	*vx += nn
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
