---
title: HTTPS
---

Restreamer currently doesn't support serving the player and its administration interface via HTTPS. You can circumvent this by e.g. putting a nginx
in front of Restreamer that has a valid certificate.

## nginx

This example for nginx assumes that the Restreamer is listening on `localhost:8080`.

```nginx
...
    server {
    	listen 443 ssl http2;
    	server_name ...;

    	[SSL configuration]

        location / {
		proxy_http_version 1.1;
		proxy_pass http://localhost:8080;
	}
    }
...
```

Check the [nginx documentation](https://nginx.org/en/docs/) for more details.