package segno

import (
	"bytes"
	"log"
)

// NESComponents are the digital devices that are part of the NES architecture.
type NESComponents struct {
	// Cartridge represents the port where the Cartridge is connected to the main bus of the NES.
	Cartridge *Cartridge
	// Mapper is the memory mapper that switches banks of the cartridge's ROMs.
	Mapper Mapper
}

// Bus is the communication channel between the NESComponents.
type Bus struct {
	logger *log.Logger

	cpuRAM     [2048]byte
	cartridge  *Cartridge
	ioRegister IORegister
	mapper     Mapper
}

// NewBus returns an initialized Bus that takes the components and maps them into the internal
// state. Additionally, it initializes a logger.
func NewBus(components *NESComponents) *Bus {
	return &Bus{
		logger:    log.New(&bytes.Buffer{}, "bus: ", log.Lmicroseconds),
		cartridge: components.Cartridge,
		mapper:    components.Mapper,
	}
}

// Read attempts to read a device located in the address.
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
		data = b.cpuRAM[addr&0x07FF]
	case addr < 0x2008: // I/O Registers
		data = b.ioRegister.Read(addr - 0x2000)
	case addr < 0x4000:
		// Mirrors ($2000-$2007)
	case addr < 0x4020: // I/O Registers
		data = b.ioRegister.Read(addr - 0x4000)
	case addr < 0x6000: // Expansion ROM
		data = b.mapper.ReadExpansion(addr - 0x4020)
	case addr < 0x8000: // SRAM
		data = b.cartridge.PrgRAM[addr-0x600]
	case addr < 0xC000: // PRG-ROM lower bank
		data = b.mapper.ReadLowerBank(addr - 0x8000)
	default: // PRG-ROM upper bank (addr < 0x10000)
		data = b.mapper.ReadUpperBank(addr - 0xC000)
	}

	return data
}

// Write attempts to write the data into the component mapped to the address.
func (b *Bus) Write(addr uint16, data uint8) {
	switch {
	case addr < 0x2000:
		// These three are written to 4 different addresses. That's why it's ANDed with 0x07FF.
		// Zero Page	0x0000 - 0x00FF
		// Stack		0x0100 - 0x01FF
		// RAM			0x0200 - 0x07FF
		// Zero Page	0x0800 - 0x08FF (Mirrored from 0x0000 - 0x00FF)
		// ...
		b.cpuRAM[addr&0x07FF] = data
	case addr < 0x2008: // I/O Registers
		b.logger.Printf("attempted to write to I/O device in address %x", addr)
	case addr < 0x4000:
		// Mirrors ($2000-$2007)
	case addr < 0x4020: // I/O Registers
		b.logger.Printf("attempted to write to I/O device in address %x", addr)
	case addr < 0x6000:
		b.logger.Printf("attempted to write to Expansion ROM in address %x", addr)
	case addr < 0x8000: // SRAM
		b.cartridge.PrgRAM[addr-0x6000] = data
	case addr < 0xC000: // PRG-ROM lower bank
		b.logger.Printf("attempted to write to PRG ROM in address %x", addr)
	default: // PRG-ROM upper bank (addr < 0x10000)
		b.logger.Printf("attempted to write to PRG ROM in address %x", addr)
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
