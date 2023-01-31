package mappers

type MapperType uint

const (
	UNROM MapperType = iota
	CNROM
	MMC1
	MMC3
)

type Mapper interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}
