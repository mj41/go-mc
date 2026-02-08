package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*ProvidesTrimMaterial)(nil)

// ProvidesTrimMaterial component (wire 53).
// Wire: {hasHolder:bool, material:switch(hasHolder){true:registryEntryHolder<ArmorTrimMaterial>, false:string}}
type ProvidesTrimMaterial struct {
	HasHolder pk.Boolean
	// if HasHolder: registryEntryHolder pattern
	HolderType pk.VarInt         // 0 = inline, >0 = registry ID (value-1)
	InlineData ArmorTrimMaterial // only if HolderType == 0
	// if !HasHolder: just a string (tag key)
	TagKey pk.String
}

// ID implements DataComponent.
func (ProvidesTrimMaterial) ID() string {
	return "minecraft:provides_trim_material"
}

// ReadFrom implements DataComponent.
func (p *ProvidesTrimMaterial) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = p.HasHolder.ReadFrom(r)
	if err != nil {
		return
	}
	if p.HasHolder {
		n2, err := p.HolderType.ReadFrom(r)
		n += n2
		if err != nil {
			return n, err
		}
		if p.HolderType == 0 {
			n2, err = p.InlineData.ReadFrom(r)
			n += n2
		}
		return n, err
	}
	n2, err := p.TagKey.ReadFrom(r)
	return n + n2, err
}

// WriteTo implements DataComponent.
func (p *ProvidesTrimMaterial) WriteTo(w io.Writer) (n int64, err error) {
	n, err = p.HasHolder.WriteTo(w)
	if err != nil {
		return
	}
	if p.HasHolder {
		n2, err := p.HolderType.WriteTo(w)
		n += n2
		if err != nil {
			return n, err
		}
		if p.HolderType == 0 {
			n2, err = p.InlineData.WriteTo(w)
			n += n2
		}
		return n, err
	}
	n2, err := p.TagKey.WriteTo(w)
	return n + n2, err
}
