package component

import pk "github.com/Tnze/go-mc/net/packet"

var _ DataComponent = (*ItemModel)(nil)

type ItemModel struct {
	pk.String
}

// ID implements DataComponent.
func (ItemModel) ID() string {
	return "minecraft:item_model"
}
