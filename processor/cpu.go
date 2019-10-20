package processor

import "fmt"

type Cycles uint64

type Register int
type Flag int
type FlagOrigin int

const ResetVector = 0xFFFE

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
	Debug       bool
	TotalCycles Cycles
	Flags       Flags
	Registers   Registers
	Memory      *MappedMemory

	addressingHandlers AddressingHandlerTable
	instructions       InstructionTable
}

func NewCPU() (cpu *CPU) {
	cpu = &CPU{
		TotalCycles: 0,
		Flags:       Flags{},
		Registers: Registers{
			PC: 0xFFFC,
			S:  0xFD,
		},
		Memory: NewMappedMemory(NewBasicMemory()),
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

	if c.Debug {
		var assembly string

		switch instruction.Variant.AddressingMode {
		case Implicit, Accumulator:
			assembly = fmt.Sprintf("%02X      ",
				c.Memory.Peek(c.Registers.PC-1),
			)
		case Immediate, ZeroPage, ZeroPageX, ZeroPageY, IndirectX, IndirectY, Relative:
			assembly = fmt.Sprintf("%02X %02X   ",
				c.Memory.Peek(c.Registers.PC-1),
				c.Memory.Peek(c.Registers.PC),
			)
		case Absolute, AbsoluteX, AbsoluteY, Indirect:
			assembly = fmt.Sprintf("%02X %02X %02X",
				c.Memory.Peek(c.Registers.PC-1),
				c.Memory.Peek(c.Registers.PC),
				c.Memory.Peek(c.Registers.PC+1),
			)
		}

		fmt.Printf("[0x%04X] %s - %s - A:%02X X:%02X Y:%02X S:%02X CYC:%d\n",
			c.Registers.PC-1, assembly, instruction.Mnemonic,
			c.Registers.A, c.Registers.X, c.Registers.Y, c.Registers.S,
			c.TotalCycles,
		)
	}

	cycles := instruction.Variant.StaticCycles
	cycles += instruction.Handler(instruction.Variant.AddressingMode)
	c.TotalCycles += cycles
}

func (c *CPU) Push(value uint8) {
	address := 0x0100 | uint16(c.Registers.S)
	c.Memory.Poke(address, value)
	c.Registers.S--
}

func (c *CPU) Push16(value uint16) {
	c.Push(uint8(value >> 8))
	c.Push(uint8(value))
}

func (c *CPU) Pop() (value uint8) {
	c.Registers.S++
	address := 0x0100 | uint16(c.Registers.S)
	value = c.Memory.Peek(address)
	return
}

func (c *CPU) Pop16() (value uint16) {
	value = uint16(c.Pop())
	value |= uint16(c.Pop()) << 8
	return
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
