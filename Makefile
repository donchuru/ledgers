.PHONY: build

build:
	cd ledgers_mac/config && go build -o ledgers-config
	cd ledgers_mac/ledger && go build -o ledger
	cd ledgers_mac/ledgers && go build -o ledgers
	./copy_binaries.sh 