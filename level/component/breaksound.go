package component

import (
	"io"
)

var _ DataComponent = (*BreakSound)(nil)

// BreakSound is an ItemSoundHolder (registryEntryHolder for sound events).
type BreakSound struct {
	Sound SoundEvent
}

// ID implements DataComponent.
func (BreakSound) ID() string {
	return "minecraft:break_sound"
}

// ReadFrom implements DataComponent.
func (b *BreakSound) ReadFrom(r io.Reader) (n int64, err error) {
	return b.Sound.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (b *BreakSound) WriteTo(w io.Writer) (n int64, err error) {
	return b.Sound.WriteTo(w)
}
