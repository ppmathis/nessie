package processor

type Opcode uint8
type OpcodeHandler func(mode AddressingMode) (extraCycles Cycles)

func (c *CPU) setZeroNegative(value uint8) {
	if value == 0 {
		c.Registers.P |= FlagZero
	} else {
		c.Registers.P &^= FlagZero
	}

	if value&0x80 == 0x80 {
		c.Registers.P |= FlagNegative
	} else {
		c.Registers.P &^= FlagNegative
	}
}

func (c *CPU) setCarry(value uint16) {
	if value>>8 != 0 {
		c.Registers.P |= FlagCarry
	} else {
		c.Registers.P &^= FlagCarry
	}
}

func (c *CPU) setOverflow(value1 uint16, value2 uint16, result uint16) {
	if ((value1 ^ result) & (value2 ^ result) & 0x80) == 0x80 {
		c.Registers.P |= FlagOverflow
	} else {
		c.Registers.P &^= FlagOverflow
	}
}

func (c *CPU) addition(value uint8) {
	value1 := uint16(c.Registers.A)
	value2 := uint16(value)

	sum := value1 + value2
	if c.Registers.P&FlagCarry == FlagCarry {
		sum++
	}
	result := uint8(sum & 0xFF)

	c.Registers.A = result
	c.setCarry(sum)
	c.setZeroNegative(result)
	c.setOverflow(value1, value2, uint16(result))
}

func (c *CPU) branch(condition bool, target uint16) (extraCycles Cycles) {
	if !condition {
		return
	}

	if !SamePage(c.Registers.PC, target) {
		extraCycles += 2
	} else {
		extraCycles += 1
	}

	c.Registers.PC = target
	return
}

func (c *CPU) compare(value1 uint8, value2 uint8) {
	result := uint16(value1) + (uint16(value2) ^ 0xFF + 1)
	c.setCarry(result)
	c.setZeroNegative(uint8(result))
}

func (c *CPU) opAND(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.Registers.A &= value
	c.setZeroNegative(c.Registers.A)
	return
}

func (c *CPU) opORA(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.Registers.A |= value
	c.setZeroNegative(c.Registers.A)
	return
}

func (c *CPU) opEOR(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.Registers.A ^= value
	c.setZeroNegative(c.Registers.A)
	return
}

func (c *CPU) opINC(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	result := c.Memory.Peek(address) + 1

	c.Memory.Poke(address, result)
	c.setZeroNegative(result)
	return
}

func (c *CPU) opINX(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.X++
	c.setZeroNegative(c.Registers.X)
	return
}

func (c *CPU) opINY(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.Y++
	c.setZeroNegative(c.Registers.Y)
	return
}

func (c *CPU) opDEC(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	result := c.Memory.Peek(address) - 1

	c.Memory.Poke(address, result)
	c.setZeroNegative(result)
	return
}

func (c *CPU) opDEX(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.X--
	c.setZeroNegative(c.Registers.X)
	return
}

func (c *CPU) opDEY(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.Y--
	c.setZeroNegative(c.Registers.Y)
	return
}

func (c *CPU) opCMP(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.compare(c.Registers.A, value)
	return
}

func (c *CPU) opCPX(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.compare(c.Registers.X, value)
	return
}

func (c *CPU) opCPY(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.compare(c.Registers.Y, value)
	return
}

func (c *CPU) opTAX(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.X = c.Registers.A
	c.setZeroNegative(c.Registers.X)
	return
}

func (c *CPU) opTXA(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.A = c.Registers.X
	c.setZeroNegative(c.Registers.A)
	return
}

func (c *CPU) opTAY(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.Y = c.Registers.A
	c.setZeroNegative(c.Registers.Y)
	return
}

func (c *CPU) opTYA(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.A = c.Registers.Y
	c.setZeroNegative(c.Registers.A)
	return
}

func (c *CPU) opTSX(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.X = c.Registers.S
	c.setZeroNegative(c.Registers.X)
	return
}

func (c *CPU) opTXS(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.S = c.Registers.X
	return
}

func (c *CPU) opBCS(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Registers.P&FlagCarry == FlagCarry, target)
}

func (c *CPU) opBCC(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Registers.P&FlagCarry != FlagCarry, target)
}

func (c *CPU) opBEQ(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Registers.P&FlagZero == FlagZero, target)
}

func (c *CPU) opBNE(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Registers.P&FlagZero != FlagZero, target)
}

func (c *CPU) opBMI(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Registers.P&FlagNegative == FlagNegative, target)
}

func (c *CPU) opBPL(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Registers.P&FlagNegative != FlagNegative, target)
}

func (c *CPU) opBVS(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Registers.P&FlagOverflow == FlagOverflow, target)
}

func (c *CPU) opBVC(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Registers.P&FlagOverflow != FlagOverflow, target)
}

func (c *CPU) opSTA(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	c.Memory.Poke(address, c.Registers.A)
	return
}

func (c *CPU) opSTX(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	c.Memory.Poke(address, c.Registers.X)
	return
}

func (c *CPU) opSTY(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	c.Memory.Poke(address, c.Registers.Y)
	return
}

func (c *CPU) opCLC(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.P &^= FlagCarry
	return
}

func (c *CPU) opSEC(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.P |= FlagCarry
	return
}

func (c *CPU) opCLD(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.P &^= FlagDecimal
	return
}

func (c *CPU) opSED(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.P |= FlagDecimal
	return
}

func (c *CPU) opCLI(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.P &^= FlagInterruptDisable
	return
}

func (c *CPU) opSEI(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.P |= FlagInterruptDisable
	return
}

func (c *CPU) opCLV(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.P &^= FlagOverflow
	return
}

func (c *CPU) opLDA(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.Registers.A = value
	c.setZeroNegative(c.Registers.A)
	return
}

func (c *CPU) opLDX(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.Registers.X = value
	c.setZeroNegative(c.Registers.X)
	return
}

func (c *CPU) opLDY(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.Registers.Y = value
	c.setZeroNegative(c.Registers.Y)
	return
}

func (c *CPU) opJMP(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	c.Registers.PC = target
	return
}

func (c *CPU) opJSR(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	c.Push16(c.Registers.PC - 1)
	c.Registers.PC = target
	return
}

func (c *CPU) opRTS(mode AddressingMode) (extraCycles Cycles) {
	target := c.Pop16() + 1
	c.Registers.PC = target
	return
}

func (c *CPU) opPHA(mode AddressingMode) (extraCycles Cycles) {
	c.Push(c.Registers.A)
	return
}

func (c *CPU) opPLA(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.A = c.Pop()
	c.setZeroNegative(c.Registers.A)
	return
}

func (c *CPU) opPHP(mode AddressingMode) (extraCycles Cycles) {
	c.Push(uint8(c.Registers.P | FlagBreak))
	return
}

func (c *CPU) opPLP(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.P = Status(c.Pop())
	c.Registers.P &^= FlagBreak
	c.Registers.P |= FlagUnused
	return
}

func (c *CPU) opNOP(mode AddressingMode) (extraCycles Cycles) {
	if mode != Implicit {
		_, extraCycles = c.lookupAddress(mode)
	}
	return
}

func (c *CPU) opRTI(mode AddressingMode) (extraCycles Cycles) {
	c.opPLP(Implicit)
	c.Registers.PC = c.Pop16()
	return
}

func (c *CPU) opBRK(mode AddressingMode) (extraCycles Cycles) {
	c.Push16(c.Registers.PC)
	c.Registers.P |= FlagInterruptDisable
	c.Registers.PC = c.Memory.Peek16(ResetVector)
	c.Push(uint8(c.Registers.P | FlagBreak))
	return
}

func (c *CPU) opBIT(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	// Set zero flag if result of AND equals zero
	result := c.Registers.A & value
	if result == 0 {
		c.Registers.P |= FlagZero
	} else {
		c.Registers.P &^= FlagZero
	}

	// Set overflow flag to value of 6th bit
	if value&0x40 == 0x40 {
		c.Registers.P |= FlagOverflow
	} else {
		c.Registers.P &^= FlagOverflow
	}

	// Set negative flag to value of 7th bit
	if value&0x80 == 0x80 {
		c.Registers.P |= FlagNegative
	} else {
		c.Registers.P &^= FlagNegative
	}

	return
}

func (c *CPU) opROL(mode AddressingMode) (extraCycles Cycles) {
	var value uint8
	var address uint16

	// Load value from accumulator or memory
	if mode == Accumulator {
		value = c.Registers.A
	} else {
		address, _ = c.lookupAddress(mode)
		value = c.Memory.Peek(address)
	}

	// Re-use bit 7 as new carry
	previousCarry := (c.Registers.P & FlagCarry) == FlagCarry
	if value&0x80 == 0x80 {
		c.Registers.P |= FlagCarry
	} else {
		c.Registers.P &^= FlagCarry
	}

	// Shift value to the left and re-add previous carry as bit 0
	value <<= 1
	if previousCarry {
		value |= 0x1
	} else {
		value &^= 0x1
	}
	c.setZeroNegative(value)

	// Store value back into accumulator or memory
	if mode == Accumulator {
		c.Registers.A = value
	} else {
		c.Memory.Poke(address, value)
	}

	return
}

func (c *CPU) opROR(mode AddressingMode) (extraCycles Cycles) {
	var value uint8
	var address uint16

	// Load value from accumulator or memory
	if mode == Accumulator {
		value = c.Registers.A
	} else {
		address, _ = c.lookupAddress(mode)
		value = c.Memory.Peek(address)
	}

	// Re-use bit 0 as new carry
	previousCarry := (c.Registers.P & FlagCarry) == FlagCarry
	if value&0x1 == 0x1 {
		c.Registers.P |= FlagCarry
	} else {
		c.Registers.P &^= FlagCarry
	}

	// Shift value to the left and re-add previous carry as bit 7
	value >>= 1
	if previousCarry {
		value |= 0x80
	} else {
		value &^= 0x80
	}
	c.setZeroNegative(value)

	// Store value back into accumulator or memory
	if mode == Accumulator {
		c.Registers.A = value
	} else {
		c.Memory.Poke(address, value)
	}

	return
}

func (c *CPU) opASL(mode AddressingMode) (extraCycles Cycles) {
	var value uint8
	var address uint16

	// Load value from accumulator or memory
	if mode == Accumulator {
		value = c.Registers.A
	} else {
		address, _ = c.lookupAddress(mode)
		value = c.Memory.Peek(address)
	}

	// Re-use bit 7 as new carry
	if value&0x80 == 0x80 {
		c.Registers.P |= FlagCarry
	} else {
		c.Registers.P &^= FlagCarry
	}

	// Shift value to the left
	value <<= 1
	c.setZeroNegative(value)

	// Store value back into accumulator or memory
	if mode == Accumulator {
		c.Registers.A = value
	} else {
		c.Memory.Poke(address, value)
	}

	return
}

func (c *CPU) opLSR(mode AddressingMode) (extraCycles Cycles) {
	var value uint8
	var address uint16

	// Load value from accumulator or memory
	if mode == Accumulator {
		value = c.Registers.A
	} else {
		address, _ = c.lookupAddress(mode)
		value = c.Memory.Peek(address)
	}

	// Re-use bit 0 as new carry
	if value&0x1 == 0x1 {
		c.Registers.P |= FlagCarry
	} else {
		c.Registers.P &^= FlagCarry
	}

	// Shift value to the right
	value >>= 1
	c.setZeroNegative(value)

	// Store value back into accumulator or memory
	if mode == Accumulator {
		c.Registers.A = value
	} else {
		c.Memory.Poke(address, value)
	}

	return
}

func (c *CPU) opADC(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	c.addition(c.Memory.Peek(address))
	return
}

func (c *CPU) opSBC(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	c.addition(^c.Memory.Peek(address))
	return
}

func (c *CPU) opLAX(mode AddressingMode) (extraCycles Cycles) {
	previousPC := c.Registers.PC
	extraCycles = c.opLDA(mode) // Re-use extra cycles from LDA opcode, as LAX behaves the same
	c.Registers.PC = previousPC
	_ = c.opLDX(mode)
	return
}

func (c *CPU) opSAX(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	c.Memory.Poke(address, c.Registers.A&c.Registers.X)
	return
}

func (c *CPU) opDCP(mode AddressingMode) (extraCycles Cycles) {
	previousPC := c.Registers.PC
	_ = c.opDEC(mode)
	c.Registers.PC = previousPC
	_ = c.opCMP(mode)
	return
}

func (c *CPU) opISC(mode AddressingMode) (extraCycles Cycles) {
	previousPC := c.Registers.PC
	_ = c.opINC(mode)
	c.Registers.PC = previousPC
	_ = c.opSBC(mode)
	return
}

func (c *CPU) opSLO(mode AddressingMode) (extraCycles Cycles) {
	previousPC := c.Registers.PC
	_ = c.opASL(mode)
	c.Registers.PC = previousPC
	_ = c.opORA(mode)
	return
}

func (c *CPU) opRLA(mode AddressingMode) (extraCycles Cycles) {
	previousPC := c.Registers.PC
	_ = c.opROL(mode)
	c.Registers.PC = previousPC
	_ = c.opAND(mode)
	return
}

func (c *CPU) opSRE(mode AddressingMode) (extraCycles Cycles) {
	previousPC := c.Registers.PC
	_ = c.opLSR(mode)
	c.Registers.PC = previousPC
	_ = c.opEOR(mode)
	return
}

func (c *CPU) opRRA(mode AddressingMode) (extraCycles Cycles) {
	previousPC := c.Registers.PC
	_ = c.opROR(mode)
	c.Registers.PC = previousPC
	_ = c.opADC(mode)
	return
}

func (c *CPU) opKIL(mode AddressingMode) (extraCycles Cycles) {
	c.Halted = true
	return
}
