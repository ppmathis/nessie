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

func testDEC(t *testing.T) {
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
		fmt.Printf("testCMP[%d, %d] =? C:%v Z:%v N:%v\n", a, b, isCarry, isZero, isNegative)
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
		fmt.Printf("testCPX[%d, %d] =? C:%v Z:%v N:%v\n", a, b, isCarry, isZero, isNegative)
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
		fmt.Printf("testCPY[%d, %d] =? C:%v Z:%v N:%v\n", a, b, isCarry, isZero, isNegative)
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
		fmt.Printf("testTAX[%d] =? Z:%v N:%v\n", value, isZero, isNegative)
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
		fmt.Printf("testTXA[%d] =? Z:%v N:%v\n", value, isZero, isNegative)
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
		fmt.Printf("testTAY[%d] =? Z:%v N:%v\n", value, isZero, isNegative)
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
		fmt.Printf("testTXY[%d] =? Z:%v N:%v\n", value, isZero, isNegative)
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

