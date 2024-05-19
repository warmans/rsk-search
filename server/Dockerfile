FROM alpine:3.19

RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.35-r1/glibc-2.35-r1.apk && \
    apk add glibc-2.35-r1.apk ffmpeg

RUN mkdir -p /opt/rsk/persistent && mkdir -p /opt/rsk/script && mkdir -p /opt/rsk/imports && chown -R nobody /opt/rsk

ARG USER=nobody
USER nobody

WORKDIR /opt/rsk

COPY --chown=nobody var var
COPY --chown=nobody bin/rsk-search .

RUN chmod +x rsk-search

EXPOSE 8888

CMD ["/opt/rsk/rsk-search", "server"]
