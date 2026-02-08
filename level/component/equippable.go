package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Equippable)(nil)

type Equippable struct {
	Slot            pk.VarInt
	Sound           SoundEvent
	Model           pk.Option[pk.String, *pk.String]
	CameraOverlay   pk.Option[pk.String, *pk.String]
	AllowedEntities pk.Option[IDSet, *IDSet]
	Dispensable     pk.Boolean
	Swappable       pk.Boolean
	Damageable      pk.Boolean
	EquipOnInteract pk.Boolean
	Shearable       pk.Boolean
	ShearingSound   SoundEvent
}

// ID implements DataComponent.
func (Equippable) ID() string {
	return "minecraft:equippable"
}

// ReadFrom implements DataComponent.
func (e *Equippable) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		&e.Slot,
		&e.Sound,
		&e.Model,
		&e.CameraOverlay,
		&e.AllowedEntities,
		&e.Dispensable,
		&e.Swappable,
		&e.Damageable,
		&e.EquipOnInteract,
		&e.Shearable,
		&e.ShearingSound,
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (e *Equippable) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		&e.Slot,
		&e.Sound,
		&e.Model,
		&e.CameraOverlay,
		&e.AllowedEntities,
		&e.Dispensable,
		&e.Swappable,
		&e.Damageable,
		&e.EquipOnInteract,
		&e.Shearable,
		&e.ShearingSound,
	}.WriteTo(w)
}
