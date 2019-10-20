package cartridge

import "nessie/processor"

type NROM struct {
	*ROMFile
}

func NewNROM(romFile *ROMFile) *NROM {
	return &NROM{ROMFile: romFile}
}

func (nrom *NROM) Mappings(mappingType processor.MappingType) (peek, poke []processor.Mapping) {
	switch mappingType {
	case processor.MappingCPU:
		if nrom.BankCountPRG > 0 {
			peek = append(peek, processor.Mapping{From: 0x8000, To: 0xFFFF})
		}

	case processor.MappingPPU:
		if nrom.BankCountCHR > 0 {
			peek = append(peek, processor.Mapping{From: 0x0000, To: 0x1FFF})
			poke = append(poke, processor.Mapping{From: 0x0000, To: 0x1FFF})
		}
	}

	return
}

func (nrom *NROM) Reset() {
}

func (nrom *NROM) Peek(address uint16) (value uint8) {
	switch {
	// PPU Memory
	case address <= 0x1FFF && nrom.BankCountCHR > 0:
		value = nrom.BanksCHR[0][address]

	// CPU Memory
	case address >= 0x8000 && nrom.BankCountPRG > 0:
		bankAddress := address & 0x3FFF

		switch {
		// First PRG Bank
		case address >= 0x8000 && address <= 0xBFFF:
			value = nrom.BanksPRG[0][bankAddress]
		// Last PRG Bank
		case address >= 0xc000:
			value = nrom.BanksPRG[nrom.BankCountPRG-1][bankAddress]
		}
	}

	return
}

func (nrom *NROM) Poke(address uint16, value uint8) (oldValue uint8) {
	switch {
	// PPU Banks
	case address <= 0x1FFF && nrom.BankCountCHR > 0:
		nrom.BanksCHR[0][address] = value
	}

	return
}
