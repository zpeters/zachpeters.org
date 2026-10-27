package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadProjectPages(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("old.md", `---
name: "Old Thing"
years: "2015"
summary: "An old thing."
---

Full text.
`)
	write("new.md", `---
name: "New Thing"
years: "2017-2019"
summary: "A new thing."
---

Newer text with a note[^1].

[^1]: The note.
`)
	write("plain.md", `---
name: "Plain"
years: "2020"
---

No summary.
`)

	got, err := LoadProjectPages(filepath.Join(dir, "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d projects, want 3: %+v", len(got), got)
	}

	// Newest first
	wantIDs := []string{"plain", "new", "old"}
	for i, id := range wantIDs {
		if got[i].ID != id {
			t.Errorf("project %d = %q, want %q", i, got[i].ID, id)
		}
	}

	byID := map[string]ProjectPage{}
	for _, p := range got {
		byID[p.ID] = p
	}
	if s := byID["new"].Summary; s != "A new thing." {
		t.Errorf("summary = %q", s)
	}
	if byID["new"].Name != "New Thing" || byID["new"].Years != "2017-2019" {
		t.Errorf("name/years = %q/%q", byID["new"].Name, byID["new"].Years)
	}
	if byID["new"].Start != 2017 {
		t.Errorf("start = %d, want 2017", byID["new"].Start)
	}
	if body := string(byID["new"].Body); !strings.Contains(body, `<span class="sidenote">The note.</span>`) {
		t.Errorf("body should render footnotes as sidenotes:\n%s", body)
	}
	if !strings.Contains(string(byID["old"].Body), "Full text.") {
		t.Error("body should be the Markdown text after the front matter")
	}
	if strings.Contains(string(byID["old"].Body), "---") {
		t.Error("front matter should not appear in the body")
	}
	if byID["plain"].Summary != "" {
		t.Errorf("project without a summary has summary %q", byID["plain"].Summary)
	}
}

func TestParseProjectPageErrors(t *testing.T) {
	cases := map[string]string{
		"no front matter": "Just text.\n",
		"not closed":      "---\nname: \"X\"\nyears: \"2020\"\n",
		"missing years":   "---\nname: \"X\"\n---\n",
		"missing name":    "---\nyears: \"2020\"\n---\n",
		"bad yaml":        "---\nname: [unclosed\n---\n",
	}
	for name, src := range cases {
		if _, err := ParseProjectPage("x", []byte(src)); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}
