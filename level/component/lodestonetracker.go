package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*LodestoneTracker)(nil)

type GlobalPos struct {
	DimensionName pk.Identifier
	Location      pk.Position
}

func (g *GlobalPos) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&g.DimensionName, &g.Location}.ReadFrom(r)
}

func (g GlobalPos) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&g.DimensionName, &g.Location}.WriteTo(w)
}

type LodestoneTracker struct {
	GlobalPosition pk.Option[GlobalPos, *GlobalPos]
	Tracked        pk.Boolean
}

// ID implements DataComponent.
func (LodestoneTracker) ID() string {
	return "minecraft:lodestone_tracker"
}

// ReadFrom implements DataComponent.
func (l *LodestoneTracker) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{&l.GlobalPosition, &l.Tracked}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (l *LodestoneTracker) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{&l.GlobalPosition, &l.Tracked}.WriteTo(w)
}
