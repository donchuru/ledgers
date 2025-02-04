.PHONY: build

build:
	cd ledgers_mac/ledger && go build
	cd ledgers_mac/ledgers && go build
	./copy_binaries.sh 