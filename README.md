# Ledgers
Ledgers is a journaling TUI application for Mac and Windows.\
*N.B. I do not update the Windows distribution frequently so it might be a little outdated*

P.S. *"ledgers" === "journals"*. I may use them interchangeably

## How to Use (Mac)
1. Download `mac_dist` from the `ledgers_mac` folder.  
2. Add it to your PATH environment variable by adding `export PATH="/your/folder/path:$PATH"` to your `.zshrc` file.

### Set up config details
Run `ledgers-config` to add config details like your name and the path to the directory where you'd like to keep your journals.

### Commands
Running `ledgers` shows you a list of your journal so far with details like the date it was last modified and the tags for each journal.

![ScreenRecording2025-02-03at7 11 01PM-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/57010e18-83ee-4aa4-9939-1c2938339314)

`ledgers -m` shows a list of your journals (and corresponding details), sorted by date of last modification.

Start typing to search: the list narrows as you type to journals whose title or tags fuzzy-match your search, best match first (e.g. `gbw` finds "on going back to work").
Use ↑/↓ to move, `enter` to open the selected journal in TextEdit, `tab` to open it in Chrome (see below), `esc` to clear the search (or quit if it's already empty), and `ctrl+c` to quit.

`ledger newDoc` creates a journal titled "newDoc"

`ledger newDoc -t random` creates a journal titled "newDoc" with the "random" tag.
To add multiple tags, separate them with spaces e.g. `ledger newDoc -t random insight`. Quote multi-word tags e.g. `ledger newDoc -t "dream interpretations" insight`.
Running `ledger` on an existing journal adds a new dated entry, and any `-t` tags are added to its header.

### Writing in the browser
Add `-web` to write in a Chrome window instead of TextEdit, e.g. `ledger newDoc -web` or `ledger newDoc -t random -web`.
The page is plain and quiet, styled like [donchuru.github.io](https://donchuru.github.io/): the title, the date and tags, then the journal itself with your cursor at the end.

- Everything you type autosaves a second after you pause, and `⌘S` saves right away. The note in the bottom-right corner says `saved`, or `not saved — retrying` in red if something's wrong.
- Edit the tags under the title (comma-separated) and they're saved to the header.
- The journal is still a plain text file, so it opens in TextEdit and shows up in `ledgers` like any other.
- The terminal is free as soon as Chrome opens. A small local server runs in the background while the window is open, then exits a few seconds after you close it.

To open an existing journal in the browser without adding a new entry, select it in `ledgers` and press `tab`.

### Journal format
Every journal starts with a header:
```
----------------------------------------------------
Title:   newDoc
Created: 2024-09-09
Tags:    random, dream interpretations
----------------------------------------------------
```
You can edit the header by hand (e.g. to add tags). Keep tags comma-separated.

Journals created before headers existed can be converted once (back up your journals folder first):
```
go run ./ledgers_mac/migrate_headers /path/to/your/journals
```


# When Developing
After developing, run `make build` from the root folder to replace binaries in `mac_dist`.
It links `ledger` and `ledgers` with the system linker (needs Xcode's command line tools), because recent macOS won't launch them otherwise. For the same reason, run their tests with `go test -ldflags=-linkmode=external ./...`.
This will automatically update the binaries running on my machine since my PATH points to a symlink of this directory.  

Test as needed.
