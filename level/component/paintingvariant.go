package component

import (
	"io"

	"github.com/Tnze/go-mc/chat"
	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*PaintingVariant)(nil)

// PaintingVariantData represents the inline data for a painting variant.
// Wire: {width:i32, height:i32, assetId:string, title:Option<anonymousNbt>, author:Option<anonymousNbt>}
type PaintingVariantData struct {
	Width   pk.Int
	Height  pk.Int
	AssetID pk.String
	Title   pk.Option[chat.Message, *chat.Message]
	Author  pk.Option[chat.Message, *chat.Message]
}

func (d *PaintingVariantData) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&d.Width, &d.Height, &d.AssetID, &d.Title, &d.Author}.ReadFrom(r)
}

func (d PaintingVariantData) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{&d.Width, &d.Height, &d.AssetID, &d.Title, &d.Author}.WriteTo(w)
}

// PaintingVariant component (wire 89).
// Wire: registryEntryHolder<EntityMetadataPaintingVariant>.
// VarInt type: 0=inline PaintingVariantData, >0=registry ref (value-1).
type PaintingVariant struct {
	Type       pk.VarInt
	InlineData PaintingVariantData // only if Type == 0
}

// ID implements DataComponent.
func (PaintingVariant) ID() string {
	return "minecraft:painting/variant"
}

// ReadFrom implements DataComponent.
func (p *PaintingVariant) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = p.Type.ReadFrom(r)
	if err != nil {
		return
	}
	if p.Type == 0 {
		n2, err := p.InlineData.ReadFrom(r)
		n += n2
		return n, err
	}
	return
}

// WriteTo implements DataComponent.
func (p *PaintingVariant) WriteTo(w io.Writer) (n int64, err error) {
	n, err = p.Type.WriteTo(w)
	if err != nil {
		return
	}
	if p.Type == 0 {
		n2, err := p.InlineData.WriteTo(w)
		n += n2
		return n, err
	}
	return
}
