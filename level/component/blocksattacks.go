package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*BlocksAttacks)(nil)

// DamageReduction represents a damage reduction entry.
type DamageReduction struct {
	HorizontalBlockingAngle pk.Float
	Type                    pk.Option[IDSet, *IDSet]
	Base                    pk.Float
	Factor                  pk.Float
}

func (d *DamageReduction) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&d.HorizontalBlockingAngle, &d.Type, &d.Base, &d.Factor}.ReadFrom(r)
}

func (d DamageReduction) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&d.HorizontalBlockingAngle, &d.Type, &d.Base, &d.Factor}.WriteTo(w)
}

// ItemDamageFunction represents item damage parameters.
type ItemDamageFunction struct {
	Threshold pk.Float
	Base      pk.Float
	Factor    pk.Float
}

func (d *ItemDamageFunction) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&d.Threshold, &d.Base, &d.Factor}.ReadFrom(r)
}

func (d ItemDamageFunction) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&d.Threshold, &d.Base, &d.Factor}.WriteTo(w)
}

// BlocksAttacks component (wire 33).
// Wire: {blockDelaySeconds:f32, disableCooldownScale:f32, damageReductions:Array<DamageReduction>,
//
//	itemDamage:ItemDamageFunction, bypassedBy:Option<string>,
//	blockSound:Option<SoundEvent>, disableSound:Option<SoundEvent>}
type BlocksAttacks struct {
	BlockDelaySeconds    pk.Float
	DisableCooldownScale pk.Float
	DamageReductions     []DamageReduction
	ItemDamage           ItemDamageFunction
	BypassedBy           pk.Option[pk.String, *pk.String]
	BlockSound           pk.Option[SoundEvent, *SoundEvent]
	DisableSound         pk.Option[SoundEvent, *SoundEvent]
}

// ID implements DataComponent.
func (BlocksAttacks) ID() string {
	return "minecraft:blocks_attacks"
}

// ReadFrom implements DataComponent.
func (b *BlocksAttacks) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		&b.BlockDelaySeconds,
		&b.DisableCooldownScale,
		pk.Array(&b.DamageReductions),
		&b.ItemDamage,
		&b.BypassedBy,
		&b.BlockSound,
		&b.DisableSound,
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (b *BlocksAttacks) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		&b.BlockDelaySeconds,
		&b.DisableCooldownScale,
		pk.Array(&b.DamageReductions),
		&b.ItemDamage,
		&b.BypassedBy,
		&b.BlockSound,
		&b.DisableSound,
	}.WriteTo(w)
}
