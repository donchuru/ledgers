package header

import (
	"slices"
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	h := Header{Title: "on going back to work", Created: "2024-05-17", Tags: []string{"career", "dream interpretations"}}
	body := "\nFriday, May 17, 2024\n\nsome text\n"

	got, gotBody, ok := Parse(h.String() + body)

	if !ok || got.Title != h.Title || got.Created != h.Created || !slices.Equal(got.Tags, h.Tags) || gotBody != body {
		t.Errorf("Parse(String()) = %+v, %q, %v; want %+v, %q, true", got, gotBody, ok, h, body)
	}
}

func TestParseIsLenient(t *testing.T) {
	content := "———\n TAGS : a,  b ,\ntitle: x: y\nmood: happy\n-----\nbody"

	h, body, ok := Parse(content)

	if !ok || h.Title != "x: y" || !slices.Equal(h.Tags, []string{"a", "b"}) || body != "body" {
		t.Errorf("Parse = %+v, %q, %v", h, body, ok)
	}
}

func TestParseWithoutHeader(t *testing.T) {
	for _, content := range []string{"", "tags: a\nhello", "---\nTitle: never closed\n"} {
		if _, body, ok := Parse(content); ok || body != content {
			t.Errorf("Parse(%q) = %q, %v; want unchanged, false", content, body, ok)
		}
	}
}

func TestAddTags(t *testing.T) {
	h := Header{Tags: []string{"career"}}

	h.AddTags("career,", "life", "dream interpretations")

	if want := []string{"career", "life", "dream interpretations"}; !slices.Equal(h.Tags, want) {
		t.Errorf("Tags = %q, want %q", h.Tags, want)
	}
}

func TestFromLegacy(t *testing.T) {
	fallback := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name, title, content string
		want                 Header
		wantBody             string
	}{
		{
			name:     "tags and title lines",
			title:    "on going back to work",
			content:  "tags: career\non going back to work\nFriday, May 17, 2024\n\nhi\n\ntags: life, idk\nThursday, May 23, 2024\n",
			want:     Header{Title: "on going back to work", Created: "2024-05-17", Tags: []string{"career", "life", "idk"}},
			wantBody: "\nFriday, May 17, 2024\n\nhi\n\nThursday, May 23, 2024\n",
		},
		{
			name:     "title line with extension",
			title:    "test",
			content:  "tags: a\ntest.txt\nWednesday, May 15, 2024\n",
			want:     Header{Title: "test", Created: "2024-05-15", Tags: []string{"a"}},
			wantBody: "\nWednesday, May 15, 2024\n",
		},
		{
			name:     "no date stamp",
			title:    "today",
			content:  "working on whatever",
			want:     Header{Title: "today", Created: "2025-01-02"},
			wantBody: "\nworking on whatever",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, body := FromLegacy(tt.title, tt.content, fallback)

			if h.Title != tt.want.Title || h.Created != tt.want.Created || !slices.Equal(h.Tags, tt.want.Tags) || body != tt.wantBody {
				t.Errorf("FromLegacy = %+v, %q; want %+v, %q", h, body, tt.want, tt.wantBody)
			}
		})
	}
}
