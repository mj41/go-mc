package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*DamageResistant)(nil)

type DamageResistant struct {
	pk.String
}

// ID implements DataComponent.
func (DamageResistant) ID() string {
	return "minecraft:damage_resistant"
}
