package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*KineticWeapon)(nil)

// KineticWeaponCondition represents a condition for kinetic weapon effects.
// Wire: {maxDurationTicks:VarInt, minSpeed:f32, minRelativeSpeed:f32}
type KineticWeaponCondition struct {
	MaxDurationTicks pk.VarInt
	MinSpeed         pk.Float
	MinRelativeSpeed pk.Float
}

func (c *KineticWeaponCondition) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&c.MaxDurationTicks, &c.MinSpeed, &c.MinRelativeSpeed}.ReadFrom(r)
}

func (c KineticWeaponCondition) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&c.MaxDurationTicks, &c.MinSpeed, &c.MinRelativeSpeed}.WriteTo(w)
}

// KineticWeapon component (wire 39).
// Wire: {contactCooldownTicks:VarInt, delayTicks:VarInt,
//
//	dismountConditions:Option<Condition>, knockbackConditions:Option<Condition>,
//	damageConditions:Option<Condition>, forwardMovement:f32, damageMultiplier:f32,
//	sound:Option<SoundEvent>, hitSound:Option<SoundEvent>}
type KineticWeapon struct {
	ContactCooldownTicks pk.VarInt
	DelayTicks           pk.VarInt
	DismountConditions   pk.Option[KineticWeaponCondition, *KineticWeaponCondition]
	KnockbackConditions  pk.Option[KineticWeaponCondition, *KineticWeaponCondition]
	DamageConditions     pk.Option[KineticWeaponCondition, *KineticWeaponCondition]
	ForwardMovement      pk.Float
	DamageMultiplier     pk.Float
	Sound                pk.Option[SoundEvent, *SoundEvent]
	HitSound             pk.Option[SoundEvent, *SoundEvent]
}

// ID implements DataComponent.
func (KineticWeapon) ID() string {
	return "minecraft:kinetic_weapon"
}

// ReadFrom implements DataComponent.
func (k *KineticWeapon) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{
		&k.ContactCooldownTicks,
		&k.DelayTicks,
		&k.DismountConditions,
		&k.KnockbackConditions,
		&k.DamageConditions,
		&k.ForwardMovement,
		&k.DamageMultiplier,
		&k.Sound,
		&k.HitSound,
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (k *KineticWeapon) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{
		&k.ContactCooldownTicks,
		&k.DelayTicks,
		&k.DismountConditions,
		&k.KnockbackConditions,
		&k.DamageConditions,
		&k.ForwardMovement,
		&k.DamageMultiplier,
		&k.Sound,
		&k.HitSound,
	}.WriteTo(w)
}
