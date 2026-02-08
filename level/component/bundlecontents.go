package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*BundleContents)(nil)

type BundleContents struct {
	Contents []SlotData
}

// ID implements DataComponent.
func (BundleContents) ID() string {
	return "minecraft:bundle_contents"
}

// ReadFrom implements DataComponent.
func (b *BundleContents) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&b.Contents).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (b *BundleContents) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&b.Contents).WriteTo(w)
}
