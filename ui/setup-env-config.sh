#!/bin/sh

output_file="$APP_DIR/env-config.js"

for var in $(env | grep '^REACT_APP_'); do
  key=$(echo "$var" | cut -d '=' -f 1)
  value=$(echo "$var" | cut -d '=' -f 2-)
  value=$(echo "$value" | sed 's/"/\\"/g')
  echo "window.REACT_APP_ENV.$key = \"$value\";" >> "$output_file"
done