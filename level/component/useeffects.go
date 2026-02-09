package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*UseEffects)(nil)

// UseEffects component (wire 5).
// Wire: {canSprint:bool, interactVibrations:bool, speedMultiplier:f32}
type UseEffects struct {
	CanSprint          pk.Boolean
	InteractVibrations pk.Boolean
	SpeedMultiplier    pk.Float
}

// ID implements DataComponent.
func (UseEffects) ID() string {
	return "minecraft:use_effects"
}

// ReadFrom implements DataComponent.
func (u *UseEffects) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&u.CanSprint, &u.InteractVibrations, &u.SpeedMultiplier}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (u *UseEffects) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&u.CanSprint, &u.InteractVibrations, &u.SpeedMultiplier}.WriteTo(w)
}
