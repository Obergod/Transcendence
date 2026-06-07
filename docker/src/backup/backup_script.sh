#!/bin/bash

# Configuration
BACKUP_DIR="/mnt/backup"
DIRS_TO_BACKUP=("/mnt/postgres" "/mnt/prometheus" "/mnt/esdata")
MAX_BACKUPS=5
TEMP_TAR="/tmp/temp.tar"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Function to create tar backup and manage retention
backup_and_prune() {
    local dir="$1"
    local base_name=$(basename "$dir")

    # Create the tar archive
    tar -cf "$BACKUP_DIR/${base_name}_$(date +%Y%m%d_%H%M%S).tar" -C "/mnt" "$base_name"

    # List all tar files for this directory, sorted by modification time (newest first)
    local tar_files=($(ls -t "$BACKUP_DIR/${base_name}_"*.tar 2>/dev/null))

    # If we have more than MAX_BACKUPS files, delete the oldest ones
    if [ ${#tar_files[@]} -gt $MAX_BACKUPS ]; then
        # Delete the oldest file (last in the sorted array)
        rm -f "${tar_files[$((MAX_BACKUPS))]}"
        echo "Deleted oldest backup for ${base_name}. Keeping ${MAX_BACKUPS} newest ones."
    fi
}

# Process each directory
for dir in "${DIRS_TO_BACKUP[@]}"; do
    if [ -d "$dir" ]; then
        backup_and_prune "$dir"
    else
        echo "Warning: Directory $dir does not exist. Skipping..."
    fi
done

echo "Backup process completed."