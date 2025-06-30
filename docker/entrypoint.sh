#!/bin/sh

# This script runs as root first (before USER blotils takes effect for the CMD).
# It ensures that the BloTils.db file in the mounted volume has the correct
# ownership and permissions for the 'blotils' user (UID 1001).

BLOTILS_DIR="/blotils"
BLOTILS_UID=1001
BLOTILS_GID=1001 # Assuming group 'blotils' has GID 1001

echo "Ensuring permissions for UID ${BLOTILS_UID} on ${BLOTILS_DIR}..."

# Check if the database file exists (e.g., if it's already mounted).
# If it exists, ensure its ownership is correct.
if [ -d "$BLOTILS_DIR" ]; then
    echo "Blotils directory $BLOTILS_DIR found. Setting ownership..."
    chown -R "${BLOTILS_UID}:${BLOTILS_GID}" "$BLOTILS_DIR"
    chmod 770 "/blotils"
    # chmod 660 "$DB_PATH" # Or 640 depending on your exact needs (user/group r/w, others no access)
fi

echo "Permissions adjusted. Starting application..."

# Execute the original CMD as the 'blotils' user
exec su-exec blotils "$@"
