package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Consumable)(nil)

type Consumable struct {
	ConsumeSeconds pk.Float
	Animation      pk.VarInt
	Sound          SoundEvent
	MakesParticles pk.Boolean
	Effects        []ItemConsumeEffect
}

// ID implements DataComponent.
func (Consumable) ID() string {
	return "minecraft:consumable"
}

// ReadFrom implements DataComponent.
func (c *Consumable) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		&c.ConsumeSeconds,
		&c.Animation,
		&c.Sound,
		&c.MakesParticles,
		pk.Array(&c.Effects),
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (c *Consumable) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		&c.ConsumeSeconds,
		&c.Animation,
		&c.Sound,
		&c.MakesParticles,
		pk.Array(&c.Effects),
	}.WriteTo(w)
}
