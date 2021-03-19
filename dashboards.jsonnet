local grafana = import 'grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local singlestat = grafana.singlestat;
local prometheus = grafana.prometheus;

{
  "nginx.json": dashboard.new(
    'Nginx',
    schemaVersion=16,
    tags=['nginx'],
    editable=true
  )
  .addPanel(
    singlestat.new(
      'Avg req per minute',
      format='s',
      datasource='Prometheus',
      span=2,
      valueName='current',
    )
    .addTarget(
      prometheus.target(
        'rate(nginx_http_requests_total[5m])/5',
      )
    ), gridPos={
      x: 0,
      y: 0,
      w: 24,
      h: 3,
    },
  )
}