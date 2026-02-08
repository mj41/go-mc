package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*TooltipDisplay)(nil)

type TooltipDisplay struct {
	HideTooltip      pk.Boolean
	HiddenComponents []pk.VarInt
}

// ID implements DataComponent.
func (TooltipDisplay) ID() string {
	return "minecraft:tooltip_display"
}

// ReadFrom implements DataComponent.
func (t *TooltipDisplay) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		&t.HideTooltip,
		pk.Array(&t.HiddenComponents),
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (t *TooltipDisplay) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		&t.HideTooltip,
		pk.Array(&t.HiddenComponents),
	}.WriteTo(w)
}
