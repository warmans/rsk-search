FROM alpine:3.13

RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.34-r0/glibc-2.34-r0.apk && \
    apk add glibc-2.34-r0.apk

RUN wget -q -O /tmp/ffmpeg.tar.xz https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    cd /tmp && \
    tar xvf ffmpeg.tar.xz && \
    cd ffmpeg-7.0-amd64-static && \
    mv ffmpeg /usr/bin/ffmpeg && mv ffprobe /usr/bin/ffprobe \

RUN mkdir -p /opt/rsk/persistent && mkdir -p /opt/rsk/script && mkdir -p /opt/rsk/imports && chown -R nobody /opt/rsk

ARG USER=nobody
USER nobody

WORKDIR /opt/rsk

COPY --chown=nobody var var
COPY --chown=nobody bin/rsk-search .

RUN chmod +x rsk-search

EXPOSE 8888

CMD ["/opt/rsk/rsk-search", "server"]
