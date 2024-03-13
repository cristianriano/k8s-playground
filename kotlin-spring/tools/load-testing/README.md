# Load Testing

It uses [K6](https://k6.io/) to run a simple load test.

```shell
docker run --rm --name load-test grafana/k6 run script.js --vus 5 --duration 30s
```