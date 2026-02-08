package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*Weapon)(nil)

type Weapon struct {
	ItemDamagePerAttack       pk.VarInt
	DisableBlockingForSeconds pk.Float
}

// ID implements DataComponent.
func (Weapon) ID() string {
	return "minecraft:weapon"
}

// ReadFrom implements DataComponent.
func (w *Weapon) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{&w.ItemDamagePerAttack, &w.DisableBlockingForSeconds}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (wp *Weapon) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{&wp.ItemDamagePerAttack, &wp.DisableBlockingForSeconds}.WriteTo(w)
}
