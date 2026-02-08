package component

import (
	"io"

	"github.com/Tnze/go-mc/chat"
	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Instrument)(nil)

// InstrumentData is the inline data for an instrument when not a registry reference.
// Wire: {soundEvent:ItemSoundHolder, useDuration:f32, range:f32, description:anonymousNbt}
type InstrumentData struct {
	SoundEvent  SoundEvent
	UseDuration pk.Float
	Range       pk.Float
	Description chat.Message
}

func (d *InstrumentData) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&d.SoundEvent, &d.UseDuration, &d.Range, &d.Description}.ReadFrom(r)
}

func (d InstrumentData) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&d.SoundEvent, &d.UseDuration, &d.Range, &d.Description}.WriteTo(w)
}

// Instrument component.
// Wire: {hasHolder:bool, data:switch(hasHolder){true:registryEntryHolder<InstrumentData>, false:string}}
type Instrument struct {
	HasHolder pk.Boolean
	// if HasHolder: registryEntryHolder pattern
	HolderType pk.VarInt      // 0 = inline, >0 = registry ID (value-1)
	InlineData InstrumentData // only if HolderType == 0
	// if !HasHolder: just a string (tag key)
	TagKey pk.String
}

// ID implements DataComponent.
func (Instrument) ID() string {
	return "minecraft:instrument"
}

// ReadFrom implements DataComponent.
func (i *Instrument) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = i.HasHolder.ReadFrom(r)
	if err != nil {
		return
	}
	if i.HasHolder {
		// registryEntryHolder: varint, if 0 read inline, else registry ref (value-1)
		n2, err := i.HolderType.ReadFrom(r)
		n += n2
		if err != nil {
			return n, err
		}
		if i.HolderType == 0 {
			n2, err = i.InlineData.ReadFrom(r)
			n += n2
		}
		return n, err
	}
	n2, err := i.TagKey.ReadFrom(r)
	return n + n2, err
}

// WriteTo implements DataComponent.
func (i *Instrument) WriteTo(w io.Writer) (n int64, err error) {
	n, err = i.HasHolder.WriteTo(w)
	if err != nil {
		return
	}
	if i.HasHolder {
		n2, err := i.HolderType.WriteTo(w)
		n += n2
		if err != nil {
			return n, err
		}
		if i.HolderType == 0 {
			n2, err = i.InlineData.WriteTo(w)
			n += n2
		}
		return n, err
	}
	n2, err := i.TagKey.WriteTo(w)
	return n + n2, err
}

// SoundEvent represents an ItemSoundHolder (registryEntryHolder for sound events).
// Wire: soundId:VarInt. If 0, read inline {soundName:string, fixedRange:Option<f32>}.
// If >0, value is soundId-1 (registry reference).
type SoundEvent struct {
	Type       pk.VarInt
	SoundName  pk.Identifier
	FixedRange pk.Option[pk.Float, *pk.Float]
}

func (s *SoundEvent) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{
		&s.Type,
		pk.Opt{
			Has: func() bool { return s.Type == 0 },
			Field: pk.Tuple{
				&s.SoundName,
				&s.FixedRange,
			},
		},
	}.ReadFrom(r)
}

func (s SoundEvent) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{
		&s.Type,
		pk.Opt{
			Has: func() bool { return s.Type == 0 },
			Field: pk.Tuple{
				&s.SoundName,
				&s.FixedRange,
			},
		},
	}.WriteTo(w)
}
