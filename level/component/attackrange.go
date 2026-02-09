package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*AttackRange)(nil)

// AttackRange component (wire 30).
// Wire: {minRange:f32, maxRange:f32, minCreativeRange:f32, maxCreativeRange:f32, hitboxMargin:f32, mobFactor:f32}
type AttackRange struct {
	MinRange         pk.Float
	MaxRange         pk.Float
	MinCreativeRange pk.Float
	MaxCreativeRange pk.Float
	HitboxMargin     pk.Float
	MobFactor        pk.Float
}

// ID implements DataComponent.
func (AttackRange) ID() string {
	return "minecraft:attack_range"
}

// ReadFrom implements DataComponent.
func (a *AttackRange) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&a.MinRange, &a.MaxRange, &a.MinCreativeRange, &a.MaxCreativeRange, &a.HitboxMargin, &a.MobFactor}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (a *AttackRange) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&a.MinRange, &a.MaxRange, &a.MinCreativeRange, &a.MaxCreativeRange, &a.HitboxMargin, &a.MobFactor}.WriteTo(w)
}
