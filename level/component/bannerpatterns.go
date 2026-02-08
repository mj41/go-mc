package component

import (
	"io"

	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*BannerPatterns)(nil)

// BannerPatternLayer represents a single banner pattern layer.
// Wire: {pattern:registryEntryHolder<BannerPattern>, colorId:VarInt}
type BannerPatternLayer struct {
	// registryEntryHolder: VarInt, if 0 read inline BannerPattern, else registry ref (value-1)
	PatternType pk.VarInt
	InlineData  BannerPattern // only if PatternType == 0
	ColorID     pk.VarInt
}

func (l *BannerPatternLayer) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = l.PatternType.ReadFrom(r)
	if err != nil {
		return
	}
	if l.PatternType == 0 {
		n2, err := l.InlineData.ReadFrom(r)
		n += n2
		if err != nil {
			return n, err
		}
	}
	n2, err := l.ColorID.ReadFrom(r)
	n += n2
	return
}

func (l BannerPatternLayer) WriteTo(w io.Writer) (n int64, err error) {
	n, err = l.PatternType.WriteTo(w)
	if err != nil {
		return
	}
	if l.PatternType == 0 {
		n2, err := l.InlineData.WriteTo(w)
		n += n2
		if err != nil {
			return n, err
		}
	}
	n2, err := l.ColorID.WriteTo(w)
	n += n2
	return
}

// BannerPatterns component (wire 63).
// Wire: {layers:Array<BannerPatternLayer>}
type BannerPatterns struct {
	Layers []BannerPatternLayer
}

// ID implements DataComponent.
func (BannerPatterns) ID() string {
	return "minecraft:banner_patterns"
}

// ReadFrom implements DataComponent.
func (b *BannerPatterns) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Array(&b.Layers).ReadFrom(r)
}

// WriteTo implements DataComponent.
func (b *BannerPatterns) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Array(&b.Layers).WriteTo(w)
}
