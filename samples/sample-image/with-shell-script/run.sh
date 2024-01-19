#!/bin/bash

# Start the salad-job-queue-worker
/usr/local/bin/salad-job-queue-worker &

# Start the rest service
uvicorn main:app --host '*' --port 80 &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?
