// Package header reads and writes the header block at the top of every journal:
//
//	----------------------------------------------------
//	Title:   on going back to work
//	Created: 2024-05-17
//	Tags:    career, life
//	----------------------------------------------------
package header

import (
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	DateLayout  = "2006-01-02"
	StampLayout = "Monday, Jan 2, 2006" // date stamp written above each entry

	border = "----------------------------------------------------"
)

type Header struct {
	Title   string
	Created string
	Tags    []string
}

func (h Header) String() string {
	return border + "\n" +
		field("Title", h.Title) +
		field("Created", h.Created) +
		field("Tags", strings.Join(h.Tags, ", ")) +
		border + "\n"
}

func field(key, value string) string {
	return strings.TrimRight(fmt.Sprintf("%-9s%s", key+":", value), " ") + "\n"
}

// AddTags adds tags not already present. Each argument may hold several
// comma-separated tags.
func (h *Header) AddTags(tags ...string) {
	for _, tag := range parseTags(strings.Join(tags, ",")) {
		if !slices.Contains(h.Tags, tag) {
			h.Tags = append(h.Tags, tag)
		}
	}
}

// Parse splits content into its header and the body that follows it. ok is
// false if content doesn't start with a header.
//
// Parsing is lenient since headers get edited by hand: a border is any line of
// 3+ dashes (TextEdit's smart dashes may turn them into em dashes), keys are
// case-insensitive, and unknown keys are ignored.
func Parse(content string) (h Header, body string, ok bool) {
	lines := strings.SplitAfter(content, "\n")
	if !isBorder(lines[0]) {
		return Header{}, content, false
	}
	for i, line := range lines[1:] {
		if isBorder(line) {
			return h, strings.Join(lines[i+2:], ""), true
		}
		key, value, _ := strings.Cut(line, ":")
		value = strings.TrimSpace(value)
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "title":
			h.Title = value
		case "created":
			h.Created = value
		case "tags":
			h.Tags = parseTags(value)
		}
	}
	return Header{}, content, false
}

// FromLegacy builds a header for a journal written before headers existed.
// Those started with an optional "tags: ..." line and an optional title line,
// and appends could add more "tags: ..." lines mid-file. Tags lines are merged
// into the header, the title line is dropped, and Created is taken from the
// first date stamp, falling back to fallback.
func FromLegacy(title, content string, fallback time.Time) (Header, string) {
	h := Header{Title: title}
	var kept []string
	for _, line := range strings.SplitAfter(content, "\n") {
		if tags, found := strings.CutPrefix(line, "tags: "); found {
			h.AddTags(tags)
			continue
		}
		kept = append(kept, line)
	}

	if len(kept) > 0 && strings.TrimSuffix(strings.TrimSpace(kept[0]), ".txt") == title {
		kept = kept[1:]
	}
	for _, line := range kept {
		if stamp, err := time.Parse(StampLayout, strings.TrimSpace(line)); err == nil {
			h.Created = stamp.Format(DateLayout)
			break
		}
	}
	if h.Created == "" {
		h.Created = fallback.Format(DateLayout)
	}
	return h, "\n" + strings.TrimLeft(strings.Join(kept, ""), "\n")
}

func isBorder(line string) bool {
	line = strings.TrimSpace(line)
	return utf8.RuneCountInString(line) >= 3 && strings.Trim(line, "-–—") == ""
}

func parseTags(s string) []string {
	var tags []string
	for _, tag := range strings.Split(s, ",") {
		if tag = strings.TrimSpace(tag); tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}
