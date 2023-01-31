package segno

// CPU6502 is
type CPU6502 struct {
	AccumulatorRegister     uint8
	XIndexRegister          uint8
	YIndexRegister          uint8
	ProcessorStatusRegister uint8
	StackPointer            uint8
	ProgramCounter          uint16

	instructionSet map[uint8]func()
}

func (c *CPU6502) Cycle() {
	c.instructionSet[0]()
}
