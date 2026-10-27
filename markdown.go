package main

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// RenderMarkdown renders a blog post. Markdown footnotes ("text[^1]" and
// "[^1]: note") become Tufte CSS sidenotes.
func RenderMarkdown(src []byte) []byte {
	return renderMarkdown(src, "sn-")
}

// renderMarkdown is RenderMarkdown with the sidenote ids starting with
// idPrefix, for pages that show more than one rendered document.
func renderMarkdown(src []byte, idPrefix string) []byte {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Typographer),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(util.Prioritized(&sidenoteTransformer{}, 1000)),
		),
		// Raw HTML is passed through, because project pages use figures.
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			renderer.WithNodeRenderers(util.Prioritized(&sidenoteRenderer{idPrefix}, 0)),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil
	}
	return buf.Bytes()
}

// kindSidenote is a footnote reference that is drawn as a sidenote.
var kindSidenote = ast.NewNodeKind("Sidenote")

// sidenote holds the text of one footnote, placed where its reference was.
type sidenote struct {
	ast.BaseInline
	Index int
}

func (n *sidenote) Kind() ast.NodeKind { return kindSidenote }

func (n *sidenote) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Index": fmt.Sprint(n.Index)}, nil)
}

// sidenoteTransformer replaces each footnote reference with a sidenote that
// holds the footnote's inline text, then removes the footnote list.
type sidenoteTransformer struct{}

func (t *sidenoteTransformer) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	var list *east.FootnoteList
	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		if l, ok := c.(*east.FootnoteList); ok {
			list = l
		}
	}
	if list == nil {
		return
	}

	notes := map[int]*east.Footnote{}
	for c := list.FirstChild(); c != nil; c = c.NextSibling() {
		f := c.(*east.Footnote)
		notes[f.Index] = f
	}

	var refs []*east.FootnoteLink
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if l, ok := n.(*east.FootnoteLink); ok && entering {
			refs = append(refs, l)
		}
		return ast.WalkContinue, nil
	})

	for _, ref := range refs {
		note, ok := notes[ref.Index]
		if !ok {
			continue
		}
		s := &sidenote{Index: ref.Index}
		for p := note.FirstChild(); p != nil; p = p.NextSibling() {
			if s.HasChildren() {
				s.AppendChild(s, ast.NewString([]byte(" ")))
			}
			for c := p.FirstChild(); c != nil; {
				next := c.NextSibling()
				p.RemoveChild(p, c)
				if c.Kind() != east.KindFootnoteBacklink {
					s.AppendChild(s, c)
				}
				c = next
			}
		}
		parent := ref.Parent()
		parent.ReplaceChild(parent, ref, s)
	}

	doc.RemoveChild(doc, list)
}

// sidenoteRenderer writes a sidenote in the markup Tufte CSS expects.
type sidenoteRenderer struct {
	idPrefix string
}

func (r *sidenoteRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindSidenote, r.render)
}

func (r *sidenoteRenderer) render(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.WriteString("</span>")
		return ast.WalkContinue, nil
	}
	id := fmt.Sprintf("%s%d", r.idPrefix, n.(*sidenote).Index)
	_, _ = fmt.Fprintf(w, `<label for="%s" class="margin-toggle sidenote-number"></label>`+
		`<input type="checkbox" id="%s" class="margin-toggle"/><span class="sidenote">`, id, id)
	return ast.WalkContinue, nil
}
