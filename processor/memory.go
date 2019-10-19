package processor

const DefaultMemorySize = 0x10000

type Memory interface {
	Clear()
	Peek(address uint16) (value uint8)
	Peek16(address uint16) (value uint16)
	Poke(address uint16, value uint8) (oldValue uint8)
	Poke16(address uint16, value uint16) (oldValue uint16)
}

type BasicMemory struct {
	data []uint8
}

func NewBasicMemory() *BasicMemory {
	return &BasicMemory{data: make([]uint8, DefaultMemorySize)}
}

func (m *BasicMemory) Clear() {
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

func SamePage(address1 uint16, address2 uint16) bool {
	return (address1^address2)>>8 == 0
}
