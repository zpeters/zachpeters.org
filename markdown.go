package main

import (
	"fmt"
	"io"

	"github.com/russross/blackfriday/v2"
)

// RenderMarkdown renders a blog post. Markdown footnotes ("text[^1]" and
// "[^1]: note") become Tufte CSS sidenotes.
func RenderMarkdown(src []byte) []byte {
	r := &sidenoteRenderer{blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags,
	})}
	return blackfriday.Run(src,
		blackfriday.WithRenderer(r),
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.Footnotes))
}

// sidenoteRenderer writes each footnote inline at its reference, in the
// sidenote markup Tufte CSS expects, and drops the footnote list.
type sidenoteRenderer struct {
	*blackfriday.HTMLRenderer
}

func (r *sidenoteRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch {
	case node.Type == blackfriday.List && node.IsFootnotesList:
		return blackfriday.SkipChildren
	case node.Type == blackfriday.Link && node.NoteID != 0:
		if entering {
			id := fmt.Sprintf("sn-%d", node.NoteID)
			_, _ = fmt.Fprintf(w, `<label for="%s" class="margin-toggle sidenote-number"></label>`+
				`<input type="checkbox" id="%s" class="margin-toggle"/><span class="sidenote">`, id, id)
			r.renderNote(w, node.Footnote)
			_, _ = io.WriteString(w, "</span>")
		}
		return blackfriday.SkipChildren
	}
	return r.HTMLRenderer.RenderNode(w, node, entering)
}

// renderNote writes a footnote's content without its <p> tags, because a
// sidenote is a <span> inside the paragraph.
func (r *sidenoteRenderer) renderNote(w io.Writer, item *blackfriday.Node) {
	first := true
	item.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		switch {
		case n == item:
			return blackfriday.GoToNext
		case n.Type == blackfriday.Paragraph:
			if entering && !first {
				_, _ = io.WriteString(w, " ")
			}
			first = false
			return blackfriday.GoToNext
		}
		return r.RenderNode(w, n, entering)
	})
}
