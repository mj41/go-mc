package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*ChickenVariant)(nil)

// ChickenVariant component (wire 93).
// Wire: EitherHolder — Boolean(true) + VarInt (registry ID), or Boolean(false) + String (resource key).
//
// Java: EitherHolder<net.minecraft.world.entity.animal.chicken.ChickenVariant>.
// Uses ByteBufCodecs.either(holderRegistry, ResourceKey.streamCodec).
type ChickenVariant struct {
	EitherHolder
}

// ID implements DataComponent.
func (ChickenVariant) ID() string {
	return "minecraft:chicken/variant"
}

var _ pk.Field = (*ChickenVariant)(nil)
