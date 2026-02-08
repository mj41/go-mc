package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*PotDecorations)(nil)

type PotDecorations struct {
	Decorations []pk.VarInt
}

// ID implements DataComponent.
func (PotDecorations) ID() string {
	return "minecraft:pot_decorations"
}

// ReadFrom implements DataComponent.
func (p *PotDecorations) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&p.Decorations).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (p *PotDecorations) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&p.Decorations).WriteTo(w)
}
