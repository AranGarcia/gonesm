package segno

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestCalculateMirroringType(t *testing.T) {
	tests := []struct {
		name string
		bit3 bool
		bit0 bool
		want MirroringType
	}{
		{
			name: "horizontal",
			bit3: false,
			bit0: false,
			want: HorizontalMirroring,
		},
		{
			name: "vertical",
			bit3: false,
			bit0: true,
			want: VerticalMirroring,
		},
		{
			name: "four screen",
			bit3: true,
			bit0: false,
			want: FourScreenMirroring,
		},
		{
			name: "vertical(overwritten) = four screen",
			bit3: true,
			bit0: true,
			want: FourScreenMirroring,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateMirroringType(tt.bit3, tt.bit0); got != tt.want {
				t.Errorf("CalculateMirroringType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadCartridge(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.Reader
		want    *Cartridge
		wantErr bool
	}{
		{
			name:    "empty input",
			wantErr: true,
		},
		{
			name:    "invalid iNES header",
			reader:  strings.NewReader("invalid"),
			wantErr: true,
		},
		{
			name: "no cartridge data",
			reader: bytes.NewBuffer([]byte{
				'N', 'E', 'S', 0x1A,
			}),
			want: &Cartridge{},
		},
		{
			name: "no ROM or RAM",
			reader: bytes.NewBuffer([]byte{
				'N', 'E', 'S', 0x1A,
				0, 0, 240, 240, 0,
			}),
			want: &Cartridge{
				PrgROMBanks:        0,
				ChrROMBanks:        0,
				MirroringType:      HorizontalMirroring,
				HasBatterBackedRAM: false,
				Has512Trainer:      false,
				MapperNumber:       255,
				PRGRAMBanks:        1,

				PrgRAM: make([]byte, PrgRAMBankSize),
				ChrROM: []byte{},
				PrgROM: []byte{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := LoadCartridge(test.reader)
			if (err != nil) != test.wantErr {
				t.Errorf("LoadCartridge() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("LoadCartridge() = %v, want %v", got, test.want)
			}
		})
	}
}
