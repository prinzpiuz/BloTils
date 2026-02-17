#!/bin/sh
set -e

BLOTILS_DIR="${BLOTILS_DIR:-/blotils}"

echo "Starting BloTils..."

if [ "$(id -u)" = "0" ]; then
    echo "Running as root, adjusting permissions..."

    # Fix ownership of data directory
    if [ -d "$BLOTILS_DIR" ]; then
        chown blotils:blotils "$BLOTILS_DIR"
    fi

    # Fix ownership of database file
    if [ -f "$BLOTILS_DIR/BloTils.db" ]; then
        chown blotils:blotils "$BLOTILS_DIR/BloTils.db"
        chmod 660 "$BLOTILS_DIR/BloTils.db"
    fi

    # Fix ownership of config file
    if [ -f "$BLOTILS_DIR/config.json" ]; then
        chown blotils:blotils "$BLOTILS_DIR/config.json"
        chmod 640 "$BLOTILS_DIR/config.json"
    fi

    echo "Permissions adjusted. Switching to blotils user..."
    exec su-exec blotils "$@"
else
    exec "$@"
fi
