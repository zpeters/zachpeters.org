package main

import (
	"strings"
	"testing"
)

func TestRenderMarkdownSidenotes(t *testing.T) {
	src := "Helix uses selection first[^1], then vim[^vim].\n\n" +
		"[^1]: See `hx` and *Migrating from Vim*.\n" +
		"[^vim]: The classic.\n"
	got := string(RenderMarkdown([]byte(src)))

	for _, want := range []string{
		`first<label for="sn-1" class="margin-toggle sidenote-number"></label>`,
		`<input type="checkbox" id="sn-1" class="margin-toggle"/>`,
		`<span class="sidenote">See <code>hx</code> and <em>Migrating from Vim</em>.</span>`,
		`id="sn-2"`,
		`<span class="sidenote">The classic.</span>.`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"footnotes", "<sup", "<hr", "<ol"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("output contains %q\n%s", unwanted, got)
		}
	}
}

func TestRenderMarkdownPlain(t *testing.T) {
	got := string(RenderMarkdown([]byte("Just *text*.\n")))
	if got != "<p>Just <em>text</em>.</p>\n" {
		t.Errorf("got %q", got)
	}
}
