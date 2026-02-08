package component

import (
	"io"

	"github.com/Tnze/go-mc/chat"
	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*JukeboxPlayable)(nil)

// JukeboxSongData is inline jukebox song data.
type JukeboxSongData struct {
	SoundEvent       SoundEvent
	Description      chat.Message
	LengthInSeconds  pk.Float
	ComparatorOutput pk.VarInt
}

func (d *JukeboxSongData) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&d.SoundEvent, &d.Description, &d.LengthInSeconds, &d.ComparatorOutput}.ReadFrom(r)
}

func (d JukeboxSongData) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&d.SoundEvent, &d.Description, &d.LengthInSeconds, &d.ComparatorOutput}.WriteTo(w)
}

type JukeboxPlayable struct {
	HasHolder  pk.Boolean
	HolderType pk.VarInt       // registryEntryHolder varint (if HasHolder)
	InlineData JukeboxSongData // only if HasHolder && HolderType == 0
	TagKey     pk.String       // only if !HasHolder
}

// ID implements DataComponent.
func (JukeboxPlayable) ID() string {
	return "minecraft:jukebox_playable"
}

// ReadFrom implements DataComponent.
func (j *JukeboxPlayable) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = j.HasHolder.ReadFrom(r)
	if err != nil {
		return
	}
	if j.HasHolder {
		n2, err := j.HolderType.ReadFrom(r)
		n += n2
		if err != nil {
			return n, err
		}
		if j.HolderType == 0 {
			n2, err = j.InlineData.ReadFrom(r)
			n += n2
		}
		return n, err
	}
	n2, err := j.TagKey.ReadFrom(r)
	return n + n2, err
}

// WriteTo implements DataComponent.
func (j *JukeboxPlayable) WriteTo(w io.Writer) (n int64, err error) {
	n, err = j.HasHolder.WriteTo(w)
	if err != nil {
		return
	}
	if j.HasHolder {
		n2, err := j.HolderType.WriteTo(w)
		n += n2
		if err != nil {
			return n, err
		}
		if j.HolderType == 0 {
			n2, err = j.InlineData.WriteTo(w)
			n += n2
		}
		return n, err
	}
	n2, err := j.TagKey.WriteTo(w)
	return n + n2, err
}
