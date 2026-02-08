package component

import "io"

var _ DataComponent = (*Glider)(nil)

type Glider struct{}

// ID implements DataComponent.
func (Glider) ID() string {
	return "minecraft:glider"
}

// ReadFrom implements DataComponent.
func (g *Glider) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

// WriteTo implements DataComponent.
func (g *Glider) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}
