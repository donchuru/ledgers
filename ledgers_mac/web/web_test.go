package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/donchuru/ledgers/ledgers_mac/header"
)

const token = "secret"

// newServer serves a journal with the given contents from a temp dir
func newServer(t *testing.T, contents string) (*httptest.Server, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "dreams.txt")
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer((&journal{path: path}).routes(token))
	t.Cleanup(srv.Close)
	return srv, path
}

var existing = header.Header{Title: "dreams", Created: "2024-05-17", Tags: []string{"life"}}.String() +
	"\nFriday, May 17, 2024\nflying again\n"

func TestSaveWritesTagsAndBody(t *testing.T) {
	srv, path := newServer(t, existing)
	body := "Friday, May 17, 2024\nflying again\n\nMonday, Oct 27, 2026\n<b>falling</b> this time\n"

	edit, _ := json.Marshal(map[string]string{"tags": "dream interpretations, life, life", "body": body})

	res, err := http.Post(srv.URL+"/"+token+"/save", "application/json", bytes.NewReader(edit))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusNoContent)
	}
	raw, _ := os.ReadFile(path)
	h, gotBody, ok := header.Parse(string(raw))
	if !ok {
		t.Fatalf("saved journal has no header:\n%s", raw)
	}
	want := header.Header{Title: "dreams", Created: "2024-05-17", Tags: []string{"dream interpretations", "life"}}
	if h.Title != want.Title || h.Created != want.Created || !slices.Equal(h.Tags, want.Tags) {
		t.Errorf("header = %+v, want %+v", h, want)
	}
	if gotBody != "\n"+body {
		t.Errorf("body = %q, want %q", gotBody, "\n"+body)
	}
}

func TestSaveNeedsToken(t *testing.T) {
	srv, path := newServer(t, existing)

	res, err := http.Post(srv.URL+"/wrong/save", "application/json",
		strings.NewReader(`{"tags": "", "body": "gone"}`))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNotFound)
	}
	if raw, _ := os.ReadFile(path); string(raw) != existing {
		t.Errorf("journal changed:\n%s", raw)
	}
}

func TestPageShowsJournal(t *testing.T) {
	srv, _ := newServer(t, header.Header{Title: "dreams", Created: "2024-05-17", Tags: []string{"life", "a&b"}}.String()+
		"\n\nstarts with a blank line\n<script>\n")

	res, err := http.Get(srv.URL + "/" + token + "/")
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(res.Body)
	res.Body.Close()
	html := string(raw)

	for _, want := range []string{
		"<h1>dreams</h1>",
		"2024-05-17",
		`value="life, a&amp;b"`,
		// the browser drops the newline after <textarea>, leaving the body's own blank line
		"<textarea id=\"body\">\n\nstarts with a blank line\n&lt;script&gt;\n</textarea>",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("page is missing %q", want)
		}
	}
}

func TestPageBuildsHeaderForLegacyJournal(t *testing.T) {
	srv, _ := newServer(t, "tags: career\ndreams\nFriday, May 17, 2024\nflying again\n")

	res, err := http.Get(srv.URL + "/" + token + "/")
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(res.Body)
	res.Body.Close()
	html := string(raw)

	for _, want := range []string{"2024-05-17", `value="career"`, "<textarea id=\"body\">\nFriday, May 17, 2024\nflying again\n</textarea>"} {
		if !strings.Contains(html, want) {
			t.Errorf("page is missing %q", want)
		}
	}
}
