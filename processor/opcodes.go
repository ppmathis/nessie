package processor

type Opcode uint8
type OpcodeHandler func(mode AddressingMode) (extraCycles Cycles)

func (c *CPU) setFlagsZN(result uint8) uint8 {
	c.Flags.Zero = result == 0
	c.Flags.Negative = (result>>7)&0x1 == 1
	return result
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
	var value uint8

	if c.Flags.Carry {
		value |= 0x1 << 0
	}
	if c.Flags.Zero {
		value |= 0x1 << 1
	}
	if c.Flags.InterruptDisable {
		value |= 0x1 << 2
	}
	if c.Flags.Decimal {
		value |= 0x1 << 3
	}
	if c.Flags.Overflow {
		value |= 0x1 << 6
	}
	if c.Flags.Negative {
		value |= 0x1 << 7
	}

	switch c.Flags.Origin {
	case FlagOriginPHP:
		value |= 0x3 << 4
	case FlagOriginBRK:
		value |= 0x3 << 4
		c.Flags.InterruptDisable = true
	case FlagOriginIRQ, FlagOriginNMI:
		value |= 0x2 << 4
		c.Flags.InterruptDisable = true
	}

	c.Push(value)
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
	return
}

func (c *CPU) opRTI(mode AddressingMode) (extraCycles Cycles) {
	c.opPLP(Implicit)
	c.Registers.PC = c.Pop16()
	return
}

func (c *CPU) opBRK(mode AddressingMode) (extraCycles Cycles) {
	c.Push16(c.Registers.PC)
	c.opPHP(Implicit)
	c.Registers.PC = c.Memory.Peek16(ResetVector)
	c.Flags.Origin = FlagOriginBRK
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
