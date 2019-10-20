package cartridge

import (
	"errors"
	"fmt"
	"nessie/processor"
)

const headerLength = 16
const trainerLength = 512
const prgBankLength = 16384
const chrBankLength = 8192

const f6HasBattery = 1 << 1
const f6HasTrainer = 1 << 2
const f6ForceFourScreen = 1 << 3
const f6MapperLow = 0xF << 4
const f7MapperHigh = 0xF << 4

type ROM interface {
	processor.MemoryMapper
}

type ROMFile struct {
	BankCountPRG    uint8
	BankCountCHR    uint8
	HasBattery      bool
	HasTrainer      bool
	ForceFourScreen bool
	MapperID        uint8

	Trainer  []byte
	BanksPRG [][]byte
	BanksCHR [][]byte
	BanksRAM [][]byte
}

func NewROM(buffer []byte) (ROM, error) {
	romFile, err := NewROMFile(buffer)
	if err != nil {
		return nil, fmt.Errorf("could not load rom file: %v", err)
	}

	switch romFile.MapperID {
	case 0x00, 0x40, 0x41:
		return NewNROM(romFile), nil
	default:
		return nil, fmt.Errorf("unsupported mapper type: 0x%02X", romFile.MapperID)
	}
}

func NewROMFile(buffer []byte) (*ROMFile, error) {
	romFile := &ROMFile{}

	if err := romFile.parseHeader(buffer); err != nil {
		return nil, err
	}
	if err := romFile.loadTrainer(buffer); err != nil {
		return nil, err
	}
	if err := romFile.loadBanks(buffer); err != nil {
		return nil, err
	}

	return romFile, nil
}

func (r *ROMFile) String() string {
	return fmt.Sprintf("ROMFile[M=%d,PRG=%d,CHR=%d]", r.MapperID, r.BankCountPRG, r.BankCountCHR)
}

func (r *ROMFile) parseHeader(buffer []byte) error {
	// Verify complete iNES header is available
	if len(buffer) < headerLength {
		return errors.New("missing 16 byte header")
	}

	// Check if magic constant (first 3 bytes) is correct
	if string(buffer[0:3]) != "NES" || buffer[3] != 0x1A {
		return errors.New("missing NES magic constant")
	}

	// Parse iNES header
	r.BankCountPRG = buffer[4]
	r.BankCountCHR = buffer[5]
	r.HasBattery = (buffer[6] & f6HasBattery) == f6HasBattery
	r.HasTrainer = (buffer[6] & f6HasTrainer) == f6HasTrainer
	r.ForceFourScreen = (buffer[6] & f6ForceFourScreen) == f6ForceFourScreen
	r.MapperID = (buffer[6] & f6MapperLow) | ((buffer[7] & f7MapperHigh) << 4)

	return nil
}

func (r *ROMFile) loadTrainer(buffer []byte) error {
	// Skip if no trainer is available
	if !r.HasTrainer {
		return nil
	}

	// Ensure we have enough data for trainer
	if len(buffer) < (headerLength + trainerLength) {
		return errors.New("not enough bytes available for trainer data")
	}

	// Load trainer data
	r.Trainer = buffer[headerLength : headerLength+trainerLength]
	return nil
}

func (r *ROMFile) loadBanks(buffer []byte) error {
	// Calculate offset to first BanksPRG bank
	offset := headerLength
	if r.HasTrainer {
		offset += trainerLength
	}

	// Ensure we have enough data for BanksPRG  banks available
	if len(buffer) < (offset + (int(r.BankCountPRG) * prgBankLength)) {
		return errors.New("not enough bytes available for BanksPRG data")
	}

	// Load data into BanksPRG banks
	r.BanksPRG = make([][]byte, r.BankCountPRG)
	for bank := 0; bank < int(r.BankCountPRG); bank++ {
		r.BanksPRG[bank] = buffer[offset : offset+prgBankLength]
		offset += prgBankLength
	}

	// Ensure we have enough data for BanksCHR banks available
	if len(buffer) < (offset + (int(r.BankCountCHR) * chrBankLength)) {
		return errors.New("not enough bytes available for BanksCHR data")
	}

	// Load data into BanksCHR banks
	r.BanksCHR = make([][]byte, r.BankCountCHR)
	for bank := 0; bank < int(r.BankCountCHR); bank++ {
		r.BanksCHR[bank] = buffer[offset : offset+chrBankLength]
		offset += chrBankLength
	}

	return nil
}
