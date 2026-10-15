# DrayTek Exporter

[![CI](https://github.com/SuperQ/draytek_exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/SuperQ/draytek_exporter/actions/workflows/ci.yml)

Prometheus exporter for DrayTek Vigor modems/routers. Collects DSL line status and performance metrics via the DrayTek web API.

## Supported Devices

* DrayTek Vigor 167 (v5 firmware)

## Installation

### Pre-built binaries

Download from the [releases page](https://github.com/SuperQ/draytek_exporter/releases).

### Docker

```sh
docker run -e DRAYTEK_PASSWORD=yourpassword quay.io/superq/draytek-exporter --target=192.168.1.1
```

Multi-arch images are available for `amd64`, `armv7`, and `arm64`.

### Building from source

Requires Go 1.25+.

```sh
make build
```

## Usage

```sh
export DRAYTEK_PASSWORD=yourpassword
./draytek_exporter --target=192.168.1.1
```

The exporter authenticates to the router's web interface, scrapes DSL status data, and exposes it as Prometheus metrics on `:9103/metrics`.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--target` | `192.168.1.1` | Router IP address or hostname |
| `--username` | `monitor` | Username for router authentication |
| `--password-env` | `DRAYTEK_PASSWORD` | Environment variable containing the password |
| `--web.listen-address` | `:9103` | Address to listen on for HTTP requests |
| `--web.telemetry-path` | `/metrics` | Path under which to expose metrics |

### Router Setup

Create a dedicated monitoring user on your DrayTek device (e.g. `monitor`) with read-only access. Avoid using the admin account.

## Metrics

All metrics use the `draytek_` prefix.

| Metric | Type | Description |
|--------|------|-------------|
| `draytek_up` | Gauge | Whether the DSL link is up (1) or down (0) |
| `draytek_info` | Gauge | Device info labels: `dsl_version`, `mode`, `profile` |
| `draytek_downstream_actual_bps` | Gauge | Current downstream sync speed (bits/sec) |
| `draytek_upstream_actual_bps` | Gauge | Current upstream sync speed (bits/sec) |
| `draytek_downstream_attainable_bps` | Gauge | Maximum attainable downstream speed (bits/sec) |
| `draytek_upstream_attainable_bps` | Gauge | Maximum attainable upstream speed (bits/sec) |
| `draytek_downstream_snr_margin_db` | Gauge | Downstream signal-to-noise ratio margin (dB) |
| `draytek_upstream_snr_margin_db` | Gauge | Upstream signal-to-noise ratio margin (dB) |
| `draytek_downstream_actual_psd_db` | Gauge | Downstream power spectrum density (dB) |
| `draytek_upstream_actual_psd_db` | Gauge | Upstream power spectrum density (dB) |
| `draytek_downstream_interleave_depth` | Gauge | Downstream interleaving depth |
| `draytek_upstream_interleave_depth` | Gauge | Upstream interleaving depth |
| `draytek_{near,far}_end_attenuation_db` | Gauge | Line attenuation (dB) |
| `draytek_{near,far}_end_bitswap_active` | Gauge | Whether bitswap is active |
| `draytek_{near,far}_end_retx_active` | Gauge | Whether retransmission (G.INP) is active |
| `draytek_{near,far}_end_crc_errors_total` | Counter | CRC error count |
| `draytek_{near,far}_end_errored_seconds_total` | Counter | Errored seconds |
| `draytek_{near,far}_end_severely_errored_seconds_total` | Counter | Severely errored seconds |
| `draytek_{near,far}_end_unavailable_seconds_total` | Counter | Unavailable seconds |
| `draytek_{near,far}_end_hec_errors_total` | Counter | Header error correction errors |
| `draytek_{near,far}_end_los_failures_total` | Counter | Loss of Signal failures |
| `draytek_{near,far}_end_lof_failures_total` | Counter | Loss of Frame failures |
| `draytek_{near,far}_end_lpr_failures_total` | Counter | Loss of Power failures |
| `draytek_{near,far}_end_lcd_failures_total` | Counter | Loss of Cell Delineation failures |
| `draytek_{near,far}_end_reed_solomon_forward_error_corrections_total` | Counter | Reed-Solomon FEC corrections |

## License

Apache License 2.0. See [LICENSE](LICENSE).
