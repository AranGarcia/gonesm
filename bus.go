package segno

type Bus struct {
	// FIXME: abstract these into proper definitions outside the Bus.
	CPU_RAM   []uint8
	CPU_Stack []uint8

	Cartridge  *Cartridge
	IORegister IORegister
	Mapper     Mapper
}

func (b *Bus) Read(addr uint16) uint8 {
	var data uint8

	switch {
	case addr < 0x2000: // CPU RAM
		// These three are written to 4 different addresses. That's why it's ANDed with 0x07FF.
		// Zero Page	0x0000 - 0x00FF
		// Stack		0x0100 - 0x01FF
		// RAM			0x0200 - 0x07FF
		// Zero Page	0x0800 - 0x08FF (Mirrored from 0x0000 - 0x00FF
		// ...
		data = b.CPU_RAM[addr&0x07FF]
	case addr < 0x2008: // I/O Registers
		data = b.IORegister.Read(addr - 0x2000)
	case addr < 0x4000:
		// Mirrors ($2000-$2007)
	case addr < 0x4020: // I/O Registers
		data = b.IORegister.Read(addr - 0x4000)
	case addr < 0x6000: // Expansion ROM
		data = b.Mapper.ReadExpansion(addr - 0x4020)
	case addr < 0x8000: // SRAM
		data = b.Cartridge.PrgRAM[addr-0x600]
	case addr < 0xC000: // PRG-ROM lower bank
		data = b.Mapper.ReadLowerBank(addr - 0x8000)
	default: // PRG-ROM upper bank (addr < 0x10000)
		data = b.Mapper.ReadUpperBank(addr - 0xC000)
	}

	return data
}

func (b *Bus) Write(addr uint16, data uint8) {
	switch {
	case addr < 0x2000:
		// These three are written to 4 different addresses. That's why it's ANDed with 0x07FF.
		// Zero Page	0x0000 - 0x00FF
		// Stack		0x0100 - 0x01FF
		// RAM			0x0200 - 0x07FF
		// Zero Page	0x0800 - 0x08FF (Mirrored from 0x0000 - 0x00FF)
		// ...
		b.CPU_RAM[addr&0x07FF] = data
	case addr < 0x2008: // I/O Registers
		// TODO: Log: can't write to I/O device
	case addr < 0x4000:
		// Mirrors ($2000-$2007)
	case addr < 0x4020: // I/O Registers
		// TODO: Log: can't write to I/O device
	case addr < 0x6000:
		// TODO: Log: Can't write to ROM
	case addr < 0x8000: // SRAM
		b.Cartridge.PrgRAM[addr-0x6000] = data
	case addr < 0xC000: // PRG-ROM lower bank
		// TODO: Log: Can't write to ROM
	default: // PRG-ROM upper bank (addr < 0x10000)
		// TODO: Log: Can't write to ROM
	}
}

type Mapper interface {
	ReadExpansion(addr uint16) uint8
	ReadLowerBank(addr uint16) uint8
	ReadUpperBank(addr uint16) uint8
}

type IORegister interface {
	Read(addr uint16) uint8
}
