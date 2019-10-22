package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"nessie/cartridge"
	"nessie/processor"
	"os"
	"regexp"
)

func main() {
	data, err := ioutil.ReadFile("roms/nestest.nes")
	if err != nil {
		panic(err)
	}

	rom, err := cartridge.NewROM(data)
	if err != nil {
		panic(err)
	}

	cpu := processor.NewCPU()
	if err := cpu.Memory.AddMappings(rom, processor.MappingCPU); err != nil {
		panic(err)
	}

	file, err := os.Open("roms/nestest.log")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	testRegex := regexp.MustCompile(
		`^(?P<PC>[[:xdigit:]]{4})\s+(?P<bytes>(?:[[:xdigit:]]{2} ){0,2}[[:xdigit:]]{2})\s+(?P<instr>.*?)\s+` +
			`A:(?P<A>[[:xdigit:]]{2}) X:(?P<X>[[:xdigit:]]{2}) Y:(?P<Y>[[:xdigit:]]{2}) P:(?P<P>[[:xdigit:]]{2}) ` +
			`SP:(?P<S>[[:xdigit:]]{2}) PPU:\s*(?P<ppu1>\d+),\s*(?P<ppu2>\d+) CYC:(?P<cyc>\d+)`)

	cpu.Debug = true
	cpu.Registers.PC = 0xC000
	cpu.Registers.P = 0x24
	cpu.TotalCycles = 7

	for !cpu.Halted && scanner.Scan() {
		testLine := scanner.Text()
		testData := mapRegexpSubs(testRegex.FindStringSubmatch(testLine), testRegex.SubexpNames())

		fmt.Printf("NESTest: [0x%s] %-8s - %-40s - A:%s X:%s Y:%s S:%s P:%s CYC:%s\n",
			testData["PC"], testData["bytes"], testData["instr"],
			testData["A"], testData["X"], testData["Y"], testData["S"], testData["P"], testData["cyc"],
		)
		if act, exp := fmt.Sprintf("%04X", cpu.Registers.PC), testData["PC"]; exp != act {
			fmt.Printf("NESTest: program counter mismatch, expected 0x%s != 0x%s actual\n", exp, act)
			panic(fmt.Errorf("test mismatch"))
		}
		if act, exp := fmt.Sprintf("%d", cpu.TotalCycles), testData["cyc"]; exp != act {
			fmt.Printf("NESTest: cycle count mismatch, expected %s != %s actual\n", exp, act)
			panic(fmt.Errorf("test mismatch"))
		}
		if act, exp := fmt.Sprintf("%02X", cpu.Registers.A), testData["A"]; exp != act {
			fmt.Printf("NESTest: register A mismatch, expected 0x%s != 0x%s actual\n", exp, act)
			panic(fmt.Errorf("test mismatch"))
		}
		if act, exp := fmt.Sprintf("%02X", cpu.Registers.X), testData["X"]; exp != act {
			fmt.Printf("NESTest: register X mismatch, expected 0x%s != 0x%s actual\n", exp, act)
			panic(fmt.Errorf("test mismatch"))
		}
		if act, exp := fmt.Sprintf("%02X", cpu.Registers.Y), testData["Y"]; exp != act {
			fmt.Printf("NESTest: register Y mismatch, expected 0x%s != 0x%s actual\n", exp, act)
			panic(fmt.Errorf("test mismatch"))
		}
		if act, exp := fmt.Sprintf("%02X", cpu.Registers.P), testData["P"]; exp != act {
			fmt.Printf("NESTest: cpu flags mismatch, expected 0x%s != 0x%s actual\n", exp, act)
			panic(fmt.Errorf("test mismatch"))
		}

		cpu.Execute()
		fmt.Printf("Nessie:  %s\n", cpu.Disassembly())
	}

	result1 := cpu.Memory.Peek(0x02)
	result2 := cpu.Memory.Peek(0x03)
	fmt.Printf("NESTest Result: 0x%02X 0x%02X", result1, result2)
}

func mapRegexpSubs(matches, names []string) map[string]string {
	matches, names = matches[1:], names[1:]
	r := make(map[string]string, len(matches))
	for i := range names {
		r[names[i]] = matches[i]
	}
	return r
}
