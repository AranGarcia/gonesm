package mappers

type MapperType uint

const (
	NROM MapperType = iota
)

// MMC is the memory management controller interface. Dependending on the implementation,
// read/write operations may or may not include memory bank switching.
type MMC interface {
	// ReadCHR returns data located in the address of CHR memory.
	ReadCHR(uint16) uint8
	// ReadPRG returns data located in the address of PRG memory.
	ReadPRG(uint16) uint8
	// WritePRG saves data in the address of the PRG memory.
	WritePRG(uint16, uint8)
}

type Specs struct {
	// PRGBanks specifies the amount of PRG banks.
	PRGBanks uint
	// CHRBanks specifies the amount of CHR banks.
	CHRBanks uint
}

// NewMapper returns a memory management controller implementation according to the specified
// mapper number.
func NewMapper(mt MapperType, s Specs) MMC {
	var mmc MMC
	switch mt {
	case NROM:
		mmc = &nromMapper{
			prg: make([]byte, 0x4000*s.PRGBanks),
			chr: make([]byte, 0x2000*s.CHRBanks),
		}
	}

	return mmc
}

// nromMapper is MMC implementation for the NROM (no mapper).
type nromMapper struct {
	prg []byte
	chr []byte
}

func (m *nromMapper) ReadCHR(addr uint16) uint8 { return m.chr[addr] }

func (m *nromMapper) ReadPRG(addr uint16) uint8 { return m.prg[addr] }

func (m *nromMapper) WritePRG(addr uint16, data uint8) { m.prg[addr] = data }
