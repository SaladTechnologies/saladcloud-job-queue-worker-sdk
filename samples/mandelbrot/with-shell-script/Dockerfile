# syntax=docker/dockerfile:1

FROM mandelbrot:latest

WORKDIR /app

COPY --chmod=755 ./run.sh ./
COPY --chmod=755 ./salad-http-job-queue-worker /usr/local/bin/

CMD ["/app/run.sh"]
