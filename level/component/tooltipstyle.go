package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*TooltipStyle)(nil)

type TooltipStyle struct {
	pk.String
}

// ID implements DataComponent.
func (TooltipStyle) ID() string {
	return "minecraft:tooltip_style"
}
