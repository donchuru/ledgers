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

`ledger newDoc` creates a journal titled "newDoc"

`ledger newDoc -t random` creates a journal titled "newDoc" with the "random" tag.
As of now, tags can only be one word. To add multiple word tags, include a space before each tag e.g. `ledger newDoc -t random insight`
A work around for multi-word tags is to underscore them e.g. `ledger newDoc -t dream_interpretations day_after_insomnia`


# When Developing
After developing, run `make build` from the root folder to replace binaries in `mac_dist`.
This will automatically update the binaries running on my machine since my PATH points to a symlink of this directory.  

Test as needed.
