package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*StoredEnchantments)(nil)

type StoredEnchantments struct {
	Enchantments []EnchantmentEntry
}

// ID implements DataComponent.
func (StoredEnchantments) ID() string {
	return "minecraft:stored_enchantments"
}

// ReadFrom implements DataComponent.
func (s *StoredEnchantments) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&s.Enchantments).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (s *StoredEnchantments) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&s.Enchantments).WriteTo(w)
}
