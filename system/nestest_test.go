package system

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"nessie/cartridge"
	"nessie/processor"
	"os"
	"regexp"
	"testing"
)

func TestNESTest(t *testing.T) {
	// Initialize CPU specifically for NESTest
	cpu := processor.NewCPU()
	cpu.Registers.PC = 0xC000
	cpu.Registers.P = 0x24
	cpu.TotalCycles = 7

	// Attempt to load ROM
	rom, err := cartridge.LoadROM("roms/nestest.nes")
	assert.NoError(t, err)

	// Add ROM memory mappings to CPU
	assert.NoError(t, cpu.Memory.AddMappings(rom, processor.MappingCPU))

	// Open golden log and instantiate line-by-line scanning
	logFile, err := os.Open("roms/nestest.log")
	assert.NoError(t, err)
	scanner := bufio.NewScanner(logFile)
	scanner.Split(bufio.ScanLines)

	// Define regular expression for parsing golden log
	logRegexp := regexp.MustCompile(
		`^(?P<PC>[[:xdigit:]]{4})\s+(?P<bytes>(?:[[:xdigit:]]{2} ){0,2}[[:xdigit:]]{2})\s+(?P<instr>.*?)\s+` +
			`A:(?P<A>[[:xdigit:]]{2}) X:(?P<X>[[:xdigit:]]{2}) Y:(?P<Y>[[:xdigit:]]{2}) P:(?P<P>[[:xdigit:]]{2}) ` +
			`SP:(?P<S>[[:xdigit:]]{2}) PPU:\s*(?P<ppu1>\d+),\s*(?P<ppu2>\d+) CYC:(?P<cycles>\d+)`,
	)

	// Iterate until CPU is halted, end of log is reached or a test failed
	for !cpu.Halted && scanner.Scan() && !t.Failed() {
		logLine := scanner.Text()
		logData := mapRegexpSubs(logRegexp.FindStringSubmatch(logLine), logRegexp.SubexpNames())

		assert.Equal(t, logData["PC"], fmt.Sprintf("%04X", cpu.Registers.PC), "unexpected program counter")
		assert.Equal(t, logData["cycles"], fmt.Sprintf("%d", cpu.TotalCycles), "unexpected cycle count")
		assert.Equal(t, logData["A"], fmt.Sprintf("%02X", cpu.Registers.A), "unexpected value of register A")
		assert.Equal(t, logData["X"], fmt.Sprintf("%02X", cpu.Registers.X), "unexpected value of register X")
		assert.Equal(t, logData["Y"], fmt.Sprintf("%02X", cpu.Registers.Y), "unexpected value of register Y")
		assert.Equal(t, logData["P"], fmt.Sprintf("%02X", cpu.Registers.P), "unexpected cpu flags")

		cpu.Execute()
	}
}

func mapRegexpSubs(matches, names []string) map[string]string {
	matches, names = matches[1:], names[1:]
	r := make(map[string]string, len(matches))
	for i := range names {
		r[names[i]] = matches[i]
	}
	return r
}
