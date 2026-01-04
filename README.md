# 2b2t Heatmaps Service 🗺️

A high-performance Tile CDN service designed for 2b2t heatmap data. Built with Go, it serves map tiles from the file
system with aggressive in-memory caching and smart ETag validation to minimize latency and bandwidth.

## 🚀 Features

* **Fast I/O:** Optimized file handling with single-descriptor `Stat` and `Read` operations.
* **In-Memory Cache:** Powered by **Ristretto** (TinyLFU policy) to prevent RAM exhaustion on budget VPS.
* **Smart ETagging:** Weak ETag implementation (`W/"size-mtime"`) handles `304 Not Modified` responses before reading
  file data.
* **Zero-Config SSL:** Automated Let's Encrypt / ZeroSSL via **Caddy**.
* **Modern Compression:** Support for **Zstd** and **Gzip** via the edge proxy.
* **Observability:** Real-time metrics (Cache Hit/Miss, ETag stats) exported via StatsD.

---

## 🛠️ Tech Stack

* **Language:** Go 1.25+
* **Framework:** Gin Gonic
* **Cache:** Ristretto
* **Proxy:** Caddy 2
* **Metrics:** StatsD (VictoriaMetrics compatible), Vector (UDP proxy for VictoriaLogs)
* **Dev Ops:** Docker, Air (live-reload), Make

---

## 🚦 API Endpoints

| Endpoint               | Description               | Parameters                                                                               |
|------------------------|---------------------------|------------------------------------------------------------------------------------------|
| `/{world}/{z}/{x}/{y}` | Fetch heatmap tile image  | `world`: nether, ~~overworld~~, ~~end~~<br>`z`: zoom level<br>`x`, `y`: tile coordinates |
| `/settings`            | Retrieve heatmap settings | None                                                                                     |

Example: https://api.2b2theatmap.info/api/nether/3/4/2

--- 

## 💻 Development

1. Prerequisites
    - Docker & Docker Compose
    - Make

2. Setup
    ```bash
    cp .env.example .env
    ```
3. Run with Live-Reload
    ```bash
    make up 
    ```
    - API: http://localhost:8000
    - Debugger (Delve): localhost:5864

--- 

## ⚡ Performance Logic

1. **Request Arrival**: Caddy handles SSL and compression.
2. **ETag Check**: If `If-None-Match` header matches the file's `mtime/size`, the service immediately returns `304 Not Modified`.
3. **Cache Lookup**: Ristretto checks if tile bytes are in RAM.
4. **Disk Fallback**: On a cache miss, the file is opened once, read into memory, stored in the cache, and served.

---

## 📊 Monitoring & Metrics

The service exports real-time telemetry to **StatsD** (integrated with VictoriaMetrics/Grafana). This allows for deep visibility into both application performance and CDN efficiency.

### HTTP Traffic & Load
* `http.requests.total` - Total incoming requests counter.
* `http.responses.total` - Total responses sent.
* `http.inflight` - **Gauge** of active requests currently being processed. Critical for monitoring server saturation.
* `http.responses.status.{code}` - Response counter partitioned by HTTP status (e.g., `200`, `304`, `404`, `500`).
* `http.errors.total` - Counter for all requests with status code `>= 400`.
* `http.duration_ms` - Histogram of request processing times in milliseconds.
* `http.duration_micros` - Histogram of request processing times in microseconds.

### CDN & Cache Efficiency
* `cache.hit` / `cache.miss`: Tracks **Ristretto** performance. A high hit ratio indicates optimal RAM utilization.
* `etag.hit` / `etag.miss`: Tracks **304 Not Modified** efficiency. High hit rates mean clients are successfully using local browser caches, saving your server's bandwidth.
* `cache.cost_bytes`: Gauge of current RAM consumption by the tile cache.

### Health Indicators
* `http.request.400.*`: Malformed requests (invalid coordinates or world names).
* `http.request.500.*`: Server-side issues (file system errors or crashes).

### Go Runtime & System (Collector)
* `cpu.percent` - CPU usage by the service process.
* `mem.rss_bytes` - Resident Set Size (actual RAM used by the process).
* `go.heap_alloc_bytes` / `go.heap_sys_bytes` - Go-specific memory management stats.
* `go.gc_pause_ns` - GC pause duration (crucial for latency-sensitive tile serving).
* `go.goroutines` / `go.threads` - Concurrency and scheduler health.
* `go.uptime_sec` - Total service uptime.

--- 

## 🔗 Links

* [Monitoring Repository](https://github.com/m1n64/2b2t-heatmap-monitoring) *Currently private*
* [2b2t Heatmaps Website](https://2b2theatmap.info)
* [2b2t Heatmaps API](https://api.2b2theatmap.info)
* [2b2t Heatmaps Data Source (nocom data)](https://github.com/nerdsinspace/nocom-explanation/blob/main/README.md)