#!/bin/bash

# Ensure the script exits on any error
set -e

# Define source and destination paths
SOURCE_DIR="./ledgers_mac"
DEST_DIR="./ledgers_mac/mac_dist"

# Copy config binary
if [ -f "$SOURCE_DIR/config/ledgers-config" ]; then
    echo "Copying ledgers-config binary..."
    cp "$SOURCE_DIR/config/ledgers-config" "$DEST_DIR/"
fi

# Copy ledger binary
if [ -f "$SOURCE_DIR/ledger/ledger" ]; then
    echo "Copying ledger binary..."
    cp "$SOURCE_DIR/ledger/ledger" "$DEST_DIR/"
fi

# Copy ledgers binary
if [ -f "$SOURCE_DIR/ledgers/ledgers" ]; then
    echo "Copying ledgers binary..."
    cp "$SOURCE_DIR/ledgers/ledgers" "$DEST_DIR/"
fi

echo "Binary copy complete!" 