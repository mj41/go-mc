package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*DamageType)(nil)

// DamageType component (wire 8).
// Wire: EitherHolder — Boolean(true) + VarInt (registry ID), or Boolean(false) + String (resource key).
//
// Note: Named DamageType (not DamageTypeComponent) to match the gen-component naming convention.
// The Java class is EitherHolder<net.minecraft.world.damagesource.DamageType>.
type DamageType struct {
	EitherHolder
}

// ID implements DataComponent.
func (DamageType) ID() string {
	return "minecraft:damage_type"
}

// Ensure DamageType satisfies pk.Field via embedded EitherHolder.
var _ pk.Field = (*DamageType)(nil)
