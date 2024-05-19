FROM debian:stable-slim

RUN apt update && apt install -y gcc libfreetype-dev ffmpeg

#RUN wget -q -O /tmp/ffmpeg.tar.xz https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
#    cd /tmp && \
#    mkdir ffmpeg-static && \
#    tar xvf ffmpeg.tar.xz -C /tmp/ffmpeg-static --strip-components=1 && \
#    cd ffmpeg-static && \
#    mv ffmpeg /usr/bin/ffmpeg && mv ffprobe /usr/bin/ffprobe

RUN mkdir -p /opt/rsk/persistent && mkdir -p /opt/rsk/script && mkdir -p /opt/rsk/imports && chown -R nobody /opt/rsk

WORKDIR /opt/rsk

COPY var var
COPY bin/rsk-search .

RUN chmod +x rsk-search

EXPOSE 8888

CMD ["/opt/rsk/rsk-search", "server"]
