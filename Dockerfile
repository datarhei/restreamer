FROM datarhei/nginx-rtmp:1.13.4-dev as builder

FROM node:8.1-alpine

MAINTAINER datarhei <info@datarhei.org>

ENV RS_USERNAME=admin \
    RS_PASSWORD=datarhei \
    SRC=/usr/local \
    LD_LIBRARY_PATH=/usr/local/lib \
    PKG_CONFIG_PATH=/usr/local/lib/pkgconfig

COPY --from=builder /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=builder /usr/local/bin/ffprobe /usr/local/bin/ffprobe
COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /usr/local/nginx /usr/local/nginx

COPY . /restreamer
WORKDIR /restreamer

RUN apk add --no-cache --update git libssl1.0 pcre && \
    ffmpeg -buildconf && \
    npm install -g bower grunt grunt-cli nodemon eslint && \
    npm install && \
    grunt build && \
    bower cache clean --allow-root && \
    npm uninstall -g bower grunt grunt-cli nodemon eslint && \
    npm prune --production && \
    npm cache clean --force && \
    apk del --no-cache --purge git && \
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

EXPOSE 8080

VOLUME ["/restreamer/db"]

CMD ["./run.sh"]
