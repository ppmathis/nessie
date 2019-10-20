package processor

import (
	"fmt"
)

type MappingType int
type Mapping struct {
	From uint32
	To   uint32
}

const DefaultMemorySize = 0x10000
const (
	MappingCPU MappingType = iota
	MappingPPU
)

type Memory interface {
	Reset()
	Peek(address uint16) (value uint8)
	Peek16(address uint16) (value uint16)
	Poke(address uint16, value uint8) (oldValue uint8)
	Poke16(address uint16, value uint16) (oldValue uint16)
}

type MemoryMapper interface {
	Reset()
	Peek(address uint16) (value uint8)
	Poke(address uint16, value uint8) (oldValue uint8)
	Mappings(mapping MappingType) (peek, poke []Mapping)
}

type BasicMemory struct {
	data []uint8
}

type MappedMemory struct {
	Memory
	peek [DefaultMemorySize]MemoryMapper
	poke [DefaultMemorySize]MemoryMapper
}

func NewBasicMemory() *BasicMemory {
	return &BasicMemory{data: make([]uint8, DefaultMemorySize)}
}

func (m *BasicMemory) Reset() {
	for i := range m.data {
		m.data[i] = 0x00
	}
}

func (m *BasicMemory) Peek(address uint16) (value uint8) {
	value = m.data[address]
	return
}

func (m *BasicMemory) Peek16(address uint16) (value uint16) {
	lowByte := m.data[address]
	highByte := m.data[address+1]
	value = (uint16(highByte) << 8) | uint16(lowByte)
	return
}

func (m *BasicMemory) Poke(address uint16, value uint8) (oldValue uint8) {
	oldValue = m.data[address]
	m.data[address] = value
	return
}

func (m *BasicMemory) Poke16(address uint16, value uint16) (oldValue uint16) {
	oldValue = m.Peek16(address)
	m.data[address] = uint8(value & 0xFF)
	m.data[address+1] = uint8((value >> 8) & 0xFF)
	return
}

func NewMappedMemory(base Memory) *MappedMemory {
	return &MappedMemory{Memory: base}
}

func (m *MappedMemory) AddMappings(mapper MemoryMapper, mappingType MappingType) error {
	peekMappings, pokeMappings := mapper.Mappings(mappingType)

	for _, peekMapping := range peekMappings {
		for address := peekMapping.From; address <= peekMapping.To; address++ {
			if m.peek[address] != nil {
				return fmt.Errorf("can not map peek@0x%04X to %v, already in use", address, mapper)
			}

			m.peek[address] = mapper
		}
	}

	for _, pokeMapping := range pokeMappings {
		for address := pokeMapping.From; address <= pokeMapping.To; address++ {
			if m.poke[address] != nil {
				return fmt.Errorf("can not map poke@0x%04X to %v, already in use", address, mapper)
			}

			m.poke[address] = mapper
		}
	}

	return nil
}

func (m *MappedMemory) Peek(address uint16) (value uint8) {
	if mapping := m.peek[address]; mapping != nil {
		value = mapping.Peek(address)
	} else {
		value = m.Memory.Peek(address)
	}

	return
}

func (m *MappedMemory) Peek16(address uint16) (value uint16) {
	lowByte := m.Peek(address)
	highByte := m.Peek(address + 1)
	value = (uint16(highByte) << 8) | uint16(lowByte)
	return
}

func (m *MappedMemory) Poke(address uint16, value uint8) (oldValue uint8) {
	if mapping := m.poke[address]; mapping != nil {
		oldValue = mapping.Poke(address, value)
	} else {
		oldValue = m.Memory.Poke(address, value)
	}

	return
}

func (m *MappedMemory) Poke16(address uint16, value uint16) (oldValue uint16) {
	oldValue = m.Peek16(oldValue)
	m.Poke(address, uint8(value&0xFF))
	m.Poke(address+1, uint8((value>>8)&0xFF))
	return
}

func SamePage(address1 uint16, address2 uint16) bool {
	return (address1^address2)>>8 == 0
}
