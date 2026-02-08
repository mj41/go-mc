package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Enchantments)(nil)

type Enchantments struct {
	Enchantments []EnchantmentEntry
}

type EnchantmentEntry struct {
	ID    pk.VarInt
	Level pk.VarInt
}

// ID implements DataComponent.
func (Enchantments) ID() string {
	return "minecraft:enchantments"
}

// ReadFrom implements DataComponent.
func (e *Enchantments) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&e.Enchantments).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (e *Enchantments) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&e.Enchantments).WriteTo(w)
}
