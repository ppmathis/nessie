package main

import (
	"fmt"
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

	fmt.Printf("0xC000 = 0x%02X", cpu.Memory.Peek(0xC000))
}
