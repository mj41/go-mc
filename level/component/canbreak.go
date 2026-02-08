package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*CanBreak)(nil)

type CanBreak struct {
	Predicates []ItemBlockPredicate
}

// ID implements DataComponent.
func (CanBreak) ID() string {
	return "minecraft:can_break"
}

// ReadFrom implements DataComponent.
func (c *CanBreak) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&c.Predicates).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (c *CanBreak) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&c.Predicates).WriteTo(w)
}
