package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Food)(nil)

type Food struct {
	Nutrition    pk.VarInt
	Saturation   pk.Float
	CanAlwaysEat pk.Boolean
}

// ID implements DataComponent.
func (Food) ID() string {
	return "minecraft:food"
}

// ReadFrom implements DataComponent.
func (f *Food) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{&f.Nutrition, &f.Saturation, &f.CanAlwaysEat}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (f *Food) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{&f.Nutrition, &f.Saturation, &f.CanAlwaysEat}.WriteTo(w)
}
