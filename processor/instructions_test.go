package processor

import (
	"fmt"
	"testing"
)

const TestAbsoluteAddress = 0x4242

var cpu *CPU

func setup() {
	cpu = NewCPU()
}

type testFlags struct {
	enabled  []Flag
	disabled []Flag
}

type testRegister struct {
	register Register
	expected uint16
}

func (f *testFlags) Add(flag Flag, isEnabled bool) {
	if isEnabled {
		f.enabled = append(f.enabled, flag)
	} else {
		f.disabled = append(f.disabled, flag)
	}
}

func assertCycles(t *testing.T, expected Cycles) {
	actual := cpu.TotalCycles

	if actual != expected {
		t.Errorf("expected [%d] cpu cycles, got [%d]", expected, actual)
	}
}

func assertFlags(t *testing.T, enabledFlags []Flag, disabledFlags []Flag) {
	for _, flag := range enabledFlags {
		if !cpu.GetFlag(flag) {
			t.Errorf("expected flag [%s] to be set", cpu.GetFlagName(flag))
		}
	}

	for _, flag := range disabledFlags {
		if cpu.GetFlag(flag) {
			t.Errorf("expected flag [%s] to be unset", cpu.GetFlagName(flag))
		}
	}
}

func assertRegister(t *testing.T, register Register, expected uint16) {
	actual := cpu.GetRegister(register)
	registerName := cpu.GetRegisterName(register)

	if actual != expected {
		t.Errorf("expected value [0x%04X] for register %s, got [0x%04X]", expected, registerName, actual)
	}
}

func assertMemory(t *testing.T, address uint16, expected uint8) {
	actual := cpu.Memory.Peek(address)

	if actual != expected {
		t.Errorf("expected address [0x%04X] to contain [0x%02X], got [0x%02X]", address, expected, actual)
	}
}

func assertCPU(t *testing.T, expectedCycles Cycles, flags testFlags, registers ...testRegister) {
	assertCycles(t, expectedCycles)
	assertFlags(t, flags.enabled, flags.disabled)

	for _, assert := range registers {
		assertRegister(t, assert.register, assert.expected)
	}
}

func testImplicit(opcode uint8, registers Registers) {
	setup()

	cpu.Registers = registers
	cpu.Registers.PC = 0x0100
	cpu.Memory.Poke(0x0100, opcode)
	cpu.Execute()
}

func testImplicitFlags(opcode uint8, flags Flags) {
	setup()

	cpu.Flags = flags
	cpu.Registers.PC = 0x0100
	cpu.Memory.Poke(0x0100, opcode)
	cpu.Execute()
}

func testRelative(opcode uint8, offset int8, flags Flags) {
	setup()

	var value uint8
	if value < 0 {
		value = 0xFF - uint8(-offset)
	} else {
		value = uint8(offset)
	}

	cpu.Flags = flags
	cpu.Registers.PC = 0x0100
	cpu.Memory.Poke(0x0100, opcode)
	cpu.Memory.Poke(0x0101, value)
	cpu.Execute()
}

func testImmediate(opcode uint8, argument uint8, registers Registers) {
	setup()

	cpu.Registers = registers
	cpu.Registers.PC = 0x0100
	cpu.Memory.Poke(0x0100, opcode)
	cpu.Memory.Poke(0x0101, argument)
	cpu.Execute()
}

func testAbsolute(opcode uint8, value uint8, registers Registers) {
	setup()

	cpu.Registers = registers
	cpu.Registers.PC = 0x0100
	cpu.Memory.Poke(0x0100, opcode)
	cpu.Memory.Poke16(0x0101, TestAbsoluteAddress)
	cpu.Memory.Poke(TestAbsoluteAddress, value)
	cpu.Execute()
}

func testAbsoluteAddress(opcode uint8, address uint16) {
	setup()

	cpu.Registers.PC = 0x0100
	cpu.Memory.Poke(0x0100, opcode)
	cpu.Memory.Poke16(0x0101, address)
	cpu.Execute()
}

func TestAND(t *testing.T) {
	testAND := func(a uint8, b uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testAND[%d, %d] =? %d\n", a, b, result)
		testImmediate(0x29, b, Registers{A: a})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterA, expected: uint16(result)})
	}

	testAND(0x10, 0x10, 0x10, false, false)
	testAND(0x10, 0x01, 0x00, true, false)
	testAND(0xFF, 0x80, 0x80, false, true)
}

func TestORA(t *testing.T) {
	testORA := func(a uint8, b uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testORA[%d, %d] =? %d\n", a, b, result)
		testImmediate(0x09, b, Registers{A: a})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterA, expected: uint16(result)})
	}

	testORA(0x10, 0x01, 0x11, false, false)
	testORA(0x00, 0x00, 0x00, true, false)
	testORA(0x80, 0x01, 0x81, false, true)
}

func TestEOR(t *testing.T) {
	testEOR := func(a uint8, b uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testEOR[%d, %d] =? %d\n", a, b, result)
		testImmediate(0x49, b, Registers{A: a})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterA, expected: uint16(result)})
	}

	testEOR(0x10, 0x01, 0x11, false, false)
	testEOR(0x10, 0x10, 0x00, true, false)
	testEOR(0x81, 0x01, 0x80, false, true)
}

func TestINC(t *testing.T) {
	testINC := func(value uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testINC[%d] =? %d\n", value, result)
		testAbsolute(0xEE, value, Registers{})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 6, flags)
		assertMemory(t, TestAbsoluteAddress, result)
	}

	testINC(0x0F, 0x10, false, false)
	testINC(0xFF, 0x00, true, false)
	testINC(0x7F, 0x80, false, true)
}

func TestINX(t *testing.T) {
	testINX := func(value uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testINX[%d] =? %d\n", value, result)
		testImplicit(0xE8, Registers{X: value})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterX, expected: uint16(result)})
	}

	testINX(0x0F, 0x10, false, false)
	testINX(0xFF, 0x00, true, false)
	testINX(0x7F, 0x80, false, true)
}

func TestINY(t *testing.T) {
	testINY := func(value uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testINY[%d] =? %d\n", value, result)
		testImplicit(0xC8, Registers{Y: value})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterY, expected: uint16(result)})
	}

	testINY(0x0F, 0x10, false, false)
	testINY(0xFF, 0x00, true, false)
	testINY(0x7F, 0x80, false, true)
}

func TestDEC(t *testing.T) {
	testDEC := func(value uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testDEC[%d] =? %d\n", value, result)
		testAbsolute(0xCE, value, Registers{})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 6, flags)
		assertMemory(t, TestAbsoluteAddress, result)
	}

	testDEC(0x01, 0x00, true, false)
	testDEC(0x00, 0xFF, false, true)
	testDEC(0x80, 0x7F, false, false)
}

func TestDEX(t *testing.T) {
	testDEX := func(value uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testINX[%d] =? %d\n", value, result)
		testImplicit(0xCA, Registers{X: value})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterX, expected: uint16(result)})
	}

	testDEX(0x01, 0x00, true, false)
	testDEX(0x00, 0xFF, false, true)
	testDEX(0x80, 0x7F, false, false)
}

func TestDEY(t *testing.T) {
	testDEY := func(value uint8, result uint8, isZero bool, isNegative bool) {
		fmt.Printf("testDEY[%d] =? %d\n", value, result)
		testImplicit(0x88, Registers{Y: value})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterY, expected: uint16(result)})
	}

	testDEY(0x01, 0x00, true, false)
	testDEY(0x00, 0xFF, false, true)
	testDEY(0x80, 0x7F, false, false)
}

func TestCMP(t *testing.T) {
	testCMP := func(a uint8, b uint8, isCarry bool, isZero bool, isNegative bool) {
		fmt.Printf("testCMP[%d, %d] =? C:%t Z:%t N:%t\n", a, b, isCarry, isZero, isNegative)
		testImmediate(0xC9, b, Registers{A: a})

		flags := testFlags{}
		flags.Add(FlagCarry, isCarry)
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterA, expected: uint16(a)})
	}

	testCMP(0x01, 0xFF, false, false, false)
	testCMP(0x7F, 0x80, false, false, true)
	testCMP(0x40, 0x20, true, false, false)
	testCMP(0x42, 0x42, true, true, false)
}

func TestCPX(t *testing.T) {
	testCPX := func(a uint8, b uint8, isCarry bool, isZero bool, isNegative bool) {
		fmt.Printf("testCPX[%d, %d] =? C:%t Z:%t N:%t\n", a, b, isCarry, isZero, isNegative)
		testImmediate(0xE0, b, Registers{X: a})

		flags := testFlags{}
		flags.Add(FlagCarry, isCarry)
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterX, expected: uint16(a)})
	}

	testCPX(0x01, 0xFF, false, false, false)
	testCPX(0x7F, 0x80, false, false, true)
	testCPX(0x40, 0x20, true, false, false)
	testCPX(0x42, 0x42, true, true, false)
}

func TestCPY(t *testing.T) {
	testCPY := func(a uint8, b uint8, isCarry bool, isZero bool, isNegative bool) {
		fmt.Printf("testCPY[%d, %d] =? C:%t Z:%t N:%t\n", a, b, isCarry, isZero, isNegative)
		testImmediate(0xC0, b, Registers{Y: a})

		flags := testFlags{}
		flags.Add(FlagCarry, isCarry)
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterY, expected: uint16(a)})
	}

	testCPY(0x01, 0xFF, false, false, false)
	testCPY(0x7F, 0x80, false, false, true)
	testCPY(0x40, 0x20, true, false, false)
	testCPY(0x42, 0x42, true, true, false)
}

func TestTAX(t *testing.T) {
	testTAX := func(value uint8, isZero bool, isNegative bool) {
		fmt.Printf("testTAX[%d] =? Z:%t N:%t\n", value, isZero, isNegative)
		testImplicit(0xAA, Registers{A: value})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterX, expected: uint16(value)})
	}

	testTAX(0x42, false, false)
	testTAX(0x00, true, false)
	testTAX(0x80, false, true)
}

func TestTXA(t *testing.T) {
	testTXA := func(value uint8, isZero bool, isNegative bool) {
		fmt.Printf("testTXA[%d] =? Z:%t N:%t\n", value, isZero, isNegative)
		testImplicit(0x8A, Registers{X: value})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterA, expected: uint16(value)})
	}

	testTXA(0x42, false, false)
	testTXA(0x00, true, false)
	testTXA(0x80, false, true)
}

func TestTAY(t *testing.T) {
	testTAY := func(value uint8, isZero bool, isNegative bool) {
		fmt.Printf("testTAY[%d] =? Z:%t N:%t\n", value, isZero, isNegative)
		testImplicit(0xA8, Registers{A: value})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterY, expected: uint16(value)})
	}

	testTAY(0x42, false, false)
	testTAY(0x00, true, false)
	testTAY(0x80, false, true)
}

func TestTXY(t *testing.T) {
	testTXY := func(value uint8, isZero bool, isNegative bool) {
		fmt.Printf("testTXY[%d] =? Z:%t N:%t\n", value, isZero, isNegative)
		testImplicit(0x98, Registers{Y: value})

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 2, flags, testRegister{register: RegisterA, expected: uint16(value)})
	}

	testTXY(0x42, false, false)
	testTXY(0x00, true, false)
	testTXY(0x80, false, true)
}

func TestBCS(t *testing.T) {
	testBCS := func(offset int8, carryFlag bool, isSuccessful bool, cycles Cycles) {
		fmt.Printf("testBCS[C:%t]\n", carryFlag)
		testRelative(0xB0, offset, Flags{Carry: carryFlag})

		if isSuccessful {
			address := uint16(0x0102 + int16(offset))
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: address})
		} else {
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: 0x102})
		}
	}

	testBCS(+10, false, false, 2)
	testBCS(+10, true, true, 3)
	testBCS(-10, true, true, 4)
}

func TestBCC(t *testing.T) {
	testBCC := func(offset int8, carryFlag bool, isSuccessful bool, cycles Cycles) {
		fmt.Printf("testBCC[C:%t]\n", carryFlag)
		testRelative(0x90, offset, Flags{Carry: carryFlag})

		if isSuccessful {
			address := uint16(0x0102 + int16(offset))
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: address})
		} else {
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: 0x102})
		}
	}

	testBCC(+10, true, false, 2)
	testBCC(+10, false, true, 3)
	testBCC(-10, false, true, 4)
}

func TestBEQ(t *testing.T) {
	testBEQ := func(offset int8, zeroFlag bool, isSuccessful bool, cycles Cycles) {
		fmt.Printf("testBEQ[Z:%t]\n", zeroFlag)
		testRelative(0xF0, offset, Flags{Zero: zeroFlag})

		if isSuccessful {
			address := uint16(0x0102 + int16(offset))
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: address})
		} else {
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: 0x102})
		}
	}

	testBEQ(+10, false, false, 2)
	testBEQ(+10, true, true, 3)
	testBEQ(-10, true, true, 4)
}

func TestBNE(t *testing.T) {
	testBNE := func(offset int8, zeroFlag bool, isSuccessful bool, cycles Cycles) {
		fmt.Printf("testBNE[Z:%t]\n", zeroFlag)
		testRelative(0xD0, offset, Flags{Zero: zeroFlag})

		if isSuccessful {
			address := uint16(0x0102 + int16(offset))
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: address})
		} else {
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: 0x102})
		}
	}

	testBNE(+10, true, false, 2)
	testBNE(+10, false, true, 3)
	testBNE(-10, false, true, 4)
}

func TestBMI(t *testing.T) {
	testBMI := func(offset int8, negativeFlag bool, isSuccessful bool, cycles Cycles) {
		fmt.Printf("testBMI[N:%t]\n", negativeFlag)
		testRelative(0x30, offset, Flags{Negative: negativeFlag})

		if isSuccessful {
			address := uint16(0x0102 + int16(offset))
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: address})
		} else {
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: 0x102})
		}
	}

	testBMI(+10, false, false, 2)
	testBMI(+10, true, true, 3)
	testBMI(-10, true, true, 4)
}

func TestBPL(t *testing.T) {
	testBPL := func(offset int8, negativeFlag bool, isSuccessful bool, cycles Cycles) {
		fmt.Printf("testBPL[N:%t]\n", negativeFlag)
		testRelative(0x10, offset, Flags{Negative: negativeFlag})

		if isSuccessful {
			address := uint16(0x0102 + int16(offset))
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: address})
		} else {
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: 0x102})
		}
	}

	testBPL(+10, true, false, 2)
	testBPL(+10, false, true, 3)
	testBPL(-10, false, true, 4)
}

func TestBVS(t *testing.T) {
	testBVS := func(offset int8, overflowFlag bool, isSuccessful bool, cycles Cycles) {
		fmt.Printf("testBVS[O:%t]\n", overflowFlag)
		testRelative(0x70, offset, Flags{Overflow: overflowFlag})

		if isSuccessful {
			address := uint16(0x0102 + int16(offset))
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: address})
		} else {
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: 0x102})
		}
	}

	testBVS(+10, false, false, 2)
	testBVS(+10, true, true, 3)
	testBVS(-10, true, true, 4)
}

func TestBVC(t *testing.T) {
	testBVC := func(offset int8, overflowFlag bool, isSuccessful bool, cycles Cycles) {
		fmt.Printf("testBVC[O:%t]\n", overflowFlag)
		testRelative(0x50, offset, Flags{Overflow: overflowFlag})

		if isSuccessful {
			address := uint16(0x0102 + int16(offset))
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: address})
		} else {
			assertCPU(t, cycles, testFlags{}, testRegister{register: RegisterPC, expected: 0x102})
		}
	}

	testBVC(+10, true, false, 2)
	testBVC(+10, false, true, 3)
	testBVC(-10, false, true, 4)
}

func TestSTA(t *testing.T) {
	testSTA := func(value uint8) {
		fmt.Printf("testSTA[A=0x%02X]\n", value)
		testAbsolute(0x8D, value, Registers{A: value})

		assertCPU(t, 4, testFlags{}, testRegister{register: RegisterA, expected: uint16(value)})
		assertMemory(t, TestAbsoluteAddress, value)
	}

	testSTA(0x13)
	testSTA(0x37)
}

func TestSTX(t *testing.T) {
	testSTX := func(value uint8) {
		fmt.Printf("testSTX[A=0x%02X]\n", value)
		testAbsolute(0x8E, value, Registers{X: value})

		assertCPU(t, 4, testFlags{}, testRegister{register: RegisterX, expected: uint16(value)})
		assertMemory(t, TestAbsoluteAddress, value)
	}

	testSTX(0x13)
	testSTX(0x37)
}

func TestSTY(t *testing.T) {
	testSTY := func(value uint8) {
		fmt.Printf("testSTY[A=0x%02X]\n", value)
		testAbsolute(0x8C, value, Registers{Y: value})

		assertCPU(t, 4, testFlags{}, testRegister{register: RegisterY, expected: uint16(value)})
		assertMemory(t, TestAbsoluteAddress, value)
	}

	testSTY(0x13)
	testSTY(0x37)
}

func TestCLC(t *testing.T) {
	flags := testFlags{}
	flags.Add(FlagCarry, false)

	testImplicitFlags(0x18, Flags{Carry: true})
	assertCPU(t, 2, flags)
}

func TestSEC(t *testing.T) {
	flags := testFlags{}
	flags.Add(FlagCarry, true)

	testImplicitFlags(0x38, Flags{Carry: false})
	assertCPU(t, 2, flags)
}

func TestCLD(t *testing.T) {
	flags := testFlags{}
	flags.Add(FlagDecimal, false)

	testImplicitFlags(0xD8, Flags{Decimal: true})
	assertCPU(t, 2, flags)
}

func TestSED(t *testing.T) {
	flags := testFlags{}
	flags.Add(FlagDecimal, true)

	testImplicitFlags(0xF8, Flags{Decimal: false})
	assertCPU(t, 2, flags)
}

func TestCLI(t *testing.T) {
	flags := testFlags{}
	flags.Add(FlagInterruptDisable, false)

	testImplicitFlags(0x58, Flags{InterruptDisable: true})
	assertCPU(t, 2, flags)
}

func TestSEI(t *testing.T) {
	flags := testFlags{}
	flags.Add(FlagInterruptDisable, true)

	testImplicitFlags(0x78, Flags{InterruptDisable: false})
	assertCPU(t, 2, flags)
}

func TestCLV(t *testing.T) {
	flags := testFlags{}
	flags.Add(FlagOverflow, false)

	testImplicitFlags(0xB8, Flags{Overflow: true})
	assertCPU(t, 2, flags)
}

func TestLDA(t *testing.T) {
	testLDA := func(value uint8) {
		fmt.Printf("testLDA[0x%02X]\n", value)
		testImmediate(0xA9, value, Registers{})

		assertCPU(t, 2, testFlags{}, testRegister{register: RegisterA, expected: uint16(value)})
	}

	testLDA(0x13)
	testLDA(0x37)
}

func TestLDX(t *testing.T) {
	testLDX := func(value uint8) {
		fmt.Printf("testLDX[0x%02X]\n", value)
		testImmediate(0xA2, value, Registers{})

		assertCPU(t, 2, testFlags{}, testRegister{register: RegisterX, expected: uint16(value)})
	}

	testLDX(0x13)
	testLDX(0x37)
}

func TestLDY(t *testing.T) {
	testLDY := func(value uint8) {
		fmt.Printf("testLDY[0x%02X]\n", value)
		testImmediate(0xA0, value, Registers{})

		assertCPU(t, 2, testFlags{}, testRegister{register: RegisterY, expected: uint16(value)})
	}

	testLDY(0x13)
	testLDY(0x37)
}

func TestJMP(t *testing.T) {
	testAbsoluteAddress(0x4C, 0xBEEF)
	assertCPU(t, 3, testFlags{}, testRegister{register: RegisterPC, expected: 0xBEEF})
}

func TestJSR(t *testing.T) {
	testAbsoluteAddress(0x20, 0xBEEF)
	assertCPU(t, 6, testFlags{},
		testRegister{register: RegisterPC, expected: 0xBEEF},
		testRegister{register: RegisterS, expected: 0xFB},
	)
	assertMemory(t, 0x01FD, 0x01)
	assertMemory(t, 0x01FC, 0x02)
}

func TestRTS(t *testing.T) {
	setup()

	cpu.Registers.PC = 0x0100
	cpu.Registers.S = 0xFB
	cpu.Memory.Poke(0x0100, 0x60)
	cpu.Memory.Poke16(0x01FC, 0x1233)
	cpu.Execute()

	assertCPU(t, 6, testFlags{},
		testRegister{register: RegisterPC, expected: 0x1234},
		testRegister{register: RegisterS, expected: 0xFD},
	)
	assertMemory(t, 0x01FD, 0x12)
	assertMemory(t, 0x01FC, 0x33)
}

func TestPHA(t *testing.T) {
	testImplicit(0x48, Registers{A: 0x42, S: 0xFD})

	assertCPU(t, 3, testFlags{},
		testRegister{register: RegisterS, expected: 0xFC},
	)
	assertMemory(t, 0x01FD, 0x42)
}

func TestPLA(t *testing.T) {
	testPLA := func(value uint8, isZero bool, isNegative bool) {
		fmt.Printf("testPLA[0x%02X] =? Z:%t N:%t\n", value, isZero, isNegative)

		setup()
		cpu.Registers.PC = 0x0100
		cpu.Registers.S = 0xFC
		cpu.Memory.Poke(0x0100, 0x68)
		cpu.Memory.Poke(0x01FD, value)
		cpu.Execute()

		flags := testFlags{}
		flags.Add(FlagZero, isZero)
		flags.Add(FlagNegative, isNegative)

		assertCPU(t, 4, flags, testRegister{register: RegisterA, expected: uint16(value)})
	}

	testPLA(0x42, false, false)
	testPLA(0x00, true, false)
	testPLA(0x80, false, true)
}

func TestPHP(t *testing.T) {
	testPHP := func(flags Flags, expected uint8) {
		fmt.Printf("testPHP[C:%t Z:%t I:%t D:%t B:%d V:%t N:%t]\n",
			flags.Carry, flags.Zero, flags.InterruptDisable, flags.Decimal,
			flags.Origin, flags.Overflow, flags.Negative)
		testImplicitFlags(0x08, flags)

		assertCPU(t, 3, testFlags{})
		assertMemory(t, 0x01FD, expected)
	}

	testPHP(Flags{
		Carry: true, Zero: true, InterruptDisable: true, Decimal: true,
		Origin: FlagOriginPHP, Overflow: true, Negative: true,
	}, 0xFF)
	testPHP(Flags{
		Carry: false, Zero: false, InterruptDisable: false, Decimal: false,
		Origin: FlagOriginNMI, Overflow: false, Negative: false,
	}, 0x20)
}

func TestPLP(t *testing.T) {
	setup()

	cpu.Registers.PC = 0x0100
	cpu.Registers.S = 0xFC
	cpu.Memory.Poke(0x0100, 0x28)
	cpu.Memory.Poke(0x01FD, 0xFF)
	cpu.Execute()

	flags := testFlags{enabled: []Flag{
		FlagCarry, FlagZero, FlagInterruptDisable, FlagDecimal, FlagOverflow, FlagNegative,
	}}

	assertCPU(t, 4, flags)
	if cpu.Flags.Origin != 0x00 {
		t.Errorf("expected flag origin to be 0b00 (ignored by PLP)\n")
	}
}
