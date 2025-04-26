# Status-page
Monitoring your server using prometheus exporter metrics with a single page html.

## Usage
```
./status-page --config.file=server.yaml
```

## Concept

![concept](/image/concept.png "")

<h3 id="Config-yaml">Config.yaml</h3>
<table>
<thead>
<tr>
<th>Name</th>
<th>Description</th>
<th>Value Example</th>
<th>Required</th>
</tr>
</thead>
<tbody><tr>
<td>port</td>
<td>Metric port address</td>
<td>9126</td>
<td>Yes</td>
</tr>
<tr>
<td>timeout</td>
<td>Maximum duration for reading the entire request, including the body. A zero or negative value means there will be no timeout.</td>
<td>5</td>
<td>Yes</td>
</tr>
<tr>
<td>scrape_interval</td>
<td>collect metrics for every *insert your value* seconds.</td>
<td>15</td>
<td>Yes</td>
</tr>
<tr>
<td>buffer_size</td>
<td>size of the write/read buffer used when writing to the transport (Go-httptraffic to client). Values are buffer_size * 1024.</td>
<td>1024</td>
<td>Yes</td>
</tr>
<tr>
<td>utc</td>
<td>show latest time for collect metrics.</td>
<td>-7</td>
<td>Yes</td>
</tr>
<tr>
<td>certfile</td>
<td>Certificate file for TLS connection</td>
<td>server.crt</td>
<td>No</td>
</tr>
<tr>
<td>keyfile</td>
<td>Key file for TLS connection</td>
<td>server.key</td>
<td>No</td>
</tr>
</tbody></table>

<h3 id="Config">Config</h3>
<table>
<thead>
<tr>
<th>Name</th>
<th>Description</th>
<th>Value Example</th>
<th>Required</th>
</tr>
</thead>
<tbody><tr>
<td>target</td>
<td>exporter url to collect metrics</td>
<td>http://localhost:9126/metrics</td>
<td>Yes</td>
</tr>
<tr>
<td>server_name</td>
<td>alias name of server for collect metrics</td>
<td>Server Node 1</td>
<td>Yes</td>
</tr>
<tr>
<td>tls_skip_verify</td>
<td>allow to collect metric without verify tls</td>
<td>true</td>
<td>No</td>
</tr>
<tr>
<td>username</td>
<td>Username for basic auth</td>
<td>root</td>
<td>No</td>
</tr>
<tr>
<td>password</td>
<td>Password for basic auth</td>
<td>password</td>
<td>No</td>
</tr>
<tr>
<td>bearer_token</td>
<td>token for bearer auth</td>
<td>eyJhbGc..............G4gRG9</td>
<td>No</td>
</tr>
</tbody></table>



<h3 id="List-metrics">List Metric</h3>
<table>
<thead>
<tr>
<th>Name</th>
<th>Description</th>
<th>Value Example</th>
<th>Required</th>
</tr>
</thead>
<tbody><tr>
<td>metric_name</td>
<td>prometheus metric name</td>
<td>windows_cpu_cstate_seconds_total{core="0,0",state="c1"}</td>
<td>Yes</td>
</tr>
<tr>
<td>name</td>
<td>alias of metric name</td>
<td>Windows CPU Usage Total</td>
<td>Yes</td>
</tr>
<tr>
<td>size</td>
<td>size of metric value. example %, bytes, 'C, 'F</td>
<td>Bytes</td>
<td>No</td>
</tr>
<tr>
<td>max</td>
<td>maximum value that will be the indicator for the metric.</td>
<td>1000</td>
<td>Yes</td>
</tr>
</tbody></table>

## Example Config file
```
port: 8082
timeout: 30
scrape_interval: 1
buffer_size: 1024
utc: -7
config:
  - target: "http://localhost:9126/metrics"
    server_name: "Server Node 1"
    list_metrics:
    - metric_name: "cpu_usage"
      name: "CPU Usage"
      size: "%"
      max: 100
    - metric_name: "ram_usage"
      name: "RAM Usage"
      size: "%"
      max: 100
    - metric_name: "cpu_temp"
      name: "CPU Temp"
      size: "'C"
      max: 100
    - metric_name: "disk_usage"
      name: "Disk Usage"
      size: "%"
      max: 100
```

## Screenshot

![page](/image/screenshot.png "")