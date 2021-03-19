local grafana = import 'grafonnet-7.0/grafana.libsonnet';
local dashboard = grafana.dashboard;
local graph = grafana.panel.graph;
local prometheus = grafana.target.prometheus;

{
  "nginx.json": dashboard.new(
    title='Nginx',
    tags=['nginx'],
    refresh="30s"
  )
  .setTime(
    from="now-1h"
  )
  .addPanel(
    graph.new(
      datasource="Prometheus",
      title="Avg req per min"
    )
    .setGridPos(
      h=9,
      w=12,
    )
    .addTarget(
      prometheus.new(
        expr="rate(nginx_http_requests_total[5m])/5",
        legendFormat="POD: {{pod_template_hash}}"
      )
    )
  )
}