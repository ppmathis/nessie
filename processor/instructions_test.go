package processor

import (
	"fmt"
	"testing"
)

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

func assertCPU(t *testing.T, expectedCycles Cycles, flags testFlags, registers ...testRegister) {
	assertCycles(t, expectedCycles)
	assertFlags(t, flags.enabled, flags.disabled)

	for _, assert := range registers {
		assertRegister(t, assert.register, assert.expected)
	}
}

func testImmediate(opcode uint8, argument uint8, registers Registers) {
	setup()

	cpu.Registers = registers
	cpu.Registers.PC = 0x0100
	cpu.Memory.Poke(0x0100, opcode)
	cpu.Memory.Poke(0x0101, argument)
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
