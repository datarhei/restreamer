---
title: HTTPS
---

For enabling HTTPS in Restreamer you have two options. First, you can provide your certificate directly to Restreamer or you can
proxy the requests to Restreamer to a HTTPS enabled webserver, e.g. nginx.

## BYOC ("Bring Your Own Certificate")

You will need a valid certificate for the IP or the FQDN, i.e. (sub-)domain, your Restreamer will be reachable at.

The certificate has to be in [PEM format](https://en.wikipedia.org/wiki/X.509#Certificate_filename_extensions). Regarding
intermediate certificates the [nginx documentation explains](http://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl_certificate):
If intermediate certificates should be specified in addition to a primary certificate, they should be specified in the same
file in the following order: the primary certificate comes first, then the intermediate certificates.

The key also has to be in PEM format.

Place the certificate and key file as `cert.pem` and `key.pem` in a subdirectory called `ssl` into the directory that you
mount to `/restreamer/db`, i.e. inside the Docker container Restreamer expects the certificate and key file as `/restreamer/db/ssl/cert.pem`
and `/restreamer/db/ssl/key.pem`.

In order enable HTTPS and make Restreamer use the certificate, set the environment variable `RS_HTTPS` to `true` and
expose the port `8181`. Restreamer will listen on port `8181` for incoming HTTPS requests. The normal HTTP port `8080` is still available
and can be exposed as well. 

The Docker command line would look like:

```bash
docker run ... -p 443:8181 ... -e RS_HTTPS=true ...
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
