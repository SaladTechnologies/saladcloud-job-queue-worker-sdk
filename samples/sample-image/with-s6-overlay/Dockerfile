FROM foobar

WORKDIR /app

# install the s6-overlay
ARG S6_OVERLAY_VERSION=3.1.6.2

RUN apt-get update \
    && apt-get install -y xz-utils \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz

ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-x86_64.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-x86_64.tar.xz

# copy previously downloaded salad-queue-client
COPY salad-queue-client /usr/local/bin/

# copy the s6-overlay configuration
COPY --chmod=755 s6 /etc/

ENV S6_KEEP_ENV=1

ENTRYPOINT ["/init"]
CMD ["/usr/local/bin/salad-queue-client"]

