package processor

import "fmt"

type Cycles uint64

type Register int
type Status uint8

const ResetVector = 0xFFFE

const (
	FlagCarry Status = 1 << iota
	FlagZero
	FlagInterruptDisable
	FlagDecimal
	FlagBreak
	FlagUnused
	FlagOverflow
	FlagNegative
)

type Registers struct {
	PC uint16
	P  Status
	S  uint8
	A  uint8
	X  uint8
	Y  uint8
}

type CPU struct {
	Debug       bool
	Halted      bool
	TotalCycles Cycles
	Registers   Registers
	Memory      *MappedMemory

	lastDisassembly    string
	addressingHandlers AddressingHandlerTable
	instructions       InstructionTable
}

func NewCPU() (cpu *CPU) {
	cpu = &CPU{
		Halted:      false,
		TotalCycles: 0,
		Registers: Registers{
			PC: 0xFFFC,
			P:  FlagInterruptDisable | FlagBreak | FlagUnused,
			S:  0xFD,
		},
		Memory: NewMappedMemory(NewBasicMemory()),
	}

	cpu.registerAddressingHandlers()
	cpu.registerInstructions()

	return
}

func (c *CPU) Execute() {
	if c.Halted {
		panic("cpu is halted")
	}

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

func (c *CPU) Disassembly() string {
	return c.lastDisassembly
}
