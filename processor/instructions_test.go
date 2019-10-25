package processor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type memoryDiff struct {
	Address  uint16
	Expected uint8
	Actual   uint8
}

type cpuTestState struct {
	Cycles    Cycles
	Registers Registers
	Memory    []uint8
}

type cpuTestFunc func(cpu *CPU, state *cpuTestState)

const MemoryTestLocation uint16 = 0xC000
const AbsoluteTestLocation uint16 = 0xD000

func diffMemory(expected []uint8, actual []uint8) (diffs []memoryDiff) {
	for address := 0; address < len(expected); address++ {
		if expected[address] != actual[address] {
			diffs = append(diffs, memoryDiff{
				Address:  uint16(address),
				Expected: expected[address],
				Actual:   actual[address],
			})
		}
	}

	return
}

func testCPU(t *testing.T, testFunc cpuTestFunc) {
	cpu := NewCPU()
	expectedState := &cpuTestState{
		Cycles:    cpu.TotalCycles,
		Registers: cpu.Registers,
		Memory:    cpu.Memory.Dump(),
	}

	testFunc(cpu, expectedState)
	cpu.Execute()

	actualState := &cpuTestState{
		Cycles:    cpu.TotalCycles,
		Registers: cpu.Registers,
		Memory:    cpu.Memory.Dump(),
	}
	memoryDiff := diffMemory(expectedState.Memory, actualState.Memory)

	assert.Equal(t, expectedState.Cycles, actualState.Cycles, "unexpected cycle count")
	assert.EqualValues(t, expectedState.Registers, cpu.Registers, "unexpected cpu registers")
	assert.Emptyf(t, memoryDiff, "unexpected memory changes: %+v", memoryDiff)
}

func (s *cpuTestState) expectFlag(flag Status, isEnabled bool) {
	if isEnabled {
		s.Registers.P |= flag
	} else {
		s.Registers.P &^= flag
	}
}

func (s *cpuTestState) expectStack(address uint8, value uint8) {
	s.Memory[0x0100|uint16(address)] = value
}

func (s *cpuTestState) expectStack16(lowAddress uint8, value uint16) {
	s.expectStack(lowAddress, uint8(value&0xFF))
	s.expectStack(lowAddress+1, uint8((value>>8)&0xFF))
}

func testCommon(cpu *CPU, state *cpuTestState, opcode uint8) {
	// prepare
	cpu.Registers.PC = MemoryTestLocation
	cpu.Memory.Poke(cpu.Registers.PC, opcode)

	// verify
	state.Registers.PC = MemoryTestLocation + 1
	state.Memory[cpu.Registers.PC] = opcode
	if instruction, ok := cpu.instructions[Opcode(opcode)]; ok {
		state.Cycles = instruction.Variant.StaticCycles
	}
}

func testImplicit(cpu *CPU, state *cpuTestState, opcode uint8) {
	testCommon(cpu, state, opcode)
}

func testImmediate(cpu *CPU, state *cpuTestState, opcode uint8, value uint8) {
	testCommon(cpu, state, opcode)

	// prepare
	cpu.Memory.Poke(cpu.Registers.PC+1, value)

	// verify
	state.Registers.PC++
	state.Memory[cpu.Registers.PC+1] = value
}

func testAbsolute(cpu *CPU, state *cpuTestState, opcode uint8, value uint8) {
	testCommon(cpu, state, opcode)

	// prepare
	cpu.Memory.Poke16(cpu.Registers.PC+1, AbsoluteTestLocation)
	cpu.Memory.Poke(AbsoluteTestLocation, value)

	// verify
	state.Registers.PC += 2
	state.Memory[cpu.Registers.PC+1] = uint8(AbsoluteTestLocation & 0xFF)
	state.Memory[cpu.Registers.PC+2] = uint8((AbsoluteTestLocation >> 8) & 0xFF)
	state.Memory[AbsoluteTestLocation] = value
}

func testAbsoluteDirect(cpu *CPU, state *cpuTestState, opcode uint8, address uint16) {
	testCommon(cpu, state, opcode)

	// prepare
	cpu.Memory.Poke16(cpu.Registers.PC+1, address)

	// verify
	state.Registers.PC += 2
	state.Memory[cpu.Registers.PC+1] = uint8(address & 0xFF)
	state.Memory[cpu.Registers.PC+2] = uint8((address >> 8) & 0xFF)
}

func testRelative(cpu *CPU, state *cpuTestState, opcode uint8, offset int8) {
	testCommon(cpu, state, opcode)

	// convert
	var value uint8
	if value < 0 {
		value = 0xFF - uint8(-offset)
	} else {
		value = uint8(offset)
	}

	// prepare
	cpu.Memory.Poke(cpu.Registers.PC+1, value)

	// verify
	state.Registers.PC++
	state.Memory[cpu.Registers.PC+1] = value
}

func TestAND(t *testing.T) {
	testAND := func(a uint8, b uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0x29, b)
			cpu.Registers.A = a

			// verify
			state.Registers.A = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testAND(0x10, 0x10, 0x10, false, false)
	testAND(0x10, 0x01, 0x00, true, false)
	testAND(0xFF, 0x80, 0x80, false, true)
}

func TestORA(t *testing.T) {
	testORA := func(a uint8, b uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0x09, b)
			cpu.Registers.A = a

			// verify
			state.Registers.A = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testORA(0x10, 0x01, 0x11, false, false)
	testORA(0x00, 0x00, 0x00, true, false)
	testORA(0x80, 0x01, 0x81, false, true)
}

func TestEOR(t *testing.T) {
	testEOR := func(a uint8, b uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0x49, b)
			cpu.Registers.A = a

			// verify
			state.Registers.A = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testEOR(0x10, 0x01, 0x11, false, false)
	testEOR(0x10, 0x10, 0x00, true, false)
	testEOR(0x81, 0x01, 0x80, false, true)
}

func TestINC(t *testing.T) {
	testINC := func(value uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0xEE, value)

			// verify
			state.Memory[AbsoluteTestLocation] = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testINC(0x0F, 0x10, false, false)
	testINC(0xFF, 0x00, true, false)
	testINC(0x7F, 0x80, false, true)
}

func TestINX(t *testing.T) {
	testINX := func(value uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0xE8)
			cpu.Registers.X = value

			// verify
			state.Registers.X = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testINX(0x0F, 0x10, false, false)
	testINX(0xFF, 0x00, true, false)
	testINX(0x7F, 0x80, false, true)
}

func TestINY(t *testing.T) {
	testINY := func(value uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0xC8)
			cpu.Registers.Y = value

			// verify
			state.Registers.Y = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testINY(0x0F, 0x10, false, false)
	testINY(0xFF, 0x00, true, false)
	testINY(0x7F, 0x80, false, true)
}

func TestDEC(t *testing.T) {
	testDEC := func(value uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0xCE, value)

			// verify
			state.Memory[AbsoluteTestLocation] = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testDEC(0x01, 0x00, true, false)
	testDEC(0x00, 0xFF, false, true)
	testDEC(0x80, 0x7F, false, false)
}

func TestDEX(t *testing.T) {
	testDEX := func(value uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0xCA)
			cpu.Registers.X = value

			// verify
			state.Registers.X = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testDEX(0x01, 0x00, true, false)
	testDEX(0x00, 0xFF, false, true)
	testDEX(0x80, 0x7F, false, false)
}

func TestDEY(t *testing.T) {
	testDEY := func(value uint8, result uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x88)
			cpu.Registers.Y = value

			// verify
			state.Registers.Y = result
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testDEY(0x01, 0x00, true, false)
	testDEY(0x00, 0xFF, false, true)
	testDEY(0x80, 0x7F, false, false)
}

func TestCMP(t *testing.T) {
	testCMP := func(a uint8, b uint8, isCarry bool, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0xC9, b)
			cpu.Registers.A = a

			// verify
			state.Registers.A = a
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testCMP(0x01, 0xFF, false, false, false)
	testCMP(0x7F, 0x80, false, false, true)
	testCMP(0x40, 0x20, true, false, false)
	testCMP(0x42, 0x42, true, true, false)
}

func TestCPX(t *testing.T) {
	testCPX := func(a uint8, b uint8, isCarry bool, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0xE0, b)
			cpu.Registers.X = a

			// verify
			state.Registers.X = a
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testCPX(0x01, 0xFF, false, false, false)
	testCPX(0x7F, 0x80, false, false, true)
	testCPX(0x40, 0x20, true, false, false)
	testCPX(0x42, 0x42, true, true, false)
}

func TestCPY(t *testing.T) {
	testCPY := func(a uint8, b uint8, isCarry bool, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0xC0, b)
			cpu.Registers.Y = a

			// verify
			state.Registers.Y = a
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testCPY(0x01, 0xFF, false, false, false)
	testCPY(0x7F, 0x80, false, false, true)
	testCPY(0x40, 0x20, true, false, false)
	testCPY(0x42, 0x42, true, true, false)
}

func TestTAX(t *testing.T) {
	testTAX := func(value uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0xAA)
			cpu.Registers.A = value

			// verify
			state.Registers.A = value
			state.Registers.X = value
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testTAX(0x42, false, false)
	testTAX(0x00, true, false)
	testTAX(0x80, false, true)
}

func TestTXA(t *testing.T) {
	testTXA := func(value uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x8A)
			cpu.Registers.X = value

			// verify
			state.Registers.X = value
			state.Registers.A = value
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testTXA(0x42, false, false)
	testTXA(0x00, true, false)
	testTXA(0x80, false, true)
}

func TestTAY(t *testing.T) {
	testTAY := func(value uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0xA8)
			cpu.Registers.A = value

			// verify
			state.Registers.A = value
			state.Registers.Y = value
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testTAY(0x42, false, false)
	testTAY(0x00, true, false)
	testTAY(0x80, false, true)
}

func TestTYA(t *testing.T) {
	testTYA := func(value uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x98)
			cpu.Registers.Y = value

			// verify
			state.Registers.Y = value
			state.Registers.A = value
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testTYA(0x42, false, false)
	testTYA(0x00, true, false)
	testTYA(0x80, false, true)
}

func TestBCS(t *testing.T) {
	testBCS := func(offset int8, hasCarry bool, isSuccessful bool, extraCycles Cycles) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testRelative(cpu, state, 0xB0, offset)
			if hasCarry {
				cpu.Registers.P |= FlagCarry
			} else {
				cpu.Registers.P &^= FlagCarry
			}

			// verify
			state.Cycles += extraCycles
			state.expectFlag(FlagCarry, hasCarry)
			if isSuccessful {
				state.Registers.PC = uint16(int16(state.Registers.PC) + int16(offset))
			}
		})
	}

	testBCS(+10, false, false, 0)
	testBCS(+10, true, true, 1)
	testBCS(-10, true, true, 2)
}

func TestBCC(t *testing.T) {
	testBCC := func(offset int8, hasCarry bool, isSuccessful bool, extraCycles Cycles) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testRelative(cpu, state, 0x90, offset)
			if hasCarry {
				cpu.Registers.P |= FlagCarry
			} else {
				cpu.Registers.P &^= FlagCarry
			}

			// verify
			state.Cycles += extraCycles
			state.expectFlag(FlagCarry, hasCarry)
			if isSuccessful {
				state.Registers.PC = uint16(int16(state.Registers.PC) + int16(offset))
			}
		})
	}

	testBCC(+10, true, false, 0)
	testBCC(+10, false, true, 1)
	testBCC(-10, false, true, 2)
}

func TestBEQ(t *testing.T) {
	testBEQ := func(offset int8, hasZero bool, isSuccessful bool, extraCycles Cycles) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testRelative(cpu, state, 0xF0, offset)
			if hasZero {
				cpu.Registers.P |= FlagZero
			} else {
				cpu.Registers.P &^= FlagZero
			}

			// verify
			state.Cycles += extraCycles
			state.expectFlag(FlagZero, hasZero)
			if isSuccessful {
				state.Registers.PC = uint16(int16(state.Registers.PC) + int16(offset))
			}
		})
	}

	testBEQ(+10, false, false, 0)
	testBEQ(+10, true, true, 1)
	testBEQ(-10, true, true, 2)
}

func TestBNE(t *testing.T) {
	testBNE := func(offset int8, hasZero bool, isSuccessful bool, extraCycles Cycles) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testRelative(cpu, state, 0xD0, offset)
			if hasZero {
				cpu.Registers.P |= FlagZero
			} else {
				cpu.Registers.P &^= FlagZero
			}

			// verify
			state.Cycles += extraCycles
			state.expectFlag(FlagZero, hasZero)
			if isSuccessful {
				state.Registers.PC = uint16(int16(state.Registers.PC) + int16(offset))
			}
		})
	}

	testBNE(+10, true, false, 0)
	testBNE(+10, false, true, 1)
	testBNE(-10, false, true, 2)
}

func TestBMI(t *testing.T) {
	testBMI := func(offset int8, hasNegative bool, isSuccessful bool, extraCycles Cycles) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testRelative(cpu, state, 0x30, offset)
			if hasNegative {
				cpu.Registers.P |= FlagNegative
			} else {
				cpu.Registers.P &^= FlagNegative
			}

			// verify
			state.Cycles += extraCycles
			state.expectFlag(FlagNegative, hasNegative)
			if isSuccessful {
				state.Registers.PC = uint16(int16(state.Registers.PC) + int16(offset))
			}
		})
	}

	testBMI(+10, false, false, 0)
	testBMI(+10, true, true, 1)
	testBMI(-10, true, true, 2)
}

func TestBPL(t *testing.T) {
	testBPL := func(offset int8, hasNegative bool, isSuccessful bool, extraCycles Cycles) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testRelative(cpu, state, 0x10, offset)
			if hasNegative {
				cpu.Registers.P |= FlagNegative
			} else {
				cpu.Registers.P &^= FlagNegative
			}

			// verify
			state.Cycles += extraCycles
			state.expectFlag(FlagNegative, hasNegative)
			if isSuccessful {
				state.Registers.PC = uint16(int16(state.Registers.PC) + int16(offset))
			}
		})
	}

	testBPL(+10, true, false, 0)
	testBPL(+10, false, true, 1)
	testBPL(-10, false, true, 2)
}

func TestBVS(t *testing.T) {
	testBVS := func(offset int8, hasOverflow bool, isSuccessful bool, extraCycles Cycles) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testRelative(cpu, state, 0x70, offset)
			if hasOverflow {
				cpu.Registers.P |= FlagOverflow
			} else {
				cpu.Registers.P &^= FlagOverflow
			}

			// verify
			state.Cycles += extraCycles
			state.expectFlag(FlagOverflow, hasOverflow)
			if isSuccessful {
				state.Registers.PC = uint16(int16(state.Registers.PC) + int16(offset))
			}
		})
	}

	testBVS(+10, false, false, 0)
	testBVS(+10, true, true, 1)
	testBVS(-10, true, true, 2)
}

func TestBVC(t *testing.T) {
	testBVC := func(offset int8, hasOverflow bool, isSuccessful bool, extraCycles Cycles) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testRelative(cpu, state, 0x50, offset)
			if hasOverflow {
				cpu.Registers.P |= FlagOverflow
			} else {
				cpu.Registers.P &^= FlagOverflow
			}

			// verify
			state.Cycles += extraCycles
			state.expectFlag(FlagOverflow, hasOverflow)
			if isSuccessful {
				state.Registers.PC = uint16(int16(state.Registers.PC) + int16(offset))
			}
		})
	}

	testBVC(+10, true, false, 0)
	testBVC(+10, false, true, 1)
	testBVC(-10, false, true, 2)
}

func TestSTA(t *testing.T) {
	testSTA := func(value uint8) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0x8D, value)
			cpu.Registers.A = value

			// verify
			state.Registers.A = value
			state.Memory[AbsoluteTestLocation] = value
		})
	}

	testSTA(0x13)
	testSTA(0x37)
}

func TestSTX(t *testing.T) {
	testSTX := func(value uint8) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0x8E, value)
			cpu.Registers.X = value

			// verify
			state.Registers.X = value
			state.Memory[AbsoluteTestLocation] = value
		})
	}

	testSTX(0x13)
	testSTX(0x37)
}

func TestSTY(t *testing.T) {
	testSTY := func(value uint8) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0x8C, value)
			cpu.Registers.Y = value

			// verify
			state.Registers.Y = value
			state.Memory[AbsoluteTestLocation] = value
		})
	}

	testSTY(0x13)
	testSTY(0x37)
}

func TestCLC(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0x18)

		// verify
		state.expectFlag(FlagCarry, false)
	})
}

func TestSEC(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0x38)

		// verify
		state.expectFlag(FlagCarry, true)
	})
}

func TestCLD(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0xD8)

		// verify
		state.expectFlag(FlagDecimal, false)
	})
}

func TestSED(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0xF8)

		// verify
		state.expectFlag(FlagDecimal, true)
	})
}

func TestCLI(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0x58)

		// verify
		state.expectFlag(FlagInterruptDisable, false)
	})
}

func TestSEI(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0x78)

		// verify
		state.expectFlag(FlagInterruptDisable, true)
	})
}

func TestCLV(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0xB8)

		// verify
		state.expectFlag(FlagOverflow, false)
	})
}

func TestLDA(t *testing.T) {
	testLDA := func(value uint8) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0xA9, value)

			// verify
			state.Registers.A = value
		})
	}

	testLDA(0x13)
	testLDA(0x37)
}

func TestLDX(t *testing.T) {
	testLDX := func(value uint8) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0xA2, value)

			// verify
			state.Registers.X = value
		})
	}

	testLDX(0x13)
	testLDX(0x37)
}

func TestLDY(t *testing.T) {
	testLDY := func(value uint8) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0xA0, value)

			// verify
			state.Registers.Y = value
		})
	}

	testLDY(0x13)
	testLDY(0x37)
}

func TestJMP(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testAbsoluteDirect(cpu, state, 0x4C, 0xBEEF)

		// verify
		state.Registers.PC = 0xBEEF
	})
}

func TestJSR(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testAbsoluteDirect(cpu, state, 0x20, 0xBEEF)

		// verify
		previousPC := state.Registers.PC - 1
		state.Registers.PC = 0xBEEF
		state.Registers.S -= 2
		state.expectStack16(0xFC, previousPC)
	})
}

func TestRTS(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0x60)
		cpu.Push16(0x1233)

		// verify
		state.Registers.PC = 0x1234
		state.expectStack16(0xFC, 0x1233)
	})
}

func TestPHA(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0x48)
		cpu.Registers.A = 0x42

		// verify
		state.Registers.A = 0x42
		state.Registers.S--
		state.expectStack(0xFD, 0x42)
	})
}

func TestPLA(t *testing.T) {
	testPLA := func(value uint8, isZero bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x68)
			cpu.Push(value)

			// verify
			state.Registers.A = value
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
			state.expectStack(0xFD, value)
		})
	}

	testPLA(0x42, false, false)
	testPLA(0x00, true, false)
	testPLA(0x80, false, true)
}

func TestPHP(t *testing.T) {
	testPHP := func(actual Status, expected Status) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x08)
			cpu.Registers.P = actual

			// verify
			state.Registers.P = actual
			state.Registers.S--
			state.expectStack(0xFD, uint8(expected))
		})
	}

	testPHP(
		FlagCarry|FlagZero|FlagInterruptDisable|FlagDecimal|FlagBreak|FlagUnused|FlagOverflow|FlagNegative,
		FlagCarry|FlagZero|FlagInterruptDisable|FlagDecimal|FlagBreak|FlagUnused|FlagOverflow|FlagNegative,
	)
	testPHP(0, FlagBreak)
}

func TestPLP(t *testing.T) {
	testPLP := func(actual Status, expected Status) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x28)
			cpu.Push(uint8(actual))

			// verify
			state.Registers.P = expected
			state.expectStack(0xFD, uint8(actual))
		})
	}

	testPLP(0xFF, 0xFF&^FlagBreak)
	testPLP(0, FlagUnused)
}

func TestNOP(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0xEA)
	})
}

func TestRTI(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0x40)
		cpu.Push16(0x1234)
		cpu.Push(0xFF)

		// verify
		state.Registers.PC = 0x1234
		state.Registers.P = 0xFF &^ FlagBreak
		state.expectStack(0xFB, 0xFF)
		state.expectStack16(0xFC, 0x1234)
	})
}

func TestBRK(t *testing.T) {
	testCPU(t, func(cpu *CPU, state *cpuTestState) {
		// prepare
		testImplicit(cpu, state, 0x00)
		cpu.Memory.Poke16(ResetVector, 0x1234)

		// verify
		previousPC := state.Registers.PC
		state.Registers.PC = 0x1234
		state.Registers.P |= FlagInterruptDisable
		state.Registers.S -= 3
		state.Memory[ResetVector] = 0x34
		state.Memory[ResetVector+1] = 0x12
		state.expectStack16(0xFC, previousPC)
		state.expectStack(0xFB, uint8(state.Registers.P|FlagBreak))
	})
}

func TestBIT(t *testing.T) {
	testBIT := func(a uint8, b uint8, isZero bool, isOverflow bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0x2C, b)
			cpu.Registers.A = a

			// verify
			state.Registers.A = a
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagOverflow, isOverflow)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testBIT(0x80, 0x01, true, false, false)
	testBIT(0x03, 0x01, false, false, false)
	testBIT(0x01, 0x81, false, false, true)
	testBIT(0xFF, 0x7F, false, true, false)
}

func TestROL(t *testing.T) {
	testROL := func(before uint8, after uint8, isCarry bool, isZero bool, isNegative bool) {
		// Implicit
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x2A)
			cpu.Registers.A = before

			// verify
			state.Registers.A = after
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})

		// Absolute
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0x2E, before)

			// verify
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
			state.Memory[AbsoluteTestLocation] = after
		})
	}

	testROL(0x1, 0x2, false, false, false)
	testROL(0x80, 0x0, true, true, false)
}

func TestROR(t *testing.T) {
	testROR := func(before uint8, after uint8, isCarry bool, isZero bool, isNegative bool) {
		// Implicit
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x6A)
			cpu.Registers.A = before

			// verify
			state.Registers.A = after
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})

		// Absolute
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0x6E, before)

			// verify
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
			state.Memory[AbsoluteTestLocation] = after
		})
	}

	testROR(0x2, 0x1, false, false, false)
	testROR(0x1, 0x0, true, true, false)
}

func TestASL(t *testing.T) {
	testASL := func(before uint8, after uint8, isCarry bool, isZero bool, isNegative bool) {
		// Implicit
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x0A)
			cpu.Registers.A = before

			// verify
			state.Registers.A = after
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})

		// Absolute
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0x0E, before)

			// verify
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
			state.Memory[AbsoluteTestLocation] = after
		})
	}

	testASL(0x1, 0x2, false, false, false)
	testASL(0x80, 0x0, true, true, false)
	testASL(0x7F, 0xFE, false, false, true)
}

func TestLSR(t *testing.T) {
	testLSR := func(before uint8, after uint8, isCarry bool, isZero bool, isNegative bool) {
		// Implicit
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImplicit(cpu, state, 0x4A)
			cpu.Registers.A = before

			// verify
			state.Registers.A = after
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
		})

		// Absolute
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testAbsolute(cpu, state, 0x4E, before)

			// verify
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagNegative, isNegative)
			state.Memory[AbsoluteTestLocation] = after
		})
	}

	testLSR(0x80, 0x40, false, false, false)
	testLSR(0x1, 0x0, true, true, false)
	testLSR(0x3, 0x1, true, false, false)
}

func TestADC(t *testing.T) {
	testADC := func(a uint8, b uint8, result uint8, isCarry bool, isZero bool, isOverflow bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0x69, b)
			cpu.Registers.A = a

			// verify
			state.Registers.A = result
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagOverflow, isOverflow)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testADC(0x10, 0x20, 0x30, false, false, false, false)
	testADC(0x00, 0x00, 0x00, false, true, false, false)
	testADC(0x40, 0x40, 0x80, false, false, true, true)
	testADC(0xFF, 0xFF, 0xFE, true, false, false, true)
}

func TestSBC(t *testing.T) {
	testSBC := func(a uint8, b uint8, result uint8, isCarry bool, isZero bool, isOverflow bool, isNegative bool) {
		testCPU(t, func(cpu *CPU, state *cpuTestState) {
			// prepare
			testImmediate(cpu, state, 0xE9, b)
			cpu.Registers.A = a

			// verify
			state.Registers.A = result
			state.expectFlag(FlagCarry, isCarry)
			state.expectFlag(FlagZero, isZero)
			state.expectFlag(FlagOverflow, isOverflow)
			state.expectFlag(FlagNegative, isNegative)
		})
	}

	testSBC(0x30, 0x20, 0x0F, true, false, false, false)
	testSBC(0x02, 0x01, 0x00, true, true, false, false)
	testSBC(0x10, 0x20, 0xEF, false, false, false, true)
	testSBC(0x80, 0x01, 0x7E, true, false, true, false)
}
