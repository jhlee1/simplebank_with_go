#!/bin/sh
# We are using sh instead of bash because the image is alpine and alpine does not have bash installed by default

# It exits immediately if a command exits with a non-zero status
set -e

echo "run db migration"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"

# Takes the arguments passed to the script
# From the Dockerfile,
# CMD ["/app/main"]
# ENTRYPOINT["/app/start.sh"]
# is the same as ENTRYPOINT["/app/start.sh", "/app/main"]
exec "$@"
