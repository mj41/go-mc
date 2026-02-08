package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*FireworkExplosion)(nil)

type FireworkExplosion struct {
	ItemFireworkExplosion
}

// ID implements DataComponent.
func (FireworkExplosion) ID() string {
	return "minecraft:firework_explosion"
}

// ReadFrom implements DataComponent.
func (f *FireworkExplosion) ReadFrom(r io.Reader) (n int64, err error) {
	return f.ItemFireworkExplosion.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (f *FireworkExplosion) WriteTo(w io.Writer) (n int64, err error) {
	return f.ItemFireworkExplosion.WriteTo(w)
}

var _ DataComponent = (*Fireworks)(nil)

type Fireworks struct {
	FlightDuration pk.VarInt
	Explosions     []ItemFireworkExplosion
}

// ID implements DataComponent.
func (Fireworks) ID() string {
	return "minecraft:fireworks"
}

// ReadFrom implements DataComponent.
func (f *Fireworks) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		&f.FlightDuration,
		pk.Array(&f.Explosions),
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (f *Fireworks) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		&f.FlightDuration,
		pk.Array(&f.Explosions),
	}.WriteTo(w)
}
