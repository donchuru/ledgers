package main

import (
	"slices"
	"testing"
)

func TestWebFlagAnywhere(t *testing.T) {
	tests := []struct {
		args          []string
		wantInBrowser bool
		wantFilename  string
		wantTags      []string
	}{
		{[]string{"ledger", "new Doc"}, false, "new Doc.txt", nil},
		{[]string{"ledger", "new Doc", "-web"}, true, "new Doc.txt", nil},
		{[]string{"ledger", "new Doc", "-web", "-t", "life", "dreams"}, true, "new Doc.txt", []string{"life", "dreams"}},
		{[]string{"ledger", "new Doc", "-t", "life", "dreams", "-web"}, true, "new Doc.txt", []string{"life", "dreams"}},
	}
	for _, tt := range tests {
		args, inBrowser := removeFlag(tt.args, "-web")
		entry := createLedgerEntry(args)

		if inBrowser != tt.wantInBrowser || entry.filename != tt.wantFilename || !slices.Equal(entry.tags, tt.wantTags) {
			t.Errorf("%q: got (%v, %q, %q), want (%v, %q, %q)", tt.args,
				inBrowser, entry.filename, entry.tags, tt.wantInBrowser, tt.wantFilename, tt.wantTags)
		}
	}
}

func TestWebFlagAloneUsesTodaysDate(t *testing.T) {
	args, inBrowser := removeFlag([]string{"ledger", "-web"}, "-web")

	if entry := createLedgerEntry(args); !inBrowser || entry.filename == "-web.txt" {
		t.Errorf("got (%v, %q), want a dated journal opened in the browser", inBrowser, entry.filename)
	}
}
