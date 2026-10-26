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
	write("project-old.html", `<section>
  <h3 id="old">Old Thing (2015)</h3>
  <p class="summary">An old thing.</p>
  <p>Full text.</p>
</section>`)
	write("project-new.html", `<section>
  <h3 id="new">New Thing (2017-2019)</h3>
  <p class="summary">
    A new thing.
  </p>
</section>`)
	write("project-nosummary.html", `<h3 id="plain">Plain (2020)</h3><p>No summary.</p>`)
	write("project-noheading.html", `<p class="summary">Ignored.</p>`)

	got, err := LoadProjectPages(filepath.Join(dir, "project-*.html"))
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
	if s := string(byID["new"].Summary); s != "A new thing." {
		t.Errorf("summary = %q, want trimmed text", s)
	}
	if byID["new"].Name != "New Thing" || byID["new"].Years != "2017-2019" {
		t.Errorf("name/years = %q/%q", byID["new"].Name, byID["new"].Years)
	}
	if !strings.Contains(string(byID["old"].Body), "Full text.") {
		t.Error("body should be the whole file")
	}
	if byID["plain"].Summary != "" {
		t.Errorf("project without a summary has summary %q", byID["plain"].Summary)
	}
}
