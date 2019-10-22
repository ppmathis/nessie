package processor

type Opcode uint8
type OpcodeHandler func(mode AddressingMode) (extraCycles Cycles)

func (c *CPU) setFlagsZN(result uint8) uint8 {
	c.Flags.Zero = result == 0
	c.Flags.Negative = (result>>7)&0x1 == 1
	return result
}

func (c *CPU) addition(value uint8) {
	previousA := uint16(c.Registers.A)
	summand := uint16(value)

	sum := previousA + summand
	if c.Flags.Carry {
		sum++
	}
	result := sum & 0xFF

	c.Registers.A = uint8(result)
	c.Flags.Carry = (sum & 0x100) == 0x100
	c.Flags.Zero = result == 0
	c.Flags.Overflow = ((previousA ^ result) & (summand ^ result) & 0x80) == 0x80
	c.Flags.Negative = (sum & 0x80) == 0x80
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

func (c *CPU) opAND(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := c.setFlagsZN(c.Registers.A & value)
	c.Registers.A = result
	return
}

func (c *CPU) opORA(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := c.setFlagsZN(c.Registers.A | value)
	c.Registers.A = result
	return
}

func (c *CPU) opEOR(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := c.setFlagsZN(c.Registers.A ^ value)
	c.Registers.A = result
	return
}

func (c *CPU) opINC(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := c.setFlagsZN(value + 1)
	c.Memory.Poke(address, result)
	return
}

func (c *CPU) opINX(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.X + 1)
	c.Registers.X = result
	return
}

func (c *CPU) opINY(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.Y + 1)
	c.Registers.Y = result
	return
}

func (c *CPU) opDEC(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := c.setFlagsZN(value - 1)
	c.Memory.Poke(address, result)
	return
}

func (c *CPU) opDEX(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.X - 1)
	c.Registers.X = result
	return
}

func (c *CPU) opDEY(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.Y - 1)
	c.Registers.Y = result
	return
}

func (c *CPU) opCMP(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.setFlagsZN(c.Registers.A - value)
	c.Flags.Carry = c.Registers.A >= value
	return
}

func (c *CPU) opCPX(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.setFlagsZN(c.Registers.X - value)
	c.Flags.Carry = c.Registers.X >= value
	return
}

func (c *CPU) opCPY(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	c.setFlagsZN(c.Registers.Y - value)
	c.Flags.Carry = c.Registers.Y >= value
	return
}

func (c *CPU) opTAX(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.A)
	c.Registers.X = result
	return
}

func (c *CPU) opTXA(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.X)
	c.Registers.A = result
	return
}

func (c *CPU) opTAY(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.A)
	c.Registers.Y = result
	return
}

func (c *CPU) opTYA(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.Y)
	c.Registers.A = result
	return
}

func (c *CPU) opTSX(mode AddressingMode) (extraCycles Cycles) {
	result := c.setFlagsZN(c.Registers.S)
	c.Registers.X = result
	return
}

func (c *CPU) opTXS(mode AddressingMode) (extraCycles Cycles) {
	c.Registers.S = c.Registers.X
	return
}

func (c *CPU) opBCS(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Flags.Carry, target)
}

func (c *CPU) opBCC(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(!c.Flags.Carry, target)
}

func (c *CPU) opBEQ(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Flags.Zero, target)
}

func (c *CPU) opBNE(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(!c.Flags.Zero, target)
}

func (c *CPU) opBMI(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Flags.Negative, target)
}

func (c *CPU) opBPL(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(!c.Flags.Negative, target)
}

func (c *CPU) opBVS(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(c.Flags.Overflow, target)
}

func (c *CPU) opBVC(mode AddressingMode) (extraCycles Cycles) {
	target, _ := c.lookupAddress(mode)
	return c.branch(!c.Flags.Overflow, target)
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
	c.Flags.Carry = false
	return
}

func (c *CPU) opSEC(mode AddressingMode) (extraCycles Cycles) {
	c.Flags.Carry = true
	return
}

func (c *CPU) opCLD(mode AddressingMode) (extraCycles Cycles) {
	c.Flags.Decimal = false
	return
}

func (c *CPU) opSED(mode AddressingMode) (extraCycles Cycles) {
	c.Flags.Decimal = true
	return
}

func (c *CPU) opCLI(mode AddressingMode) (extraCycles Cycles) {
	c.Flags.InterruptDisable = false
	return
}

func (c *CPU) opSEI(mode AddressingMode) (extraCycles Cycles) {
	c.Flags.InterruptDisable = true
	return
}

func (c *CPU) opCLV(mode AddressingMode) (extraCycles Cycles) {
	c.Flags.Overflow = false
	return
}

func (c *CPU) opLDA(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := c.setFlagsZN(value)
	c.Registers.A = result
	return
}

func (c *CPU) opLDX(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := c.setFlagsZN(value)
	c.Registers.X = result
	return
}

func (c *CPU) opLDY(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := c.setFlagsZN(value)
	c.Registers.Y = result
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
	result := c.setFlagsZN(c.Pop())
	c.Registers.A = result
	return
}

func (c *CPU) opPHP(mode AddressingMode) (extraCycles Cycles) {
	c.Push(c.FlagsBinary(FlagOriginPHP))
	return
}

func (c *CPU) opPLP(mode AddressingMode) (extraCycles Cycles) {
	value := c.Pop()
	c.Flags.Carry = (value>>0)&0x1 == 0x1
	c.Flags.Zero = (value>>1)&0x1 == 0x1
	c.Flags.InterruptDisable = (value>>2)&0x1 == 0x1
	c.Flags.Decimal = (value>>3)&0x1 == 0x1
	c.Flags.Overflow = (value>>6)&0x1 == 0x1
	c.Flags.Negative = (value>>7)&0x1 == 0x1

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
	c.Flags.InterruptDisable = true
	c.Registers.PC = c.Memory.Peek16(ResetVector)
	c.Push(c.FlagsBinary(FlagOriginBRK))
	return
}

func (c *CPU) opBIT(mode AddressingMode) (extraCycles Cycles) {
	address, _ := c.lookupAddress(mode)
	value := c.Memory.Peek(address)

	result := value & c.Registers.A
	c.Flags.Zero = result == 0
	c.Flags.Overflow = (value>>6)&0x1 == 0x1
	c.Flags.Negative = (value>>7)&0x1 == 0x1
	return
}

func (c *CPU) opROL(mode AddressingMode) (extraCycles Cycles) {
	var value uint8
	var address uint16

	if mode == Accumulator {
		value = c.Registers.A
	} else {
		address, _ = c.lookupAddress(mode)
		value = c.Memory.Peek(address)
	}

	previousCarry := (value>>7)&0x1 == 0x1
	value <<= 1
	if c.Flags.Carry {
		value |= 1 << 0
	} else {
		value &^= 1 << 0
	}

	if mode == Accumulator {
		c.Registers.A = value
	} else {
		c.Memory.Poke(address, value)
	}

	c.Flags.Carry = previousCarry
	c.Flags.Zero = c.Registers.A == 0
	c.Flags.Negative = (value & 0x80) == 0x80

	return
}

func (c *CPU) opROR(mode AddressingMode) (extraCycles Cycles) {
	var value uint8
	var address uint16

	if mode == Accumulator {
		value = c.Registers.A
	} else {
		address, _ = c.lookupAddress(mode)
		value = c.Memory.Peek(address)
	}

	previousCarry := (value>>0)&0x1 == 0x1
	value >>= 1
	if c.Flags.Carry {
		value |= 1 << 7
	} else {
		value &^= 1 << 7
	}

	if mode == Accumulator {
		c.Registers.A = value
	} else {
		c.Memory.Poke(address, value)
	}

	c.Flags.Carry = previousCarry
	c.Flags.Zero = c.Registers.A == 0
	c.Flags.Negative = (value & 0x80) == 0x80

	return
}

func (c *CPU) opASL(mode AddressingMode) (extraCycles Cycles) {
	var value uint8
	var address uint16

	if mode == Accumulator {
		value = c.Registers.A
	} else {
		address, _ = c.lookupAddress(mode)
		value = c.Memory.Peek(address)
	}

	c.Flags.Carry = (value>>7)&0x1 == 0x1
	value = c.setFlagsZN(value << 1)

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

	if mode == Accumulator {
		value = c.Registers.A
	} else {
		address, _ = c.lookupAddress(mode)
		value = c.Memory.Peek(address)
	}

	c.Flags.Carry = (value & 0x1) == 0x1
	value = c.setFlagsZN(value >> 1)

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
