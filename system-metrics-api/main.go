package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// Structs for JSON responses

type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

type MemoryMetrics struct {
	TotalMB     uint64  `json:"total_mb"`
	UsedMB      uint64  `json:"used_mb"`
	AvailableMB uint64  `json:"available_mb"`
	UsedPercent float64 `json:"used_percent"`
}

type CPUMetrics struct {
	ModelName   string  `json:"model_name"`
	Cores       int     `json:"cores"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskMetrics struct {
	Path        string  `json:"path"`
	TotalGB     uint64  `json:"total_gb"`
	UsedGB      uint64  `json:"used_gb"`
	FreeGB      uint64  `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
}

// SystemMetrics combines all metrics into one response for /metrics and /dashboard
type SystemMetrics struct {
	CPU    CPUMetrics    `json:"cpu"`
	Memory MemoryMetrics `json:"memory"`
	Disk   DiskMetrics   `json:"disk"`
}

// Helper functions that collect metrics and return structs.
// These are reused by the individual endpoints, /metrics, and /dashboard.

func getMemoryMetrics() (MemoryMetrics, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return MemoryMetrics{}, err
	}
	return MemoryMetrics{
		TotalMB:     v.Total / 1024 / 1024,
		UsedMB:      v.Used / 1024 / 1024,
		AvailableMB: v.Available / 1024 / 1024,
		UsedPercent: math.Round(v.UsedPercent),
	}, nil
}

func getCPUMetrics() (CPUMetrics, error) {
	info, err := cpu.Info()
	if err != nil {
		return CPUMetrics{}, err
	}

	percent, err := cpu.Percent(0, false)
	if err != nil {
		return CPUMetrics{}, err
	}

	return CPUMetrics{
		ModelName:   info[0].ModelName,
		Cores:       runtime.NumCPU(),
		UsedPercent: math.Round(percent[0]),
	}, nil
}

func getDiskMetrics() (DiskMetrics, error) {
	d, err := disk.Usage("/")
	if err != nil {
		return DiskMetrics{}, err
	}
	return DiskMetrics{
		Path:        d.Path,
		TotalGB:     d.Total / 1024 / 1024 / 1024,
		UsedGB:      d.Used / 1024 / 1024 / 1024,
		FreeGB:      d.Free / 1024 / 1024 / 1024,
		UsedPercent: math.Round(d.UsedPercent),
	}, nil
}

// getAllMetrics collects CPU, memory, and disk metrics into one struct.
func getAllMetrics() (SystemMetrics, error) {
	cpuData, err := getCPUMetrics()
	if err != nil {
		return SystemMetrics{}, fmt.Errorf("cpu: %w", err)
	}

	memData, err := getMemoryMetrics()
	if err != nil {
		return SystemMetrics{}, fmt.Errorf("memory: %w", err)
	}

	diskData, err := getDiskMetrics()
	if err != nil {
		return SystemMetrics{}, fmt.Errorf("disk: %w", err)
	}

	return SystemMetrics{
		CPU:    cpuData,
		Memory: memData,
		Disk:   diskData,
	}, nil
}

// HTTP handlers

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	response := HealthResponse{
		Status: "ok",
		Uptime: fmt.Sprintf("%s", time.Since(startTime).Round(time.Second)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func memoryHandler(w http.ResponseWriter, req *http.Request) {
	data, err := getMemoryMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func cpuHandler(w http.ResponseWriter, req *http.Request) {
	data, err := getCPUMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func diskHandler(w http.ResponseWriter, req *http.Request) {
	data, err := getDiskMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// allMetricsHandler returns CPU, memory, and disk metrics as one JSON response.
func allMetricsHandler(w http.ResponseWriter, req *http.Request) {
	metrics, err := getAllMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// dashboardHandler
func dashboardHandler(w http.ResponseWriter, req *http.Request) {
	metrics, err := getAllMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))

	err = tmpl.Execute(w, metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// startTime tracks when the server started, used for the /health uptime field.
var startTime time.Time

func main() {
	startTime = time.Now()

	http.HandleFunc("/", hello)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/memory", memoryHandler)
	http.HandleFunc("/cpu", cpuHandler)
	http.HandleFunc("/disk", diskHandler)
	http.HandleFunc("/metrics", allMetricsHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

	log.Println("starting server on :9000")
	log.Println("dashboard: http://localhost:9000/dashboard")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
