package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*BlockState)(nil)

type BlockStateProperty struct {
	Name  pk.String
	Value pk.String
}

type BlockState struct {
	Properties []BlockStateProperty
}

// ID implements DataComponent.
func (BlockState) ID() string {
	return "minecraft:block_state"
}

// ReadFrom implements DataComponent.
func (b *BlockState) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&b.Properties).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (b *BlockState) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&b.Properties).WriteTo(w)
}
