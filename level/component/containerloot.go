package component

import (
	"io"

	"github.com/Tnze/go-mc/nbt/dynbt"
	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*ContainerLoot)(nil)

type ContainerLoot struct {
	Data dynbt.Value
}

// ID implements DataComponent.
func (ContainerLoot) ID() string {
	return "minecraft:container_loot"
}

// ReadFrom implements DataComponent.
func (c *ContainerLoot) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.NBTField{V: &c.Data, AllowUnknownFields: true}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (c *ContainerLoot) WriteTo(w io.Writer) (n int64, err error) {
	return pk.NBTField{V: &c.Data, AllowUnknownFields: true}.WriteTo(w)
}
