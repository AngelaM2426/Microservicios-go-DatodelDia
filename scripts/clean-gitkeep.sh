#!/bin/bash

echo "Cleaning up .gitkeep files..."

find . -name ".gitkeep" | while read -r file; do
  dir=$(dirname "$file")
  count=$(find "$dir" -maxdepth 1 -type f ! -name ".gitkeep" | wc -l)

  if [ "$count" -gt 0 ]; then
    echo "Removing $file (other files found in $dir)"
    rm "$file"
  fi
done

echo "Cleanup complete!"
