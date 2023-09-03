#!/bin/bash

# Generate a random password and save it to secrets/postgres_password.txt
pwgen -s 32 1 > secrets/postgres_password.txt

# Start the PostgreSQL server
exec "$@"
