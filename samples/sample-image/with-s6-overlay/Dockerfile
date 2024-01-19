# syntax=docker/dockerfile:1

FROM foobar

WORKDIR /app

ARG S6_OVERLAY_VERSION=3.1.6.2
ENV S6_KEEP_ENV=1

RUN apt-get update && \
    apt-get -y install xz-utils curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp
RUN \
  tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz && \
  rm -rf /tmp/s6-overlay-noarch.tar.xz
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-x86_64.tar.xz /tmp
RUN \
  tar -C / -Jxpf /tmp/s6-overlay-x86_64.tar.xz && \
  rm -rf /tmp/s6-overlay-x86_64.tar.xz
COPY --chmod=755 s6 /etc/

COPY --chmod=755 ./salad-job-queue-worker /usr/local/bin/

ENTRYPOINT ["/init"]
CMD ["/usr/local/bin/salad-job-queue-worker"]
