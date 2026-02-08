package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Trim)(nil)

type Trim struct {
	// Material: registryEntryHolder<ArmorTrimMaterial>
	MaterialType pk.VarInt
	MaterialData ArmorTrimMaterial // only if MaterialType == 0
	// Pattern: registryEntryHolder<ArmorTrimPattern>
	PatternType pk.VarInt
	PatternData ArmorTrimPattern // only if PatternType == 0
}

// ID implements DataComponent.
func (Trim) ID() string {
	return "minecraft:trim"
}

// ReadFrom implements DataComponent.
func (t *Trim) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		&t.MaterialType,
		pk.Opt{
			Has:   func() bool { return t.MaterialType == 0 },
			Field: &t.MaterialData,
		},
		&t.PatternType,
		pk.Opt{
			Has:   func() bool { return t.PatternType == 0 },
			Field: &t.PatternData,
		},
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (t *Trim) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		&t.MaterialType,
		pk.Opt{
			Has:   func() bool { return t.MaterialType == 0 },
			Field: &t.MaterialData,
		},
		&t.PatternType,
		pk.Opt{
			Has:   func() bool { return t.PatternType == 0 },
			Field: &t.PatternData,
		},
	}.WriteTo(w)
}
