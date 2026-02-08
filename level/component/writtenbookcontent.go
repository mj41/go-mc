package component

import (
	"io"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/nbt/dynbt"
	pk "github.com/Tnze/go-mc/net/packet"
)

var _ DataComponent = (*WrittenBookContent)(nil)

type WrittenBookContent struct {
	RawTitle      pk.String
	FilteredTitle pk.Option[pk.String, *pk.String]
	Author        pk.String
	Generation    pk.VarInt
	Pages         []WrittenPage
	Resolved      pk.Boolean
}

type WrittenPage struct {
	Content chat.Message // anonymousNbt
	// filteredContent: anonOptionalNbt — optional NBT, we store as dynbt.Value
	FilteredContent dynbt.Value
}

func (p *WrittenPage) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = p.Content.ReadFrom(r)
	if err != nil {
		return
	}
	// anonOptionalNbt — read NBT that may be TagEnd
	n2, err := pk.NBTField{V: &p.FilteredContent, AllowUnknownFields: true}.ReadFrom(r)
	return n + n2, err
}

func (p WrittenPage) WriteTo(w io.Writer) (n int64, err error) {
	n, err = p.Content.WriteTo(w)
	if err != nil {
		return
	}
	n2, err := pk.NBTField{V: &p.FilteredContent, AllowUnknownFields: true}.WriteTo(w)
	return n + n2, err
}

// ID implements DataComponent.
func (WrittenBookContent) ID() string {
	return "minecraft:written_book_content"
}

// ReadFrom implements DataComponent.
func (b *WrittenBookContent) ReadFrom(r io.Reader) (n int64, err error) {
	return pk.Tuple{
		&b.RawTitle,
		&b.FilteredTitle,
		&b.Author,
		&b.Generation,
		pk.Array(&b.Pages),
		&b.Resolved,
	}.ReadFrom(r)
}

// WriteTo implements DataComponent.
func (b *WrittenBookContent) WriteTo(w io.Writer) (n int64, err error) {
	return pk.Tuple{
		&b.RawTitle,
		&b.FilteredTitle,
		&b.Author,
		&b.Generation,
		pk.Array(&b.Pages),
		&b.Resolved,
	}.WriteTo(w)
}
