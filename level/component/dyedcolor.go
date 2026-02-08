package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*DyedColor)(nil)

type DyedColor struct {
	pk.Int
}

// ID implements DataComponent.
func (DyedColor) ID() string {
	return "minecraft:dyed_color"
}
