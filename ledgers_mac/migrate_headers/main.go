// One-off: adds a header to every journal in a directory that doesn't have one,
// moving old-style "tags: ..." lines into it. Last-modified times are kept so
// `ledgers -m` sorting is unaffected. Safe to rerun.
//
//	go run ./ledgers_mac/migrate_headers /path/to/journals
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/donchuru/ledgers/ledgers_mac/header"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: migrate_headers <journals directory>")
		os.Exit(1)
	}
	dir := os.Args[1]

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if err := migrate(filepath.Join(dir, e.Name())); err != nil {
			fmt.Printf("%s: %v\n", e.Name(), err)
			os.Exit(1)
		}
	}
}

func migrate(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if _, _, ok := header.Parse(string(raw)); ok {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	title := strings.TrimSuffix(filepath.Base(path), ".txt")
	h, body := header.FromLegacy(title, string(raw), info.ModTime())
	if err := os.WriteFile(path, []byte(h.String()+body), info.Mode()); err != nil {
		return err
	}
	fmt.Println("migrated", filepath.Base(path))
	return os.Chtimes(path, info.ModTime(), info.ModTime())
}
