package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*PiercingWeapon)(nil)

// PiercingWeapon component (wire 38).
// Wire: {dealsKnockback:bool, dismounts:bool, sound:Option<SoundEvent>, hitSound:Option<SoundEvent>}
type PiercingWeapon struct {
	DealsKnockback pk.Boolean
	Dismounts      pk.Boolean
	Sound          pk.Option[SoundEvent, *SoundEvent]
	HitSound       pk.Option[SoundEvent, *SoundEvent]
}

// ID implements DataComponent.
func (PiercingWeapon) ID() string {
	return "minecraft:piercing_weapon"
}

// ReadFrom implements DataComponent.
func (p *PiercingWeapon) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&p.DealsKnockback, &p.Dismounts, &p.Sound, &p.HitSound}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (p *PiercingWeapon) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&p.DealsKnockback, &p.Dismounts, &p.Sound, &p.HitSound}.WriteTo(w)
}
