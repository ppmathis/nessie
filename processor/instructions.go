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
}

func (t *InstructionTable) register(mnemonic string, opcode Opcode, handler OpcodeHandler, mode AddressingMode, staticCycles Cycles) {
	t.registerVariant(mnemonic, handler, InstructionVariant{
		Opcode:         opcode,
		AddressingMode: mode,
		StaticCycles:   staticCycles,
	})
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
