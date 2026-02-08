package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*CanPlaceOn)(nil)

type CanPlaceOn struct {
	Predicates []ItemBlockPredicate
}

// ID implements DataComponent.
func (CanPlaceOn) ID() string {
	return "minecraft:can_place_on"
}

// ReadFrom implements DataComponent.
func (c *CanPlaceOn) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&c.Predicates).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (c *CanPlaceOn) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&c.Predicates).WriteTo(w)
}
