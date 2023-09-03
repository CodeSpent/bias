#!/bin/bash

echo "Generating random password..."
password=$(cat /run/secrets/mongodb_password | base64 | head -c 32)
echo "Generated password: $password"

export MONGO_INITDB_ROOT_PASSWORD=$password

exec "$@"
