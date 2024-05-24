#!/bin/bash
/usr/local/bin/salad-job-queue-worker &
uvicorn main:app --host '*' --port 80 &
wait -n
exit $?
