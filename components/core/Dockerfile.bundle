ARG CORE_IMAGE=datarhei/base:alpine-core-latest

ARG FFMPEG_IMAGE=datarhei/base:alpine-ffmpeg-latest

FROM $CORE_IMAGE as core

FROM $FFMPEG_IMAGE

COPY --from=core /core /core

RUN ffmpeg -buildconf

ENV CORE_CONFIGFILE=/core/config/config.json
ENV CORE_STORAGE_DISK_DIR=/core/data
ENV CORE_DB_DIR=/core/config

EXPOSE 8080/tcp
EXPOSE 8181/tcp
EXPOSE 1935/tcp
EXPOSE 1936/tcp
EXPOSE 6000/udp

VOLUME ["/core/data", "/core/config"]
ENTRYPOINT ["/core/bin/run.sh"]
WORKDIR /core
