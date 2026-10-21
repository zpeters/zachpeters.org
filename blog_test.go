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
