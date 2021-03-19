// Grafonet DOCS: https://github.com/grafana/grafonnet-lib/blob/master/grafonnet-7.0/DOCS.md
local grafana = import 'grafonnet-7.0/grafana.libsonnet';
local dashboard = grafana.dashboard;
local graph = grafana.panel.graph;
local gauge = grafana.panel.gauge;
local row = grafana.panel.row;
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
    row.new(
      title="Connections",
      collapsed=false
    )
    .setGridPos(
      w=24,
      h=1
    )
  )
  .addPanel(
    graph.new(
      datasource="Prometheus",
      title="Avg req per min"
    )
    .setGridPos(
      h=9,
      w=12,
      y=1
    )
    .addTarget(
      prometheus.new(
        expr="rate(nginx_http_requests_total[5m])/5",
        legendFormat="POD: {{pod_template_hash}}"
      )
    )
  )
  .addPanel(
    gauge.new(
      datasource="Prometheus",
      title="Active connections"
    )
    .setGridPos(
      h=9,
      w=12,
      x=12,
      y=1
    )
    .addTarget(
      prometheus.new(
        expr="sum(nginx_connections_active)"
      )
    )
  )
}