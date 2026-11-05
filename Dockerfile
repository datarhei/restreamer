# syntax=docker/dockerfile:1.7

ARG NODE_IMAGE=node:24-alpine3.22
ARG GOLANG_IMAGE=golang:1.25-alpine3.22
ARG FFMPEG_IMAGE=ghcr.io/kaibsora/restreamer-ffmpeg:hardware-latest
ARG YARN_VERSION=1.22.22

FROM --platform=$BUILDPLATFORM ${NODE_IMAGE} AS ui-builder
ARG YARN_VERSION
ENV PUBLIC_URL=./
WORKDIR /src/ui
COPY components/ui/package.json components/ui/yarn.lock components/ui/.yarnrc.yml ./
RUN corepack enable && \
    corepack prepare yarn@${YARN_VERSION} --activate && \
    yarn install --frozen-lockfile
COPY components/ui/ ./
RUN yarn test-ci && yarn build

FROM --platform=$BUILDPLATFORM ${GOLANG_IMAGE} AS core-builder
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ENV GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT}
WORKDIR /src/core
RUN apk add --no-cache git make
COPY components/core/ ./
RUN go test ./... && \
    go vet ./... && \
    make release && \
    make import && \
    make ffmigrate

FROM ${FFMPEG_IMAGE}

COPY --from=core-builder /src/core/core /core/bin/core
COPY --from=core-builder /src/core/import /core/bin/import
COPY --from=core-builder /src/core/ffmigrate /core/bin/ffmigrate
COPY --from=core-builder /src/core/mime.types /core/mime.types
COPY --from=ui-builder /src/ui/build /core/ui
COPY run.sh /core/bin/run.sh
COPY ui-root /core/ui-root
COPY scripts/hardware-report.sh /usr/local/bin/restreamer-hardware-report

RUN mkdir -p /core/config /core/data && \
    chmod +x /core/bin/run.sh /usr/local/bin/restreamer-hardware-report && \
    ffmpeg -hide_banner -buildconf

ENV CORE_CONFIGFILE=/core/config/config.json \
    CORE_DB_DIR=/core/config \
    CORE_ROUTER_UI_PATH=/core/ui \
    CORE_STORAGE_DISK_DIR=/core/data

EXPOSE 8080/tcp 8181/tcp 1935/tcp 1936/tcp 6000/udp
VOLUME ["/core/data", "/core/config"]
ENTRYPOINT ["/core/bin/run.sh"]
WORKDIR /core
