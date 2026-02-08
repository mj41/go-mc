package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*ChickenVariant)(nil)

// ChickenVariant component (wire 86).
// Wire: registryEntryHolder<string>. VarInt type: 0=inline string, >0=registry ref (value-1).
type ChickenVariant struct {
	Type       pk.VarInt
	InlineData pk.String // only if Type == 0
}

// ID implements DataComponent.
func (ChickenVariant) ID() string {
	return "minecraft:chicken/variant"
}

// ReadFrom implements DataComponent.
func (c *ChickenVariant) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = c.Type.ReadFrom(r)
	if err != nil {
		return
	}
	if c.Type == 0 {
		n2, err := c.InlineData.ReadFrom(r)
		n += n2
		return n, err
	}
	return
}

// WriteTo implements DataComponent.
func (c *ChickenVariant) WriteTo(w io.Writer) (n int64, err error) {
	n, err = c.Type.WriteTo(w)
	if err != nil {
		return
	}
	if c.Type == 0 {
		n2, err := c.InlineData.WriteTo(w)
		n += n2
		return n, err
	}
	return
}
