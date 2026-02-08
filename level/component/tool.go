package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Tool)(nil)

type ToolRule struct {
	Blocks               IDSet
	Speed                pk.Option[pk.Float, *pk.Float]
	CorrectDropForBlocks pk.Option[pk.Boolean, *pk.Boolean]
}

func (r *ToolRule) ReadFrom(rd io.Reader) (n int64, err error) {
	return pk.Tuple{&r.Blocks, &r.Speed, &r.CorrectDropForBlocks}.ReadFrom(rd)
}

func (r ToolRule) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{&r.Blocks, &r.Speed, &r.CorrectDropForBlocks}.WriteTo(w)
}

type Tool struct {
	Rules                      []ToolRule
	DefaultMiningSpeed         pk.Float
	DamagePerBlock             pk.VarInt
	CanDestroyBlocksInCreative pk.Boolean
}

// ID implements DataComponent.
func (Tool) ID() string {
	return "minecraft:tool"
}

// ReadFrom implements DataComponent.
func (t *Tool) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		pk.Array(&t.Rules),
		&t.DefaultMiningSpeed,
		&t.DamagePerBlock,
		&t.CanDestroyBlocksInCreative,
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (t *Tool) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		pk.Array(&t.Rules),
		&t.DefaultMiningSpeed,
		&t.DamagePerBlock,
		&t.CanDestroyBlocksInCreative,
	}.WriteTo(w)
}
