package processor

import "fmt"

type InstructionTable map[Opcode]Instruction
type InstructionVariant struct {
	Opcode         Opcode
	AddressingMode AddressingMode
	StaticCycles   Cycles
}
type Instruction struct {
	Mnemonic string
	Handler  OpcodeHandler
	Variant  InstructionVariant
}

func (c *CPU) registerInstructions() {
	c.instructions = InstructionTable{}

	// AND - Logical AND
	c.instructions.registerVariants("AND", c.opAND,
		InstructionVariant{0x29, Immediate, 2},
		InstructionVariant{0x25, ZeroPage, 3},
		InstructionVariant{0x35, ZeroPageX, 4},
		InstructionVariant{0x2D, Absolute, 4},
		InstructionVariant{0x3D, AbsoluteX, 4},
		InstructionVariant{0x39, AbsoluteY, 4},
		InstructionVariant{0x21, IndirectX, 6},
		InstructionVariant{0x31, IndirectY, 5},
	)

	// ORA - Logical Inclusive OR
	c.instructions.registerVariants("ORA", c.opORA,
		InstructionVariant{0x09, Immediate, 2},
		InstructionVariant{0x05, ZeroPage, 3},
		InstructionVariant{0x15, ZeroPageX, 4},
		InstructionVariant{0x0D, Absolute, 4},
		InstructionVariant{0x1D, AbsoluteX, 4},
		InstructionVariant{0x19, AbsoluteY, 4},
		InstructionVariant{0x01, IndirectX, 6},
		InstructionVariant{0x11, IndirectY, 5},
	)

	// EOR - Exclusive OR
	c.instructions.registerVariants("EOR", c.opEOR,
		InstructionVariant{0x49, Immediate, 2},
		InstructionVariant{0x45, ZeroPage, 3},
		InstructionVariant{0x55, ZeroPageX, 4},
		InstructionVariant{0x4D, Absolute, 4},
		InstructionVariant{0x5D, AbsoluteX, 4},
		InstructionVariant{0x59, AbsoluteY, 4},
		InstructionVariant{0x41, IndirectX, 6},
		InstructionVariant{0x51, IndirectY, 5},
	)

	// INC - Increment Memory
	c.instructions.registerVariants("INC", c.opINC,
		InstructionVariant{0xE6, ZeroPage, 5},
		InstructionVariant{0xF6, ZeroPageX, 6},
		InstructionVariant{0xEE, Absolute, 6},
		InstructionVariant{0xFE, AbsoluteX, 7},
	)

	// INX - Increment X Register
	c.instructions.registerVariants("INX", c.opINX,
		InstructionVariant{0xE8, Implicit, 2},
	)

	// INY - Increment Y Register
	c.instructions.registerVariants("INY", c.opINY,
		InstructionVariant{0xC8, Implicit, 2},
	)

	// DEC - Decrement Memory
	c.instructions.registerVariants("DEC", c.opDEC,
		InstructionVariant{0xC6, ZeroPage, 5},
		InstructionVariant{0xD6, ZeroPageX, 6},
		InstructionVariant{0xCE, Absolute, 6},
		InstructionVariant{0xDE, AbsoluteX, 7},
	)

	// DEX - Decrement X Register
	c.instructions.registerVariants("DEX", c.opDEX,
		InstructionVariant{0xCA, Implicit, 2},
	)

	// DEY - Decrement Y Register
	c.instructions.registerVariants("DEY", c.opDEY,
		InstructionVariant{0x88, Implicit, 2},
	)

	// CMP - Compare
	c.instructions.registerVariants("CMP", c.opCMP,
		InstructionVariant{0xC9, Immediate, 2},
		InstructionVariant{0xC5, ZeroPage, 3},
		InstructionVariant{0xD5, ZeroPageX, 4},
		InstructionVariant{0xCD, Absolute, 4},
		InstructionVariant{0xDD, AbsoluteX, 4},
		InstructionVariant{0xD9, AbsoluteY, 4},
		InstructionVariant{0xC1, IndirectX, 6},
		InstructionVariant{0xD1, IndirectY, 5},
	)

	// CPX - Compare X Register
	c.instructions.registerVariants("CPX", c.opCPX,
		InstructionVariant{0xE0, Immediate, 2},
		InstructionVariant{0xE4, ZeroPage, 3},
		InstructionVariant{0xEC, Absolute, 4},
	)

	// CPY - Compare Y Register
	c.instructions.registerVariants("CPY", c.opCPY,
		InstructionVariant{0xC0, Immediate, 2},
		InstructionVariant{0xC4, ZeroPage, 3},
		InstructionVariant{0xCC, Absolute, 4},
	)

	// TAX - Transfer Accumulator to X
	c.instructions.registerVariants("TAX", c.opTAX,
		InstructionVariant{0xAA, Implicit, 2},
	)

	// TXA - Transfer X to Accumulator
	c.instructions.registerVariants("TXA", c.opTXA,
		InstructionVariant{0x8A, Implicit, 2},
	)

	// TAY - Transfer Accumulator to Y
	c.instructions.registerVariants("TAY", c.opTAY,
		InstructionVariant{0xA8, Implicit, 2},
	)

	// TYA - Transfer Y to Accumulator
	c.instructions.registerVariants("TYA", c.opTYA,
		InstructionVariant{0x98, Implicit, 2},
	)

	// TSX - Transfer Stack Pointer to X
	c.instructions.registerVariants("TSX", c.opTSX,
		InstructionVariant{0xBA, Implicit, 2},
	)

	// TXS - Transfer X to Stack Pointer
	c.instructions.registerVariants("TXS", c.opTXS,
		InstructionVariant{0x9A, Implicit, 2},
	)

	// BCS - Branch if Carry Set
	c.instructions.registerVariants("BCS", c.opBCS,
		InstructionVariant{0xB0, Relative, 2},
	)

	// BCC - Branch if Carry Clear
	c.instructions.registerVariants("BCC", c.opBCC,
		InstructionVariant{0x90, Relative, 2},
	)

	// BEQ - Branch if Equal
	c.instructions.registerVariants("BEQ", c.opBEQ,
		InstructionVariant{0xF0, Relative, 2},
	)

	// BNE - Branch if Not Equal
	c.instructions.registerVariants("BNE", c.opBNE,
		InstructionVariant{0xD0, Relative, 2},
	)

	// BMI - Branch if Minus
	c.instructions.registerVariants("BMI", c.opBMI,
		InstructionVariant{0x30, Relative, 2},
	)

	// BPL - Branch if Positive
	c.instructions.registerVariants("BPL", c.opBPL,
		InstructionVariant{0x10, Relative, 2},
	)

	// BVS - Branch if Overflow Set
	c.instructions.registerVariants("BVS", c.opBVS,
		InstructionVariant{0x70, Relative, 2},
	)

	// BVC - Branch if Overflow Clear
	c.instructions.registerVariants("BVC", c.opBVC,
		InstructionVariant{0x50, Relative, 2},
	)

	// LDA - Load Accumulator
	c.instructions.registerVariants("LDA", c.opLDA,
		InstructionVariant{0xA9, Immediate, 2},
		InstructionVariant{0xA5, ZeroPage, 3},
		InstructionVariant{0xB5, ZeroPageX, 4},
		InstructionVariant{0xAD, Absolute, 4},
		InstructionVariant{0xBD, AbsoluteX, 4},
		InstructionVariant{0xB9, AbsoluteY, 4},
		InstructionVariant{0xA1, IndirectX, 6},
		InstructionVariant{0xB1, IndirectY, 5},
	)

	// LDX - Load X Register
	c.instructions.registerVariants("LDX", c.opLDX,
		InstructionVariant{0xA2, Immediate, 2},
		InstructionVariant{0xA6, ZeroPage, 3},
		InstructionVariant{0xB6, ZeroPageY, 4},
		InstructionVariant{0xAE, Absolute, 4},
		InstructionVariant{0xBE, AbsoluteY, 4},
	)

	// LDY - Load Y Register
	c.instructions.registerVariants("LDY", c.opLDY,
		InstructionVariant{0xA0, Immediate, 2},
		InstructionVariant{0xA4, ZeroPage, 3},
		InstructionVariant{0xB4, ZeroPageX, 4},
		InstructionVariant{0xAC, Absolute, 4},
		InstructionVariant{0xBC, AbsoluteX, 4},
	)

	// STA - Store Accumulator
	c.instructions.registerVariants("STA", c.opSTA,
		InstructionVariant{0x85, ZeroPage, 3},
		InstructionVariant{0x95, ZeroPageX, 4},
		InstructionVariant{0x8D, Absolute, 4},
		InstructionVariant{0x9D, AbsoluteX, 5},
		InstructionVariant{0x99, AbsoluteY, 5},
		InstructionVariant{0x81, IndirectX, 6},
		InstructionVariant{0x91, IndirectY, 6},
	)

	// STX - Store X Register
	c.instructions.registerVariants("STX", c.opSTX,
		InstructionVariant{0x86, ZeroPage, 3},
		InstructionVariant{0x96, ZeroPageY, 4},
		InstructionVariant{0x8E, Absolute, 4},
	)

	// STY - Store Y Register
	c.instructions.registerVariants("STY", c.opSTY,
		InstructionVariant{0x84, ZeroPage, 3},
		InstructionVariant{0x94, ZeroPageX, 4},
		InstructionVariant{0x8C, Absolute, 4},
	)

	// CLC - Clear Carry Flag
	c.instructions.registerVariants("CLC", c.opCLC,
		InstructionVariant{0x18, Implicit, 2},
	)

	// SEC - Set Carry Flag
	c.instructions.registerVariants("SEC", c.opSEC,
		InstructionVariant{0x38, Implicit, 2},
	)

	// CLD - Set Decimal Mode
	c.instructions.registerVariants("CLD", c.opCLD,
		InstructionVariant{0xD8, Implicit, 2},
	)

	// SED - Set Decimal Mode
	c.instructions.registerVariants("SED", c.opSED,
		InstructionVariant{0xF8, Implicit, 2},
	)

	// CLI - Clear Interrupt Disable
	c.instructions.registerVariants("CLI", c.opCLI,
		InstructionVariant{0x58, Implicit, 2},
	)

	// SEI - Set Interrupt Disable
	c.instructions.registerVariants("SEI", c.opSEI,
		InstructionVariant{0x78, Implicit, 2},
	)

	// CLV - Clear Overflow Flag
	c.instructions.registerVariants("CLV", c.opCLV,
		InstructionVariant{0xB8, Implicit, 2},
	)

	// JMP - Jump
	c.instructions.registerVariants("JMP", c.opJMP,
		InstructionVariant{0x4C, Absolute, 3},
		InstructionVariant{0x6C, Indirect, 5},
	)

	// JSR - Jump to Subroutine
	c.instructions.registerVariants("JSR", c.opJSR,
		InstructionVariant{0x20, Absolute, 6},
	)

	// RTS - Return from Subroutine
	c.instructions.registerVariants("RTS", c.opRTS,
		InstructionVariant{0x60, Implicit, 6},
	)

	// PHA - Push Accumulator
	c.instructions.registerVariants("PHA", c.opPHA,
		InstructionVariant{0x48, Implicit, 3},
	)

	// PLA - Pull Accumulator
	c.instructions.registerVariants("PLA", c.opPLA,
		InstructionVariant{0x68, Implicit, 4},
	)

	// PHP - Push Processor Status
	c.instructions.registerVariants("PHP", c.opPHP,
		InstructionVariant{0x08, Implicit, 3},
	)

	// PLP - Pull Processor Status
	c.instructions.registerVariants("PLP", c.opPLP,
		InstructionVariant{0x28, Implicit, 4},
	)

	// NOP - No Operation
	c.instructions.registerVariants("NOP", c.opNOP,
		InstructionVariant{0xEA, Implicit, 2},
	)

	// RTI - Return from Interrupt
	c.instructions.registerVariants("RTI", c.opRTI,
		InstructionVariant{0x40, Implicit, 6},
	)

	// BRK - Force Interrupt
	c.instructions.registerVariants("BRK", c.opBRK,
		InstructionVariant{0x00, Implicit, 7},
	)

	// BIT - Bit Test
	c.instructions.registerVariants("BIT", c.opBIT,
		InstructionVariant{0x24, ZeroPage, 3},
		InstructionVariant{0x2C, Absolute, 4},
	)

	// ROL - Rotate Left
	c.instructions.registerVariants("ROL", c.opROL,
		InstructionVariant{0x2A, Accumulator, 2},
		InstructionVariant{0x26, ZeroPage, 5},
		InstructionVariant{0x36, ZeroPageX, 6},
		InstructionVariant{0x2E, Absolute, 6},
		InstructionVariant{0x3E, AbsoluteX, 7},
	)

	// ROR - Rotate Right
	c.instructions.registerVariants("ROR", c.opROR,
		InstructionVariant{0x6A, Accumulator, 2},
		InstructionVariant{0x66, ZeroPage, 5},
		InstructionVariant{0x76, ZeroPageX, 6},
		InstructionVariant{0x6E, Absolute, 6},
		InstructionVariant{0x7E, AbsoluteX, 7},
	)

	// ASL - Arithmetic Shift Left
	c.instructions.registerVariants("ASL", c.opASL,
		InstructionVariant{0x0A, Accumulator, 2},
		InstructionVariant{0x06, ZeroPage, 5},
		InstructionVariant{0x16, ZeroPageX, 6},
		InstructionVariant{0x0E, Absolute, 6},
		InstructionVariant{0x1E, AbsoluteX, 7},
	)

	// LSR - Logical Shift Right
	c.instructions.registerVariants("LSR", c.opLSR,
		InstructionVariant{0x4A, Accumulator, 2},
		InstructionVariant{0x46, ZeroPage, 5},
		InstructionVariant{0x56, ZeroPageX, 6},
		InstructionVariant{0x4E, Absolute, 6},
		InstructionVariant{0x5E, AbsoluteX, 7},
	)

	// ADC - Add with Carry
	c.instructions.registerVariants("ADC", c.opADC,
		InstructionVariant{0x69, Immediate, 2},
		InstructionVariant{0x65, ZeroPage, 3},
		InstructionVariant{0x75, ZeroPageX, 4},
		InstructionVariant{0x6D, Absolute, 4},
		InstructionVariant{0x7D, AbsoluteX, 4},
		InstructionVariant{0x79, AbsoluteY, 4},
		InstructionVariant{0x61, IndirectX, 6},
		InstructionVariant{0x71, IndirectY, 5},
	)

	// SBC - Subtract with Carry
	c.instructions.registerVariants("SBC", c.opSBC,
		InstructionVariant{0xE9, Immediate, 2},
		InstructionVariant{0xE5, ZeroPage, 3},
		InstructionVariant{0xF5, ZeroPageX, 4},
		InstructionVariant{0xED, Absolute, 4},
		InstructionVariant{0xFD, AbsoluteX, 4},
		InstructionVariant{0xF9, AbsoluteY, 4},
		InstructionVariant{0xE1, IndirectX, 6},
		InstructionVariant{0xF1, IndirectY, 5},
	)
}

func (t *InstructionTable) registerVariant(mnemonic string, handler OpcodeHandler, variant InstructionVariant) {
	if _, ok := (*t)[variant.Opcode]; ok {
		panic(fmt.Errorf("duplicate opcode registration: 0x%02X", variant.Opcode))
	}

	(*t)[variant.Opcode] = Instruction{
		Mnemonic: mnemonic,
		Handler:  handler,
		Variant:  variant,
	}
}

func (t *InstructionTable) registerVariants(mnemonic string, handler OpcodeHandler, variants ...InstructionVariant) {
	for _, variant := range variants {
		t.registerVariant(mnemonic, handler, variant)
	}
}
