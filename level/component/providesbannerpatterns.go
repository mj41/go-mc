package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*ProvidesBannerPatterns)(nil)

type ProvidesBannerPatterns struct {
	pk.String
}

// ID implements DataComponent.
func (ProvidesBannerPatterns) ID() string {
	return "minecraft:provides_banner_patterns"
}
