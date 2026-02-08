package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*UseCooldown)(nil)

type UseCooldown struct {
	Seconds       pk.Float
	CooldownGroup pk.Option[pk.String, *pk.String]
}

// ID implements DataComponent.
func (UseCooldown) ID() string {
	return "minecraft:use_cooldown"
}

// ReadFrom implements DataComponent.
func (u *UseCooldown) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{&u.Seconds, &u.CooldownGroup}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (u *UseCooldown) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{&u.Seconds, &u.CooldownGroup}.WriteTo(w)
}
