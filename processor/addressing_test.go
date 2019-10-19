package processor

import "testing"

func testAddressingMode(t *testing.T, mode AddressingMode, address uint16, extraCycles Cycles, testFunc func(cpu *CPU)) {
	cpu := NewCPU()
	cpu.Registers.PC = 0x0100

	testFunc(cpu)
	actualAddress, actualCycles := cpu.lookupAddress(mode)
	modeName := AddressingModeName(mode)

	if actualAddress != address {
		t.Errorf("invalid address for mode [%s], expected [0x%04X] but got [0x%04X]",
			modeName, address, actualAddress)
	}
	if actualCycles != extraCycles {
		t.Errorf("invalid cycles for mode [%s], expected [%d] but got [%d]",
			modeName, extraCycles, actualCycles)
	}
}

func TestImplicit(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected implicit addressing mode to panic")
		}
	}()

	testAddressingMode(t, Implicit, 0x0, 0, func(cpu *CPU) {})
}

func TestAccumulator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected accumulator addressing mode to panic")
		}
	}()

	testAddressingMode(t, Accumulator, 0x0, 0, func(cpu *CPU) {})
}

func TestImmediate(t *testing.T) {
	testAddressingMode(t, Immediate, 0x0100, 0, func(cpu *CPU) {})
}

func TestZeroPage(t *testing.T) {
	testAddressingMode(t, ZeroPage, 0x42, 0, func(cpu *CPU) {
		cpu.Memory.Poke(0x0100, 0x42)
	})
}

func TestZeroPageX(t *testing.T) {
	testAddressingMode(t, ZeroPageX, 0x42, 0, func(cpu *CPU) {
		cpu.Registers.X = 0x10
		cpu.Memory.Poke(0x0100, 0x32)
	})
}

func TestZeroPageY(t *testing.T) {
	testAddressingMode(t, ZeroPageY, 0x42, 0, func(cpu *CPU) {
		cpu.Registers.Y = 0x10
		cpu.Memory.Poke(0x0100, 0x32)
	})
}

func TestRelativePositive(t *testing.T) {
	testAddressingMode(t, Relative, 0x123, 0, func(cpu *CPU) {
		cpu.Memory.Poke(0x0100, 0x22)
	})
}

func TestRelativeNegative(t *testing.T) {
	testAddressingMode(t, Relative, 0xDF, 0, func(cpu *CPU) {
		cpu.Memory.Poke(0x0100, 0x100-0x22)
	})
}

func TestAbsolute(t *testing.T) {
	testAddressingMode(t, Absolute, 0x1234, 0, func(cpu *CPU) {
		cpu.Memory.Poke16(0x0100, 0x1234)
	})
}

func TestAbsoluteX(t *testing.T) {
	testAddressingMode(t, AbsoluteX, 0x1234, 0, func(cpu *CPU) {
		cpu.Registers.X = 0x34
		cpu.Memory.Poke16(0x0100, 0x1200)
	})
}

func TestAbsoluteXWrap(t *testing.T) {
	testAddressingMode(t, AbsoluteX, 0x134F, 1, func(cpu *CPU) {
		cpu.Registers.X = 0xFF
		cpu.Memory.Poke16(0x0100, 0x1250)
	})
}

func TestAbsoluteY(t *testing.T) {
	testAddressingMode(t, AbsoluteY, 0x1234, 0, func(cpu *CPU) {
		cpu.Registers.Y = 0x34
		cpu.Memory.Poke16(0x0100, 0x1200)
	})
}

func TestAbsoluteYWrap(t *testing.T) {
	testAddressingMode(t, AbsoluteY, 0x134F, 1, func(cpu *CPU) {
		cpu.Registers.Y = 0xFF
		cpu.Memory.Poke16(0x0100, 0x1250)
	})
}

func TestIndirect(t *testing.T) {
	testAddressingMode(t, Indirect, 0x1234, 0, func(cpu *CPU) {
		cpu.Memory.Poke16(0x0100, 0xABCD)
		cpu.Memory.Poke16(0xABCD, 0x1234)
	})
}

func TestIndirectX(t *testing.T) {
	testAddressingMode(t, IndirectX, 0x1234, 0, func(cpu *CPU) {
		cpu.Registers.X = 0x32
		cpu.Memory.Poke(0x0100, 0x10)
		cpu.Memory.Poke16(0x42, 0x1234)
	})
}

func TestIndirectY(t *testing.T) {
	testAddressingMode(t, IndirectY, 0x42, 0, func(cpu *CPU) {
		cpu.Registers.Y = 0x32
		cpu.Memory.Poke(0x0100, 0xCC)
		cpu.Memory.Poke16(0xCC, 0x0010)
	})
}

func TestIndirectYWrap(t *testing.T) {
	testAddressingMode(t, IndirectY, 0x1234, 1, func(cpu *CPU) {
		cpu.Registers.Y = 0x34
		cpu.Memory.Poke(0x0100, 0x42)
		cpu.Memory.Poke16(0x42, 0x1200)
	})
}
