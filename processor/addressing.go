package processor

import "fmt"

type AddressingMode int
type AddressingHandler func() (address uint16, extraCycles Cycles)
type AddressingHandlerTable map[AddressingMode]AddressingHandler

const (
	Implicit    AddressingMode = iota
	Accumulator AddressingMode = iota
	Immediate   AddressingMode = iota
	ZeroPage    AddressingMode = iota
	ZeroPageX   AddressingMode = iota
	ZeroPageY   AddressingMode = iota
	Relative    AddressingMode = iota
	Absolute    AddressingMode = iota
	AbsoluteX   AddressingMode = iota
	AbsoluteY   AddressingMode = iota
	Indirect    AddressingMode = iota
	IndirectX   AddressingMode = iota
	IndirectY   AddressingMode = iota
)

func (c *CPU) amImplicit() (address uint16, extraCycles Cycles) {
	panic(fmt.Errorf("unable to lookup implicit addressing mode"))
}

func (c *CPU) amAccumulator() (address uint16, extraCycles Cycles) {
	panic(fmt.Errorf("unable to lookup accumulator addressing mode"))
}

func (c *CPU) amImmediate() (address uint16, extraCycles Cycles) {
	address = c.Registers.PC
	c.Registers.PC++
	return
}

func (c *CPU) amZeroPage() (address uint16, extraCycles Cycles) {
	address = uint16(c.Memory.Peek(c.Registers.PC))
	c.Registers.PC++
	return
}

func (c *CPU) amZeroPageX() (address uint16, extraCycles Cycles) {
	address = uint16(c.Memory.Peek(c.Registers.PC)+c.Registers.X) & 0xFF
	c.Registers.PC++
	return
}

func (c *CPU) amZeroPageY() (address uint16, extraCycles Cycles) {
	address = uint16(c.Memory.Peek(c.Registers.PC)+c.Registers.Y) & 0xFF
	c.Registers.PC++
	return
}

func (c *CPU) amRelative() (address uint16, extraCycles Cycles) {
	value := c.Memory.Peek(c.Registers.PC)
	c.Registers.PC++

	var offset int16
	if value > 0x7F {
		offset = -(0x100 - int16(value))
	} else {
		offset = int16(value & 0x7F)
	}

	address = uint16(int16(c.Registers.PC) + offset)
	return
}

func (c *CPU) amAbsolute() (address uint16, extraCycles Cycles) {
	address = c.Memory.Peek16(c.Registers.PC)
	c.Registers.PC += 2
	return
}

func (c *CPU) amAbsoluteX() (address uint16, extraCycles Cycles) {
	addressPtr := c.Memory.Peek16(c.Registers.PC)
	c.Registers.PC += 2

	address = addressPtr + uint16(c.Registers.X)
	if !SamePage(addressPtr, address) {
		extraCycles = 1
	}

	return
}

func (c *CPU) amAbsoluteY() (address uint16, extraCycles Cycles) {
	addressPtr := c.Memory.Peek16(c.Registers.PC)
	c.Registers.PC += 2

	address = addressPtr + uint16(c.Registers.Y)
	if !SamePage(addressPtr, address) {
		extraCycles = 1
	}

	return
}

func (c *CPU) amIndirect() (address uint16, extraCycles Cycles) {
	addressPtr := c.Memory.Peek16(c.Registers.PC)
	c.Registers.PC += 2

	wrappedPtr := (addressPtr & 0xFF00) | ((addressPtr + 1) & 0x00FF)
	addressLow := uint16(c.Memory.Peek(addressPtr))
	addressHigh := uint16(c.Memory.Peek(wrappedPtr))
	address = addressLow | addressHigh<<8

	return
}

func (c *CPU) amIndirectX() (address uint16, extraCycles Cycles) {
	addressPtr := uint16(c.Memory.Peek(c.Registers.PC) + c.Registers.X)
	c.Registers.PC++

	addressLow := uint16(c.Memory.Peek(addressPtr))
	addressHigh := uint16(c.Memory.Peek((addressPtr + 1) & 0xFF))
	address = addressLow | (addressHigh << 8)

	return
}

func (c *CPU) amIndirectY() (address uint16, extraCycles Cycles) {
	addressPtr := uint16(c.Memory.Peek(c.Registers.PC))
	c.Registers.PC++

	baseLow := uint16(c.Memory.Peek(addressPtr))
	baseHigh := uint16(c.Memory.Peek((addressPtr + 1) & 0xFF))
	baseAddress := baseLow | (baseHigh << 8)
	address = baseAddress + uint16(c.Registers.Y)

	if !SamePage(baseAddress, address) {
		extraCycles = 1
	}

	return
}

func (c *CPU) registerAddressingHandlers() {
	c.addressingHandlers = AddressingHandlerTable{
		Implicit:    c.amImplicit,
		Accumulator: c.amAccumulator,
		Immediate:   c.amImmediate,
		ZeroPage:    c.amZeroPage,
		ZeroPageX:   c.amZeroPageX,
		ZeroPageY:   c.amZeroPageY,
		Relative:    c.amRelative,
		Absolute:    c.amAbsolute,
		AbsoluteX:   c.amAbsoluteX,
		AbsoluteY:   c.amAbsoluteY,
		Indirect:    c.amIndirect,
		IndirectX:   c.amIndirectX,
		IndirectY:   c.amIndirectY,
	}
}

func (c *CPU) lookupAddress(mode AddressingMode) (address uint16, extraCycles Cycles) {
	address, extraCycles = c.addressingHandlers[mode]()
	return
}

func AddressingModeName(mode AddressingMode) string {
	switch mode {
	case Implicit:
		return "Implicit"
	case Accumulator:
		return "Accumulator"
	case Immediate:
		return "Immediate"
	case ZeroPage:
		return "ZeroPage"
	case ZeroPageX:
		return "ZeroPageX"
	case ZeroPageY:
		return "ZeroPageY"
	case Relative:
		return "Relative"
	case Absolute:
		return "Absolute"
	case AbsoluteX:
		return "AbsoluteX"
	case AbsoluteY:
		return "AbsoluteY"
	case Indirect:
		return "Indirect"
	case IndirectX:
		return "IndirectX"
	case IndirectY:
		return "IndirectY"
	default:
		return "<unknown>"
	}
}
