package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*DeathProtection)(nil)

type DeathProtection struct {
	Effects []ItemConsumeEffect
}

// ID implements DataComponent.
func (DeathProtection) ID() string {
	return "minecraft:death_protection"
}

// ReadFrom implements DataComponent.
func (d *DeathProtection) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&d.Effects).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (d *DeathProtection) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&d.Effects).WriteTo(w)
}
