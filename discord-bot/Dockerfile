FROM debian:stable-slim

RUN apt update && apt install -y gcc libfreetype-dev ffmpeg
RUN addgroup nobody

RUN mkdir -p /opt/rsk && chown -R nobody /opt/rsk

ARG USER=nobody
USER nobody

WORKDIR /opt/rsk

COPY --chown=nobody ./rsk-search .

RUN chmod +x rsk-search
RUN mkdir x509

EXPOSE 8888

CMD ["/opt/rsk/rsk-search", "discord-bot"]
