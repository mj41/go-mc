package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*MinimumAttackCharge)(nil)

// MinimumAttackCharge component (wire 7).
// Wire: f32 (single float, range 0.0-1.0).
type MinimumAttackCharge struct {
	pk.Float
}

// ID implements DataComponent.
func (MinimumAttackCharge) ID() string {
	return "minecraft:minimum_attack_charge"
}
