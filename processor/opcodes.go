package processor

type Opcode uint8
type OpcodeHandler func(mode AddressingMode) (extraCycles Cycles)

func (c *CPU) setZero() {
	c.Flags.Zero = c.Registers.A == 0
}

func (c *CPU) setNegative() {
	c.Flags.Negative = (c.Registers.A>>7)&0x1 == 1
}

func (c *CPU) opAND(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)
	c.Registers.A &= value

	c.setZero()
	c.setNegative()
	return
}

func (c *CPU) opORA(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)
	c.Registers.A |= value

	c.setZero()
	c.setNegative()
	return
}

func (c *CPU) opEOR(mode AddressingMode) (extraCycles Cycles) {
	address, extraCycles := c.lookupAddress(mode)
	value := c.Memory.Peek(address)
	c.Registers.A ^= value

	c.setZero()
	c.setNegative()
	return
}
