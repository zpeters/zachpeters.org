package main

import "testing"

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Coding for One":   "coding-for-one",
		"Recommended Apps": "recommended-apps",
		"Tuple spaces (or, good ideas don't always win)": "tuple-spaces-or-good-ideas-dont-always-win",
		"Wu wei & more":             "wu-wei-more",
		"  Leading and trailing!  ": "leading-and-trailing",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseFilename(t *testing.T) {
	date, err := parseFilename("2026-10-18T21:30:00-0500_Tools_and_Notes.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := date.Format("2006-01-02"); got != "2026-10-18" {
		t.Errorf("date = %s", got)
	}
	if _, err := parseFilename("no-date-here.md"); err == nil {
		t.Error("want error for a filename without a date")
	}
}

func TestParsePost(t *testing.T) {
	title, body, err := parsePost([]byte("---\ntitle: Tools_and - Notes\n---\n\nHello.\n"))
	if err != nil {
		t.Fatal(err)
	}
	if title != "Tools_and - Notes" || body != "Hello." {
		t.Errorf("title = %q, body = %q", title, body)
	}
	for name, src := range map[string]string{
		"no front matter": "Hello.",
		"not closed":      "---\ntitle: X\nHello.",
		"no title":        "---\nauthor: Me\n---\nHello.",
	} {
		if _, _, err := parsePost([]byte(src)); err == nil {
			t.Errorf("%s: want error", name)
		}
	}
}
