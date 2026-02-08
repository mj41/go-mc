package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*PotionDurationScale)(nil)

type PotionDurationScale struct {
	pk.Float
}

// ID implements DataComponent.
func (PotionDurationScale) ID() string {
	return "minecraft:potion_duration_scale"
}
