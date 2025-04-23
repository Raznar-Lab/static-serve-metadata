#!/bin/ash
cd /app
mkdir -p /app/var
if [ -f "/app/var/config.yml" ]; then
  echo "Config file already exists, skipping copy."
else
  echo "Copying config file..."
  cp /app/config.example.yml /app/var/config.yml
fi

cp /app/var/config.yml /app/config.yml
/usr/local/bin/app-binary start