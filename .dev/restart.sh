#!/bin/bash

file_path="files/etc/empty.json"

if [ -s "$file_path" ]; then
    # If file is not empty, truncate it
    > "$file_path"
else
    # If file is empty, write a character to it
    echo "X" > "$file_path"
fi
