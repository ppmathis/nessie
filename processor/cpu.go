package processor

import "fmt"

type Cycles uint64

type Register int
type Flag int
type FlagOrigin int

const (
	RegisterPC Register = iota
	RegisterS
	RegisterA
	RegisterX
	RegisterY
)

const (
	FlagCarry Flag = iota
	FlagZero
	FlagInterruptDisable
	FlagDecimal
	FlagOverflow
	FlagNegative
)

const (
	FlagOriginPHP FlagOrigin = iota
	FlagOriginBRK
	FlagOriginIRQ
	FlagOriginNMI
)

type Flags struct {
	Origin           FlagOrigin
	Carry            bool
	Zero             bool
	InterruptDisable bool
	Decimal          bool
	Overflow         bool
	Negative         bool
}

type Registers struct {
	PC uint16
	S  uint8
	A  uint8
	X  uint8
	Y  uint8
}

type CPU struct {
	TotalCycles Cycles
	Flags       Flags
	Registers   Registers
	Memory      Memory

	addressingHandlers AddressingHandlerTable
	instructions       InstructionTable
}

func NewCPU() (cpu *CPU) {
	cpu = &CPU{
		TotalCycles: 0,
		Flags:       Flags{},
		Registers:   Registers{},
		Memory:      NewBasicMemory(),
	}

	cpu.registerAddressingHandlers()
	cpu.registerInstructions()

	return
}

func (c *CPU) Execute() {
	opcode := Opcode(c.Memory.Peek(c.Registers.PC))
	c.Registers.PC++

	instruction, ok := c.instructions[opcode]
	if !ok {
		panic(fmt.Errorf("invalid opcode: 0x%02X", opcode))
	}

	cycles := instruction.Variant.StaticCycles
	cycles += instruction.Handler(instruction.Variant.AddressingMode)
	c.TotalCycles += cycles
}

func (c *CPU) GetRegister(register Register) uint16 {
	switch register {
	case RegisterPC:
		return c.Registers.PC
	case RegisterS:
		return uint16(c.Registers.S)
	case RegisterA:
		return uint16(c.Registers.A)
	case RegisterX:
		return uint16(c.Registers.X)
	case RegisterY:
		return uint16(c.Registers.Y)
	default:
		panic(fmt.Errorf("unknown cpu register: %d", register))
	}
}

func (c *CPU) GetRegisterName(register Register) string {
	switch register {
	case RegisterPC:
		return "PC"
	case RegisterS:
		return "S"
	case RegisterA:
		return "A"
	case RegisterX:
		return "X"
	case RegisterY:
		return "Y"
	default:
		panic(fmt.Errorf("unknown cpu register: %d", register))
	}
}

func (c *CPU) GetFlag(flag Flag) bool {
	switch flag {
	case FlagCarry:
		return c.Flags.Carry
	case FlagZero:
		return c.Flags.Zero
	case FlagInterruptDisable:
		return c.Flags.InterruptDisable
	case FlagDecimal:
		return c.Flags.Decimal
	case FlagOverflow:
		return c.Flags.Overflow
	case FlagNegative:
		return c.Flags.Negative
	default:
		panic(fmt.Errorf("unknown flag: %d", flag))
	}
}

func (c *CPU) GetFlagName(flag Flag) string {
	switch flag {
	case FlagCarry:
		return "Carry"
	case FlagZero:
		return "Zero"
	case FlagInterruptDisable:
		return "InterruptDisable"
	case FlagDecimal:
		return "Decimal"
	case FlagOverflow:
		return "Overflow"
	case FlagNegative:
		return "Negative"
	default:
		panic(fmt.Errorf("unknown flag: %d", flag))
	}
}
