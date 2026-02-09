package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*SwingAnimation)(nil)

// SwingAnimation component (wire 40).
// Wire: {type:VarInt (0=NONE, 1=WHACK, 2=STAB), duration:VarInt}
type SwingAnimation struct {
	Type     pk.VarInt
	Duration pk.VarInt
}

// ID implements DataComponent.
func (SwingAnimation) ID() string {
	return "minecraft:swing_animation"
}

// ReadFrom implements DataComponent.
func (s *SwingAnimation) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&s.Type, &s.Duration}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (s *SwingAnimation) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&s.Type, &s.Duration}.WriteTo(w)
}
