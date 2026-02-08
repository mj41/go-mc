package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Container)(nil)

type Container struct {
	Contents []SlotData
}

// ID implements DataComponent.
func (Container) ID() string {
	return "minecraft:container"
}

// ReadFrom implements DataComponent.
func (c *Container) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&c.Contents).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (c *Container) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&c.Contents).WriteTo(w)
}
