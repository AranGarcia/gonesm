package mappers

type MapperType uint

const (
	UNROM MapperType = iota
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
