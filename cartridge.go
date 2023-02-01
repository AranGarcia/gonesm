package segno

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"segno/mappers"
)

const (
	ChrROMBankSize = 0x2000
	PrgROMBankSize = 0x4000
	PrgRAMBankSize = 0x2000
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

// MirroringType is used to identify the cartridge's mirroring mechanism.
type MirroringType uint

const (
	HorizontalMirroring MirroringType = iota
	VerticalMirroring
	FourScreenMirroring
)

// CalculateMirroringType determines the mirroring type according to the bit values present in the
// first ROM control byte.
func CalculateMirroringType(bit3, bit0 bool) MirroringType {
	var mt MirroringType
	switch {
	case bit3:
		mt = FourScreenMirroring
	case bit0:
		mt = VerticalMirroring
	default:
		mt = HorizontalMirroring
	}

	return mt
}

// INESPrefix is the magic number used to identify an iNES file.
var INESPrefix = [4]byte{'N', 'E', 'S', 0x1A}

type Cartridge struct {
	// PrgROMBanks describes the amount of 16 KB Program ROM banks contained inside the cartridge.
	PrgROMBanks uint8
	// ChrROMBanks describes the amount of 16 KB Character ROM banks (AKA as VROM) contained inside
	// the cartridge.
	ChrROMBanks uint8
	// MirroringType indicates the type of mirroring used by the game.
	MirroringType MirroringType
	// HasBatterBackedRAM indicates the presence of batter-backed RAM in memory locations
	// $6000-$7FFF
	HasBatterBackedRAM bool
	// Has512Trainer indicates the presence of a 512-byte trainer at memory locations $7000-$71FF
	Has512Trainer bool
	// MapperNumber TODO: IDK what this does yet.
	MapperNumber uint8
	// PRGRAMBanks describes the amount of 8 KB RAM banks.
	PRGRAMBanks uint8

	MMC mappers.MMC

	// FIXME: use MMC instead of these.
	Trainer [512]byte
	ChrROM  []byte
	PrgROM  []byte
	PrgRAM  []byte
}

// LoadCartridge reads data from the input stream and parses the game data.
func LoadCartridge(reader io.Reader) (*Cartridge, error) {
	if reader == nil {
		return nil, ErrInvalidArgument
	}

	var err error
	buff := make([]byte, 4)
	if _, err = reader.Read(buff); err != nil {
		return nil, fmt.Errorf("failed to read prefix; %v", err)
	}

	if c := bytes.Compare(buff, INESPrefix[:]); c != 0 {
		return nil, fmt.Errorf("cartridge data is not in iNES format")
	}

	buff = make([]byte, 5) // Next 5 bytes are ROM and RAM bank metadata
	_, err = reader.Read(buff)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read cartridge metadata; %v", err)
	} else if err == io.EOF {
		return &Cartridge{}, nil
	}

	cartridge := &Cartridge{
		PrgROMBanks:        buff[0],
		ChrROMBanks:        buff[1],
		MirroringType:      CalculateMirroringType(buff[3]&8 == 8, buff[3]&1 == 1),
		HasBatterBackedRAM: buff[2]&2 == 2,
		Has512Trainer:      buff[2]&4 == 4,
		MapperNumber:       (buff[3] & 240) | ((buff[2] & 240) >> 4),
	}
	if buff[4] == 0 {
		cartridge.PRGRAMBanks = 1
	} else {
		cartridge.PRGRAMBanks = buff[4]
	}

	// Load 512-byte trainer
	if cartridge.Has512Trainer {
		reader.Read(cartridge.Trainer[:])
	}

	// FIXME: load CHR and PRG bank data here, not below
	cartridge.MMC = mappers.NewMapper(mappers.MapperType(cartridge.MapperNumber), mappers.Specs{})

	// CHR ROM data
	cartridge.ChrROM = make([]byte, ChrROMBankSize*int(cartridge.ChrROMBanks))
	reader.Read(cartridge.ChrROM)

	// PRG ROM and RAAM
	cartridge.PrgROM = make([]byte, PrgROMBankSize*int(cartridge.PrgROMBanks))
	reader.Read(cartridge.PrgROM)
	cartridge.PrgRAM = make([]byte, PrgRAMBankSize*int(cartridge.PRGRAMBanks))

	return cartridge, nil
}
