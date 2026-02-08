package component

import (
	"io"

	"github.com/Tnze/go-mc/nbt/dynbt"
	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Bees)(nil)

type BeeData struct {
	NBTData        dynbt.Value
	TicksInHive    pk.VarInt
	MinTicksInHive pk.VarInt
}

func (b *BeeData) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		pk.NBTField{V: &b.NBTData, AllowUnknownFields: true},
		&b.TicksInHive,
		&b.MinTicksInHive,
	}.ReadFrom(r)
}

func (b BeeData) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		pk.NBTField{V: &b.NBTData, AllowUnknownFields: true},
		&b.TicksInHive,
		&b.MinTicksInHive,
	}.WriteTo(w)
}

type Bees struct {
	Bees []BeeData
}

// ID implements DataComponent.
func (Bees) ID() string {
	return "minecraft:bees"
}

// ReadFrom implements DataComponent.
func (b *Bees) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&b.Bees).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (b *Bees) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&b.Bees).WriteTo(w)
}
