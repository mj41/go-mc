package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*CustomModelData)(nil)

type CustomModelData struct {
	Floats  []pk.Float
	Flags   []pk.Boolean
	Strings []pk.String
	Colors  []pk.Int
}

// ID implements DataComponent.
func (CustomModelData) ID() string {
	return "minecraft:custom_model_data"
}

// ReadFrom implements DataComponent.
func (c *CustomModelData) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		pk.Array(&c.Floats),
		pk.Array(&c.Flags),
		pk.Array(&c.Strings),
		pk.Array(&c.Colors),
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (c *CustomModelData) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		pk.Array(&c.Floats),
		pk.Array(&c.Flags),
		pk.Array(&c.Strings),
		pk.Array(&c.Colors),
	}.WriteTo(w)
}
