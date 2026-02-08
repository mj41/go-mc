package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*PotionContents)(nil)

type PotionContents struct {
	PotionID      pk.Option[pk.VarInt, *pk.VarInt]
	CustomColor   pk.Option[pk.Int, *pk.Int]
	CustomEffects []ItemPotionEffect
	CustomName    pk.Option[pk.String, *pk.String]
}

// ID implements DataComponent.
func (PotionContents) ID() string {
	return "minecraft:potion_contents"
}

// ReadFrom implements DataComponent.
func (p *PotionContents) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		&p.PotionID,
		&p.CustomColor,
		pk.Array(&p.CustomEffects),
		&p.CustomName,
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (p *PotionContents) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		&p.PotionID,
		&p.CustomColor,
		pk.Array(&p.CustomEffects),
		&p.CustomName,
	}.WriteTo(w)
}
