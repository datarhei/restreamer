---
title: Proxying
---

If you already have a webserver running (e.g. nginx) that serves your website, and you want Restreamer to be part of the
website, you can proxy the requests to a specific loaction to your Restreamer.

### nginx

This example for nginx assumes that the Restreamer is listening on `localhost:8080`.

```nginx
...
    server {
        listen 80;
        server_name ...;

        [your site configuration]

        location /restreamer/ {
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_pass http://localhost:8080;
        }
    }
...
```

If you have more than one Restreamer, you can expose them all under different locations, e.g.

```nginx
...
    server {
        listen 80;
        server_name ...;

        [your site configuration]

        location /camera1/ {
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_pass http://192.168.1.50:8080;
        }

        location /camera2/ {
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_pass http://192.168.1.51:8080;
        }

        location /camera3/ {
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_pass http://192.168.1.52:8080;
        }

        ...
    }
...
```

Based on the above example, with a bit more configuration it allows you to deliver the Restreamer player and admin interface with [HTTPS](guides-https.html) or to protect Restreamer
with a [login](https://nginx.org/en/docs/http/ngx_http_auth_basic_module.html), by [whitelisting IP addresses](https://nginx.org/en/docs/http/ngx_http_access_module.html),
any other authorization method, or a [combination thereof](https://nginx.org/en/docs/http/ngx_http_core_module.html#satisfy).

```nginx
...
    server {
        listen 443 ssl http2;
        server_name ...;

        [SSL configuration]

        [your site configuration]

        location /restreamer/ {
            satisfy any;

            allow 192.168.1.0/32;
            deny  all;

            auth_basic           "closed site";
            auth_basic_user_file conf/htpasswd;

            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_pass http://localhost:8080;
        }
    }
...
```

Check the [nginx documentation](https://nginx.org/en/docs/) for more details.
