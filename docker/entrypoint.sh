#!/bin/sh

# This script runs as root first (before USER blotils takes effect for the CMD).
# It ensures that the BloTils.db file in the mounted volume has the correct
# ownership and permissions for the 'blotils' user (UID 1001).

DB_PATH="/blotils/BloTils.db"
BLOTILS_UID=1001
BLOTILS_GID=1001 # Assuming group 'blotils' has GID 1001

echo "Ensuring BloTils.db permissions for UID ${BLOTILS_UID}..."

# Check if the database file exists (e.g., if it's already mounted).
# If it exists, ensure its ownership is correct.
if [ -f "$DB_PATH" ]; then
    echo "Database file $DB_PATH found. Setting ownership..."
    chown "${BLOTILS_UID}:${BLOTILS_GID}" "$DB_PATH"
    chmod 660 "$DB_PATH" # Or 640 depending on your exact needs (user/group r/w, others no access)
elif [ ! -f "$DB_PATH" ] && [ -d "/blotils" ]; then
    # If the file doesn't exist but the directory does, it might be a new volume.
    # Ensure the parent directory has correct permissions so the app can create the file.
    echo "Database file $DB_PATH not found. Ensuring parent directory permissions..."
    chown "${BLOTILS_UID}:${BLOTILS_GID}" "/blotils"
    chmod 770 "/blotils" # Allow blotils user/group to create files in /blotils
fi

echo "Permissions adjusted. Starting application..."

# Execute the original CMD as the 'blotils' user
exec "$@"
