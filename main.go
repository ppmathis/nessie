package main

import (
	"io/ioutil"
	"nessie/cartridge"
	"nessie/processor"
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

	cpu.Debug = true
	cpu.Registers.PC = 0xC000
	for true {
		cpu.Execute()
	}
}
