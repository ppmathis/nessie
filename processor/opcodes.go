package processor

type Opcode uint8
type OpcodeHandler func(mode AddressingMode) (extraCycles Cycles)

func (c *CPU) setFlagsZN(result uint8) uint8 {
	c.Flags.Zero = result == 0
	c.Flags.Negative = (result>>7)&0x1 == 1
	return result
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
