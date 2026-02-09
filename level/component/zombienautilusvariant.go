package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*ZombieNautilusVariant)(nil)

// ZombieNautilusVariant component (wire 94).
// Wire: EitherHolder — Boolean(true) + VarInt (registry ID), or Boolean(false) + String (resource key).
//
// Java: EitherHolder<net.minecraft.world.entity.animal.nautilus.ZombieNautilusVariant>.
type ZombieNautilusVariant struct {
	EitherHolder
}

// ID implements DataComponent.
func (ZombieNautilusVariant) ID() string {
	return "minecraft:zombie_nautilus/variant"
}

var _ pk.Field = (*ZombieNautilusVariant)(nil)
