local grafana = import 'grafana-builder/grafana.libsonnet';

{
  "dashboards.json": grafana.dashboard("Nginx", datasource="Prometheus")
    .addRow(
      grafana.row('')
      .addPanel(
        grafana.queryPanel("rate(nginx_http_requests_total[5m])/5", '')
      )
    )
}
