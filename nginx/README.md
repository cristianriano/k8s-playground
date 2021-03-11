# Nginx Container

It adds a metrics page on `nginx_status` via the [stub_status module](https://nginx.org/en/docs/http/ngx_http_stub_status_module.html#stub_status) on top of [nginxdemo container](https://hub.docker.com/r/nginxdemos/hello) by overwriting the default config.

## Push image

On the root directory of the project

```bash
docker login
docker build --tag nginxdemo:latest --file nginx/Dockerfile nginx/
docker tag nginxdemo:latest cristianriano/nginxdemo:latest
docker push cristianriano/nginxdemo:latest
```
