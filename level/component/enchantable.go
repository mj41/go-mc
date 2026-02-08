package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*Enchantable)(nil)

type Enchantable struct {
	pk.VarInt
}

// ID implements DataComponent.
func (Enchantable) ID() string {
	return "minecraft:enchantable"
}
