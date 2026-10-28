.PHONY: build

# The web editor pulls in Apple frameworks, and macOS won't launch binaries that
# link them without an LC_UUID, which Go 1.22's own linker doesn't write.
# The system linker does.
LINK := -ldflags=-linkmode=external

build:
	cd ledgers_mac/config && go build -o ledgers-config
	cd ledgers_mac/ledger && go build $(LINK) -o ledger
	cd ledgers_mac/ledgers && go build $(LINK) -o ledgers
	./copy_binaries.sh 