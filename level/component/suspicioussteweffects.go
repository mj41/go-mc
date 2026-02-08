package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*SuspiciousStewEffects)(nil)

type StewEffect struct {
	Effect   pk.VarInt
	Duration pk.VarInt
}

type SuspiciousStewEffects struct {
	Effects []StewEffect
}

// ID implements DataComponent.
func (SuspiciousStewEffects) ID() string {
	return "minecraft:suspicious_stew_effects"
}

// ReadFrom implements DataComponent.
func (s *SuspiciousStewEffects) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&s.Effects).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (s *SuspiciousStewEffects) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&s.Effects).WriteTo(w)
}
