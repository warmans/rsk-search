FROM alpine:3.13

RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.34-r0/glibc-2.34-r0.apk && \
    apk add glibc-2.34-r0.apk

RUN mkdir -p /opt/rsk && chown -R nobody /opt/rsk

ARG USER=nobody
USER nobody

WORKDIR /opt/rsk

COPY --chown=nobody ./bin/rsk-search .

RUN chmod +x rsk-search
RUN mkdir x509

EXPOSE 8888

CMD ["/opt/rsk/rsk-search", "discord-bot"]
