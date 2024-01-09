#!/bin/bash

# Start the salad-queue-client
/usr/local/bin/salad-queue-client &

# Start the rest service
uvicorn main:app --host '*' --port 80 &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?
