package segno

type Bus struct {
	// Data contains the information of the data bus.
	Data uint8
	// Address contains the information of the address bus.
	Address uint16

	// FIXME: abstract these into proper definitions outside the Bus.
	CPU_RAM   []uint8
	CPU_Stack []uint8

	Cartridge *Cartridge
	Mapper    *Mapper
}

func (b *Bus) Read() {
	addr := b.Address

	switch {
	case addr < 0x0100:
		// Zero Page
	case addr < 0x0200: // Stack
		b.Data = b.CPU_Stack[addr-0x0100]
	case addr < 0x0800: // CPU RAM
		b.Data = b.CPU_RAM[addr-0x0200]
	case addr < 0x2000:
		// Mirrors ($0000-$07FF)
	case addr < 0x2008:
		// I/O Registers
	case addr < 0x4000:
		// Mirrors ($2000-$2007)
	case addr < 0x4020:
		// I/O Registers
	case addr < 0x6000:
		// Expansion ROM
	case addr < 0x8000: // SRAM
		b.Data = b.Cartridge.PrgRAM[addr-0x600]
	case addr < 0xC000: // PRG-ROM lower bank
		b.Data = b.Mapper.ReadLowerBank(addr - 0x8000)
	default: // PRG-ROM upper bank (addr < 0x10000)
		b.Data = b.Mapper.ReadUpperBank(addr - 0xC000)
	}
}

func (b *Bus) Write() {
	addr := b.Address
	data := b.Data

	switch {
	case addr < 0x0100:
		// Zero Page
	case addr < 0x0200: // Stack
		b.CPU_Stack[addr-0x0100] = data
	case addr < 0x0800: // CPU RAM
		b.Data = b.CPU_RAM[addr-0x0200]
	case addr < 0x2000:
		// Mirrors ($0000-$07FF)
	case addr < 0x2008:
		// I/O Registers
	case addr < 0x4000:
		// Mirrors ($2000-$2007)
	case addr < 0x4020:
		// I/O Registers
	case addr < 0x6000:
		// Expansion ROM
	case addr < 0x8000: // SRAM
		b.Cartridge.PrgRAM[addr-0x6000] = data
	case addr < 0xC000: // PRG-ROM lower bank
		b.Data = b.Mapper.ReadLowerBank(addr - 0x8000)
	default: // PRG-ROM upper bank (addr < 0x10000)
		b.Data = b.Mapper.ReadUpperBank(addr - 0xC000)
	}
}

// FIXME: support bank switching with Cartridge
type Mapper struct {
	lowerBank []uint8
	upperBank []uint8
}

func (m *Mapper) ReadLowerBank(addr uint16) uint8 { return m.lowerBank[addr] }

func (m *Mapper) ReadUpperBank(addr uint16) uint8 { return m.upperBank[addr] }