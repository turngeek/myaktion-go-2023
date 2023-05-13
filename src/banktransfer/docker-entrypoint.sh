#!/bin/sh

set -e

# Wait for kafka
if [ -n "$KAFKA_CONNECT" ]; then
    /go/src/app/wait-for-it.sh "$KAFKA_CONNECT" -t 20
fi

# Run the main container command.
exec "$@"
