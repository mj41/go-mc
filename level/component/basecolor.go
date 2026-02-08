package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*BaseColor)(nil)

type BaseColor struct {
	pk.VarInt
}

// ID implements DataComponent.
func (BaseColor) ID() string {
	return "minecraft:base_color"
}
