package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*NoteBlockSound)(nil)

type NoteBlockSound struct {
	pk.String
}

// ID implements DataComponent.
func (NoteBlockSound) ID() string {
	return "minecraft:note_block_sound"
}
