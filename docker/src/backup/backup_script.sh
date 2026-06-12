#!/bin/bash

BACKUP_DIR="/mnt/backup"
DIRS_TO_BACKUP=("/mnt/postgres" "/mnt/prometheus" "/mnt/esdata")
MAX_BACKUPS=5
TEMP_TAR="/tmp/temp.tar"

mkdir -p "$BACKUP_DIR"

backup_and_prune() {
    local dir="$1"
    local base_name=$(basename "$dir")

    tar -cf "$BACKUP_DIR/${base_name}_$(date +%Y%m%d_%H%M%S).tar" -C "/mnt" "$base_name"

    local tar_files=($(ls -t "$BACKUP_DIR/${base_name}_"*.tar 2>/dev/null))

    if [ ${#tar_files[@]} -gt $MAX_BACKUPS ]; then
        rm -f "${tar_files[$((MAX_BACKUPS))]}"
        echo "Deleted oldest backup for ${base_name}. Keeping ${MAX_BACKUPS} newest ones."
    fi
}

for dir in "${DIRS_TO_BACKUP[@]}"; do
    if [ -d "$dir" ]; then
        backup_and_prune "$dir"
    else
        echo "Warning: Directory $dir does not exist. Skipping..."
    fi
done

echo "Backup process completed."