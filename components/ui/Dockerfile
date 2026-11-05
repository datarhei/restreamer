ARG NODE_IMAGE=node:21-alpine3.20
ARG CADDY_IMAGE=caddy:2.8.4-alpine

FROM $NODE_IMAGE AS builder

ENV PUBLIC_URL="./"

COPY . /ui

WORKDIR /ui

RUN cd /ui && \
	yarn install && \
	yarn build

FROM $CADDY_IMAGE

COPY --from=builder /ui/build /ui/build
COPY --from=builder /ui/Caddyfile /ui/Caddyfile

WORKDIR /ui

EXPOSE 3000

CMD [ "caddy", "run", "--config", "/ui/Caddyfile" ]
