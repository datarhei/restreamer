ARG CADDY_IMAGE=caddy:2.7.5-alpine
FROM $CADDY_IMAGE

COPY build /ui/build
COPY Caddyfile /ui/Caddyfile

ENV PUBLIC_URL="./"

WORKDIR /ui

EXPOSE 3000

CMD [ "caddy", "run", "--config", "/ui/Caddyfile" ]