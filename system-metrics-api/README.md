# System Metrics API

An HTTP server that exposes real-time system metrics (CPU, memory, disk) as JSON endpoints, with a live HTML dashboard.

## Features

- Individual JSON endpoints for CPU, memory, and disk metrics
- Combined `/metrics` endpoint returning all data at once
- Live HTML dashboard with auto-refresh and color-coded usage bars
- Health check endpoint with server uptime tracking

## Usage

```bash
go run main.go
```

Server starts on port 9000. Available endpoints:

| Endpoint | Description |
|---|---|
| `/health` | Health check with server uptime |
| `/cpu` | CPU model, cores, and usage % |
| `/memory` | RAM total, used, available (MB) |
| `/disk` | Disk total, used, free (GB) |
| `/metrics` | All metrics combined as JSON |
| `/dashboard` | Live HTML dashboard |

## Sample Response (`/metrics`)

```json
{
    "cpu": {
        "model_name": "Apple M2",
        "cores": 8,
        "used_percent": 10
    },
    "memory": {
        "total_mb": 8192,
        "used_mb": 6888,
        "available_mb": 1303,
        "used_percent": 84
    },
    "disk": {
        "path": "/",
        "total_gb": 228,
        "used_gb": 147,
        "free_gb": 80,
        "used_percent": 65
    }
}
```

## Dashboard

The `/dashboard` endpoint serves an HTML page with:
- Dark theme with metric cards for CPU, Memory, and Disk
- Progress bars that change color based on usage (green ≤50%, yellow ≤80%, red >80%)
- Auto-refreshes every 5 seconds

## Built With

- [Go](https://go.dev/) — `net/http`, `html/template`, `encoding/json`
- [gopsutil](https://github.com/shirou/gopsutil) — Cross-platform system metrics
