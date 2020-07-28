---
title: HTTPS
---

For enabling HTTPS in Restreamer you have three options:

- provide your certificate directly to Restreamer (BYOC),
- proxy the requests to Restreamer from a HTTPS enabled webserver, e.g. nginx,
- use the provided Docker Compose app that does that all for you

Check below for a description of all mentioned options.

## BYOC ("Bring Your Own Certificate")

You will need a valid certificate for the IP or the FQDN, i.e. (sub-)domain, your Restreamer will be reachable at.

The certificate has to be in [PEM format](https://en.wikipedia.org/wiki/X.509#Certificate_filename_extensions). Regarding
intermediate certificates the [nginx documentation explains](http://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl_certificate):
If intermediate certificates should be specified in addition to a primary certificate, they should be specified in the same
file in the following order: the primary certificate comes first, then the intermediate certificates.

The key also has to be in PEM format.

Place the certificate and key file as `cert.pem` and `key.pem` in a subdirectory called `ssl` into the directory that you
mount to `/restreamer/db`. E.g. if you mount `/mnt/restreamer/db` to  `/restreamer/db`, then place the PEM files into `/mnt/restreamer/db/ssl`.
Inside the Docker container, Restreamer expects the certificate and key file as `/restreamer/db/ssl/cert.pem`
and `/restreamer/db/ssl/key.pem`.

In order enable HTTPS and make Restreamer use the certificate, set the environment variable `RS_HTTPS` to `true` and
expose the port `8181`. Restreamer will listen on port `8181` for incoming HTTPS requests. The normal HTTP port `8080` is still available
and can be exposed as well. 

The Docker command line would look like:

```bash
docker run ... -p 443:8181 -e RS_HTTPS=true -v /mnt/restreamer/db:/restreamer/db ...
```

Of course, you can map the HTTPS port to any port you prefer.

## Proxy

If you already have a webserver running (e.g. nginx) that serves sites with HTTPS, you can proxy the requests to your Restreamer.

### nginx

This example for nginx assumes that the Restreamer is listening on `localhost:8080`.

```nginx
...
    server {
    	listen 443 ssl http2;
    	server_name ...;

    	[SSL configuration]

        location /restreamer/ {
    		proxy_http_version 1.1;
    		proxy_pass http://localhost:8080;
    	}
    }
...
```

Check the [nginx documentation](https://nginx.org/en/docs/) for more details.

## Docker Compose

Docker Compose combines different docker containers to an app. The provided `docker-compose.yml` in the root of the repository
defines such an app that makes the Restreamer available via HTTPS with almost no configuration.

Make a copy of the `docker-compose.yml`. Now you need to adjust it to your needs. In the `restreamer:`
section edit the [environment variables](references-environment-var.html) for Restreamer (at least change the default username and password).
You might also want to change the directory that is mounted into the Restreamer container to persist your settings in the `volumes:` section.

In the `https-portal` section adjust the `DOMAINS` environment variables to the (sub-)domain you want the Restreamer running on. You have to own
the (sub-)domain and it has to point to the IP where Restreamer will run on. Set the value for `STAGE` to `production`. This will tell `https-portal`
to pull a valid certificate from Let's Encrypt. If you only want to try it out without pulling a real certificate from Let's Encrypt, set the value
for `STAGE` to `local` and `https-portal` will generate a self-signed certificate.

In the console, change to the directory where you stored the `docker-compose.yml` and execute

```shell
docker-compose up
```

or

```shell
docker-compose -f /path/to/docker-compose.yml up
```

In the default configuration a self-signed certificate for `localhost` will be generated. After starting the app you can access the Restreamer by
entering `https://localhost/` into your browser.
