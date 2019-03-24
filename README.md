# php-fpm-exporter for K8S environments (avoid sidecar containers)

## Supported protocols
- fastcgi
- http (with Nginx)

## Installing the Chart
```$sh
helm install --name php-fpm-exporter deploy/php-fpm-exporter
```

### Prometheus job example

```yaml
  - job_name: 'kubernetes-php-fpm-endpoints'
    scrape_interval: 10s
    metrics_path: /metrics
    kubernetes_sd_configs:
      - role: endpoints
    relabel_configs:
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_php_fpm]
        action: keep
        regex: true
      - source_labels: [__address__]
        regex: ([^:]+)(?::\d+)?
        replacement: tcp://${1}:9000/status
        target_label: __param_target
      - target_label: __address__
        replacement: "{{ PHPFPM_EXPORTER_ADDRESS }}"
      - source_labels: [__param_target]
        target_label: target
      - action: labelmap
        regex: __meta_kubernetes_service_label_(.+)
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: kubernetes_namespace
      - source_labels: [__meta_kubernetes_service_name]
        action: replace
        target_label: kubernetes_name

```
