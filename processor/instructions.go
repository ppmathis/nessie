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
