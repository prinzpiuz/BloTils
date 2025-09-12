#!/bin/sh
BLOTILS_DIR="/blotils"
BLOTILS_UID=1001
BLOTILS_GID=1001

echo "Ensuring ownership and permissions for $BLOTILS_DIR..."

# Set correct ownership and permissions for the directory
chown "${BLOTILS_UID}:${BLOTILS_GID}" "$BLOTILS_DIR"
chmod 770 "$BLOTILS_DIR"

# (Optionally, set file permissions as before)
if [ -f "$BLOTILS_DIR/BloTils.db" ]; then
    chown "${BLOTILS_UID}:${BLOTILS_GID}" "$BLOTILS_DIR/BloTils.db"
    chmod 660 "$BLOTILS_DIR/BloTils.db"
fi

if [ -f "$BLOTILS_DIR/config.json" ]; then
    chown "${BLOTILS_UID}:${BLOTILS_GID}" "$BLOTILS_DIR/config.json"
    chmod 640 "$BLOTILS_DIR/config.json"
fi

echo "Permissions adjusted. Starting application..."

# Now switch to the proper user and run the application
exec su-exec blotils "$@"
