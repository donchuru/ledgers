# Ledgers
Ledgers is a journaling CLI tool.
This is how I manage my journals on my Windows machine.

P.S. "ledgers" === "journals". I may use them interchangeably

## How to Use
### Set up config details
Run `ledgers-config` to add config details like your name and the path to the directory where you'd like to keep your journals.
\<Insert demo pic>

### Commands
Running `ledgers` shows you a list of your journal so far with details like the date it was last modified and the tags for each journal.
\<Insert demo pic>

`ledgers -m` shows a list of your journals (and corresponding details), sorted by date of last modification.
\<Insert demo pic>

`ledger newDoc` creates a journal titled "newDoc"

`ledger newDoc -t random` creates a journal titled "newDoc" with the "random" tag.
As of now, tags can only be one word. To add multiple word tags, include a space before each tag e.g. `ledger newDoc -t random insight`
A work around for multi-word tags is to underscore them e.g. `ledger newDoc -t dream_interpretations day_after_insomnia`
