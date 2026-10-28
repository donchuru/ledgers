// Package web opens a journal in a Chrome window for reading and writing.
//
// The page is served by a copy of the running program started in the
// background, so the terminal (or the ledgers TUI) is free while you write.
// That copy autosaves to the journal and exits once the window is closed.
package web

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/donchuru/ledgers/ledgers_mac/header"
)

// serveEnv holds the journal a background copy of the program should serve
const serveEnv = "LEDGERS_WEB_JOURNAL"

const (
	connectTimeout = 30 * time.Second // how long Chrome gets to load the page
	closeGrace     = 3 * time.Second  // how long the page can be gone before exiting, so reloads survive
	maxSaveSize    = 10 << 20
)

//go:embed page.html
var pageHTML string

var page = template.Must(template.New("page").Parse(pageHTML))

// Open opens the journal at path in a Chrome window
func Open(path string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe)
	cmd.Env = append(os.Environ(), serveEnv+"="+path)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // keep serving after the terminal closes
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait() // reap it if we outlive it, e.g. in the TUI
	return nil
}

// ServeIfRequested serves a journal and exits if this process was started by
// Open. Call it first thing in main.
func ServeIfRequested() {
	path := os.Getenv(serveEnv)
	if path == "" {
		return
	}
	if err := serve(path); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func serve(path string) error {
	token, err := newToken()
	if err != nil {
		return err
	}
	// only reachable from this machine, and only with the token
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	j := &journal{path: path, windows: make(chan int)}
	srv := &http.Server{Handler: j.routes(token)}
	go srv.Serve(ln)

	url := fmt.Sprintf("http://%s/%s/", ln.Addr(), token)
	if err := exec.Command("open", "-na", "Google Chrome", "--args", "--app="+url).Run(); err != nil {
		if err := exec.Command("open", url).Run(); err != nil {
			return err
		}
	}

	j.waitUntilClosed()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func newToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type journal struct {
	path    string
	mu      sync.Mutex // serializes saves
	windows chan int   // +1 when a page connects to /alive, -1 when it goes away
}

func (j *journal) routes(token string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /"+token+"/{$}", j.page)
	mux.HandleFunc("POST /"+token+"/save", j.save)
	mux.HandleFunc("GET /"+token+"/alive", j.alive)
	return mux
}

// waitUntilClosed returns once no page has been open for closeGrace, or if
// none opens within connectTimeout
func (j *journal) waitUntilClosed() {
	open, wait := 0, connectTimeout
	for {
		var timeout <-chan time.Time // nil, so never fires, while a page is open
		if open == 0 {
			timeout = time.After(wait)
		}
		select {
		case d := <-j.windows:
			open += d
			wait = closeGrace
		case <-timeout:
			return
		}
	}
}

// read returns the journal's header and body. Journals written before headers
// existed get a header built from the old format, which is written on save.
func (j *journal) read() (header.Header, string, error) {
	info, err := os.Stat(j.path)
	if err != nil {
		return header.Header{}, "", err
	}
	raw, err := os.ReadFile(j.path)
	if err != nil {
		return header.Header{}, "", err
	}
	h, body, ok := header.Parse(string(raw))
	if !ok {
		title := strings.TrimSuffix(filepath.Base(j.path), ".txt")
		h, body = header.FromLegacy(title, string(raw), info.ModTime())
	}
	// the blank line under the header isn't part of what you write
	return h, strings.TrimPrefix(body, "\n"), nil
}

func (j *journal) page(w http.ResponseWriter, r *http.Request) {
	h, body, err := j.read()
	if err != nil {
		http.Error(w, "couldn't read journal", http.StatusInternalServerError)
		return
	}
	page.Execute(w, struct {
		Title, Created, Tags, Body string
	}{h.Title, h.Created, strings.Join(h.Tags, ", "), body})
}

func (j *journal) save(w http.ResponseWriter, r *http.Request) {
	var edit struct {
		Tags string `json:"tags"`
		Body string `json:"body"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxSaveSize)).Decode(&edit); err != nil {
		http.Error(w, "invalid save", http.StatusBadRequest)
		return
	}

	j.mu.Lock()
	defer j.mu.Unlock()
	h, _, err := j.read()
	if err != nil {
		http.Error(w, "couldn't read journal", http.StatusInternalServerError)
		return
	}
	h.Tags = nil
	h.AddTags(edit.Tags)
	if err := os.WriteFile(j.path, []byte(h.String()+"\n"+edit.Body), 0644); err != nil {
		http.Error(w, "couldn't save journal", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// alive holds an event stream open for as long as the page is open
func (j *journal) alive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	j.windows <- 1
	<-r.Context().Done()
	j.windows <- -1
}
