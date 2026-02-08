package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*ChargedProjectiles)(nil)

type ChargedProjectiles struct {
	Projectiles []SlotData
}

// ID implements DataComponent.
func (ChargedProjectiles) ID() string {
	return "minecraft:charged_projectiles"
}

// ReadFrom implements DataComponent.
func (c *ChargedProjectiles) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&c.Projectiles).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (c *ChargedProjectiles) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&c.Projectiles).WriteTo(w)
}
